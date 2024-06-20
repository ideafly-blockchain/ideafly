package npos

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"runtime/debug"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/npos/systemcontract"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	getblacklistTimer = metrics.NewRegisteredTimer("congress/blacklist/get", nil)
	getRulesTimer     = metrics.NewRegisteredTimer("congress/eventcheckrules/get", nil)
)

// Proposal is the system governance proposal info.
type Proposal struct {
	Id     *big.Int
	Action *big.Int
	From   common.Address
	To     common.Address
	Value  *big.Int
	Data   []byte
}

func (c *Npos) getPassedProposalCount(ctx *systemcontract.CallContext) (uint32, error) {
	method := "getPassedProposalCount"
	data, err := c.abi[systemcontract.SysGovContractName].Pack(method)
	if err != nil {
		log.Error("Can't pack data for getPassedProposalCount", "error", err)
		return 0, err
	}

	result, err := systemcontract.VmCall(ctx, systemcontract.SysGovContractAddr, data)
	if err != nil {
		return 0, err
	}
	// unpack data
	ret, err := c.abi[systemcontract.SysGovContractName].Unpack(method, result)
	if err != nil {
		return 0, err
	}
	if len(ret) != 1 {
		return 0, errors.New("invalid output length")
	}
	count, ok := ret[0].(uint32)
	if !ok {
		return 0, errors.New("invalid count format")
	}

	return count, nil
}

func (c *Npos) getPassedProposalByIndex(ctx *systemcontract.CallContext, idx uint32) (*Proposal, error) {
	method := "getPassedProposalByIndex"
	data, err := c.abi[systemcontract.SysGovContractName].Pack(method, idx)
	if err != nil {
		log.Error("Can't pack data for getPassedProposalByIndex", "error", err)
		return nil, err
	}

	result, err := systemcontract.VmCall(ctx, systemcontract.SysGovContractAddr, data)
	if err != nil {
		return nil, err
	}
	// unpack data
	prop := &Proposal{}
	err = c.abi[systemcontract.SysGovContractName].UnpackIntoInterface(prop, method, result)
	if err != nil {
		return nil, err
	}

	return prop, nil
}

// finishProposalById
func (c *Npos) finishProposalById(ctx *systemcontract.CallContext, id *big.Int) error {
	method := "finishProposalById"
	data, err := c.abi[systemcontract.SysGovContractName].Pack(method, id)
	if err != nil {
		log.Error("Can't pack data for getPassedProposalByIndex", "error", err)
		return err
	}

	_, err = systemcontract.VmCall(ctx, systemcontract.SysGovContractAddr, data)
	if err != nil {
		return err
	}

	return nil
}

func (c *Npos) executeProposal(ctx *systemcontract.CallContext, prop *Proposal, totalTxIndex int) (*types.Transaction, *types.Receipt, error) {
	// Even if the miner is not `running`, it's still working,
	// the 'miner.worker' will try to FinalizeAndAssemble a block,
	// in this case, the signTxFn is not set. A `non-miner node` can't execute system governance proposal.
	if c.signTxFn == nil {
		return nil, nil, errors.New("signTxFn not set")
	}

	propRLP, err := rlp.EncodeToBytes(prop)
	if err != nil {
		return nil, nil, err
	}
	//make system governance transaction
	nonce := ctx.Statedb.GetNonce(c.validator)

	tx := types.NewTransaction(nonce, systemcontract.SysGovToAddr, new(big.Int), ctx.Header.GasLimit, new(big.Int), propRLP)
	tx, err = c.signTxFn(accounts.Account{Address: c.validator}, tx, c.chainConfig.ChainID)
	if err != nil {
		return nil, nil, err
	}
	//add nonce for validator
	ctx.Statedb.SetNonce(c.validator, nonce+1)
	receipt := c.executeProposalMsg(ctx, prop, totalTxIndex, tx.Hash(), ctx.Header.Hash())

	return tx, receipt, nil
}

func (c *Npos) replayProposal(ctx *systemcontract.CallContext, prop *Proposal, totalTxIndex int, tx *types.Transaction) (*types.Receipt, error) {
	sender, err := types.Sender(c.signer, tx)
	if err != nil {
		return nil, err
	}
	if sender != ctx.Header.Coinbase {
		return nil, errors.New("invalid sender for system governance transaction")
	}
	propRLP, err := rlp.EncodeToBytes(prop)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(propRLP, tx.Data()) {
		return nil, fmt.Errorf("data missmatch, proposalID: %s, rlp: %s, txHash:%s, txData:%s", prop.Id.String(), hexutil.Encode(propRLP), tx.Hash().String(), hexutil.Encode(tx.Data()))
	}
	//make system governance transaction
	nonce := ctx.Statedb.GetNonce(sender)
	//add nonce for validator
	ctx.Statedb.SetNonce(sender, nonce+1)
	receipt := c.executeProposalMsg(ctx, prop, totalTxIndex, tx.Hash(), ctx.Header.Hash())

	return receipt, nil
}

func (c *Npos) executeProposalMsg(ctx *systemcontract.CallContext, prop *Proposal, totalTxIndex int, txHash, bHash common.Hash) *types.Receipt {
	var receipt *types.Receipt
	action := prop.Action.Uint64()
	switch action {
	case 0:
		// evm action.
		receipt = c.executeEvmCallProposal(ctx, prop, totalTxIndex, txHash, bHash)
	case 1:
		// delete code action
		ok := ctx.Statedb.Erase(prop.To)
		receipt = types.NewReceipt([]byte{}, !ok, ctx.Header.GasUsed)
		log.Info("executeProposalMsg", "action", "erase", "id", prop.Id.String(), "to", prop.To, "txHash", txHash.String(), "success", ok)
	default:
		receipt = types.NewReceipt([]byte{}, true, ctx.Header.GasUsed)
		log.Warn("executeProposalMsg failed, unsupported action", "action", action, "id", prop.Id.String(), "from", prop.From, "to", prop.To, "value", prop.Value.String(), "data", hexutil.Encode(prop.Data), "txHash", txHash.String())
	}

	receipt.TxHash = txHash
	receipt.BlockHash = bHash
	receipt.BlockNumber = ctx.Header.Number
	receipt.TransactionIndex = uint(ctx.Statedb.TxIndex())

	return receipt
}

// the returned value should not nil.
func (c *Npos) executeEvmCallProposal(ctx *systemcontract.CallContext, prop *Proposal, totalTxIndex int, txHash, bHash common.Hash) *types.Receipt {
	// actually call evm with proposal infos
	ctx.Statedb.SetTxContext(txHash, totalTxIndex)
	_, err := systemcontract.VmCallWithValue(ctx, prop.From, prop.To, prop.Data, prop.Value)
	// governance message will not actually consumes gas
	receipt := types.NewReceipt([]byte{}, err != nil, ctx.Header.GasUsed)
	// Set the receipt logs and create a bloom for filtering
	receipt.Logs = ctx.Statedb.GetLogs(txHash, bHash)
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})

	log.Info("executeProposalMsg", "action", "evmCall", "id", prop.Id.String(), "from", prop.From, "to", prop.To, "value", prop.Value.String(), "data", hexutil.Encode(prop.Data), "txHash", txHash.String(), "err", err)

	return receipt
}

// IsSysTransaction checks whether a specific transaction is a system transaction.
func (c *Npos) IsSysTransaction(sender common.Address, tx *types.Transaction, header *types.Header) (bool, error) {
	if tx.To() == nil {
		return false, nil
	}

	to := tx.To()
	if sender == header.Coinbase && *to == systemcontract.SysGovToAddr && tx.GasPrice().Sign() == 0 {
		return true, nil
	}
	// Make sure the miner can NOT call the system contract through a normal transaction.
	if sender == header.Coinbase && *to == systemcontract.SysGovContractAddr {
		return true, nil
	}
	return false, nil
}

// Methods for debug trace

// ApplySysTx applies a system-transaction using a given evm,
// the main purpose of this method is for tracing a system-transaction.
func (c *Npos) ApplySysTx(evm *vm.EVM, state *state.StateDB, txIndex int, sender common.Address, tx *types.Transaction) (ret []byte, vmerr error, err error) {
	var prop = &Proposal{}
	if err = rlp.DecodeBytes(tx.Data(), prop); err != nil {
		return
	}
	evm.Context.ExtraValidator = nil
	nonce := evm.StateDB.GetNonce(sender)
	//add nonce for validator
	evm.StateDB.SetNonce(sender, nonce+1)

	action := prop.Action.Uint64()
	switch action {
	case 0:
		// evm action.
		// actually run the governance message
		state.SetTxContext(tx.Hash(), txIndex)
		evm.TxContext = vm.TxContext{
			Origin:   prop.From,
			GasPrice: new(big.Int),
		}
		ret, _, vmerr = evm.Call(vm.AccountRef(prop.From), prop.To, prop.Data, tx.Gas(), prop.Value)
		state.Finalise(true)
	case 1:
		// delete code action
		_ = state.Erase(prop.To)
	default:
		vmerr = errors.New("unsupported action")
	}
	return
}

// ValidateTx do a consensus-related validation on the given transaction at the given header and state.
// the parentState must be the state of the header's parent block.
func (c *Npos) ValidateTx(sender common.Address, tx *types.Transaction, header *types.Header, parentState *state.StateDB) error {
	// Must use the parent state for current validation
	m, err := c.getBlacklist(header, parentState)
	if err != nil {
		return err
	}
	if d, exist := m[sender]; exist && (d != DirectionTo) {
		log.Trace("Hit blacklist", "tx", tx.Hash().String(), "addr", sender.String(), "direction", d)
		return types.ErrAddressBanned
	}
	if to := tx.To(); to != nil {
		if d, exist := m[*to]; exist && (d != DirectionFrom) {
			log.Trace("Hit blacklist", "tx", tx.Hash().String(), "addr", to.String(), "direction", d)
			return types.ErrAddressBanned
		}
	}
	return nil
}

func (c *Npos) getBlacklist(header *types.Header, parentState *state.StateDB) (map[common.Address]bannedDirection, error) {
	defer func(start time.Time) {
		getblacklistTimer.UpdateSince(start)
	}(time.Now())

	if v, ok := c.blacklists.Get(header.ParentHash); ok {
		return v.(map[common.Address]bannedDirection), nil
	}

	c.blLock.Lock()
	defer c.blLock.Unlock()
	if v, ok := c.blacklists.Get(header.ParentHash); ok {
		return v.(map[common.Address]bannedDirection), nil
	}

	// if the last updates is long ago, we don't need to get blacklist from the contract.
	num := header.Number.Uint64()
	lastUpdated := lastBlacklistUpdatedNumber(parentState)
	if num >= 2 && num > lastUpdated+1 {
		parent := c.chain.GetHeader(header.ParentHash, num-1)
		if parent != nil {
			if v, ok := c.blacklists.Get(parent.ParentHash); ok {
				m := v.(map[common.Address]bannedDirection)
				c.blacklists.Add(header.ParentHash, m)
				return m, nil
			}
		} else {
			log.Error("Unexpected error when getBlacklist, can not get parent from chain", "number", num, "blockHash", header.Hash(), "parentHash", header.ParentHash)
		}
	}

	// can't get blacklist from cache, try to call the contract
	alABI := c.abi[systemcontract.AddressListContractName]
	get := func(method string) ([]common.Address, error) {
		ret, err := c.commonCallContract(header, parentState, alABI, systemcontract.AddressListContractAddr, method, 1)
		if err != nil {
			log.Error(fmt.Sprintf("%s failed", method), "err", err)
			log.Error(string(debug.Stack()))
			return nil, err
		}

		blacks, ok := ret[0].([]common.Address)
		if !ok {
			return []common.Address{}, errors.New("invalid blacklist format")
		}
		return blacks, nil
	}
	froms, err := get("getBlacksFrom")
	if err != nil {
		return nil, err
	}
	tos, err := get("getBlacksTo")
	if err != nil {
		return nil, err
	}

	m := make(map[common.Address]bannedDirection)
	for _, from := range froms {
		m[from] = DirectionFrom
	}
	for _, to := range tos {
		if _, exist := m[to]; exist {
			m[to] = DirectionBoth
		} else {
			m[to] = DirectionTo
		}
	}
	c.blacklists.Add(header.ParentHash, m)
	return m, nil
}

func (c *Npos) CreateEvmExtraValidator(header *types.Header, parentState *state.StateDB) types.EvmExtraValidator {
	blacks, err := c.getBlacklist(header, parentState)
	if err != nil {
		log.Error("getBlacklist failed", "err", err)
		return nil
	}
	rules, err := c.getEventCheckRules(header, parentState)
	if err != nil {
		log.Error("getEventCheckRules failed", "err", err)
		return nil
	}
	return &daoRulesValidator{
		blacks: blacks,
		rules:  rules,
	}
}

func (c *Npos) getEventCheckRules(header *types.Header, parentState *state.StateDB) (map[common.Hash]*EventCheckRule, error) {
	defer func(start time.Time) {
		getRulesTimer.UpdateSince(start)
	}(time.Now())

	if v, ok := c.eventCheckRules.Get(header.ParentHash); ok {
		return v.(map[common.Hash]*EventCheckRule), nil
	}

	c.rulesLock.Lock()
	defer c.rulesLock.Unlock()
	if v, ok := c.eventCheckRules.Get(header.ParentHash); ok {
		return v.(map[common.Hash]*EventCheckRule), nil
	}

	// if the last updates is long ago, we don't need to get blacklist from the contract.
	num := header.Number.Uint64()
	lastUpdated := lastRulesUpdatedNumber(parentState)
	if num >= 2 && num > lastUpdated+1 {
		parent := c.chain.GetHeader(header.ParentHash, num-1)
		if parent != nil {
			if v, ok := c.eventCheckRules.Get(parent.ParentHash); ok {
				m := v.(map[common.Hash]*EventCheckRule)
				c.eventCheckRules.Add(header.ParentHash, m)
				return m, nil
			}
		} else {
			log.Error("Unexpected error when getEventCheckRules, can not get parent from chain", "number", num, "blockHash", header.Hash(), "parentHash", header.ParentHash)
		}
	}

	// can't get blacklist from cache, try to call the contract
	alABI := c.abi[systemcontract.AddressListContractName]
	method := "getRuleByIndex"
	get := func(i uint32) (common.Hash, int, common.AddressCheckType, error) {
		ret, err := c.commonCallContract(header, parentState, alABI, systemcontract.AddressListContractAddr, method, 3, i)
		if err != nil {
			return common.Hash{}, 0, common.CheckNone, err
		}
		sig := ret[0].([32]byte)
		idx := ret[1].(*big.Int).Uint64()
		ct := ret[2].(uint8)

		return sig, int(idx), common.AddressCheckType(ct), nil
	}

	cnt, err := c.getEventCheckRulesLen(header, parentState)
	if err != nil {
		log.Error("getEventCheckRulesLen failed", "err", err)
		return nil, err
	}
	rules := make(map[common.Hash]*EventCheckRule)
	for i := 0; i < cnt; i++ {
		sig, idx, ct, err := get(uint32(i))
		if err != nil {
			log.Error("getRuleByIndex failed", "index", i, "number", num, "blockHash", header.Hash(), "err", err)
			return nil, err
		}
		rule, exist := rules[sig]
		if !exist {
			rule = &EventCheckRule{
				EventSig: sig,
				Checks:   make(map[int]common.AddressCheckType),
			}
			rules[sig] = rule
		}
		rule.Checks[idx] = ct
	}

	c.eventCheckRules.Add(header.ParentHash, rules)
	return rules, nil
}

func (c *Npos) getEventCheckRulesLen(header *types.Header, parentState *state.StateDB) (int, error) {
	ret, err := c.commonCallContract(header, parentState, c.abi[systemcontract.AddressListContractName], systemcontract.AddressListContractAddr, "rulesLen", 1)
	if err != nil {
		return 0, err
	}
	ln, ok := ret[0].(uint32)
	if !ok {
		return 0, fmt.Errorf("unexpected output type, value: %v", ret[0])
	}
	return int(ln), nil
}

// commonCallContract only for interacting with those contracts which do no includes `BLOCKHASH` opcode
func (c *Npos) commonCallContract(header *types.Header, statedb *state.StateDB, contractABI abi.ABI, addr common.Address, method string, expectResultLen int, args ...interface{}) ([]interface{}, error) {
	data, err := contractABI.Pack(method, args...)
	if err != nil {
		log.Error("Can't pack data ", "method", method, "err", err)
		return nil, err
	}

	//the newMinimalChainContext is only for interacting with those contracts which do no includes `BLOCKHASH` opcode
	ctx := &systemcontract.CallContext{
		Statedb:      statedb,
		Header:       header,
		ChainContext: newMinimalChainContext(c),
		ChainConfig:  c.chainConfig,
	}
	result, err := systemcontract.VmCall(ctx, addr, data)

	if err != nil {
		return nil, err
	}

	// unpack data
	ret, err := contractABI.Unpack(method, result)
	if err != nil {
		return nil, err
	}
	if len(ret) != expectResultLen {
		return nil, errors.New("invalid result length")
	}
	return ret, nil
}

func lastBlacklistUpdatedNumber(state consensus.StateReader) uint64 {
	value := state.GetState(systemcontract.AddressListContractAddr, systemcontract.BlackLastUpdatedNumberPosition)
	return value.Big().Uint64()
}

func lastRulesUpdatedNumber(state consensus.StateReader) uint64 {
	value := state.GetState(systemcontract.AddressListContractAddr, systemcontract.RulesLastUpdatedNumberPosition)
	return value.Big().Uint64()
}
