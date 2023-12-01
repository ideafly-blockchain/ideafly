package npos

import (
	"errors"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/npos/systemcontract"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

// call this at epoch block to get top validators based on the state of epoch block - 1
func (c *Npos) getTopValidators(ctx *systemcontract.CallContext) ([]common.Address, error) {
	parent := ctx.ChainContext.GetHeader(ctx.Header.ParentHash, ctx.Header.Number.Uint64()-1)
	if parent == nil {
		return []common.Address{}, consensus.ErrUnknownAncestor
	}
	statedb, err := c.stateFn(parent.Root)
	if err != nil {
		return []common.Address{}, err
	}

	method := "getTopValidators"
	data, err := c.abi[systemcontract.ValidatorsContractName].Pack(method)
	if err != nil {
		log.Error("Can't pack data for getTopValidators", "error", err)
		return []common.Address{}, err
	}

	// use parent statedb
	newCtx := &systemcontract.CallContext{
		Statedb:      statedb,
		Header:       ctx.Header,
		ChainContext: ctx.ChainContext,
		ChainConfig:  ctx.ChainConfig,
	}

	result, err := systemcontract.VmCall(newCtx, systemcontract.ValidatorsContractAddr, data)
	if err != nil {
		log.Error("Can't read top validators", "err", err)
		return []common.Address{}, err
	}

	// unpack data
	ret, err := c.abi[systemcontract.ValidatorsContractName].Unpack(method, result)
	if err != nil {
		return []common.Address{}, err
	}
	if len(ret) != 1 {
		return []common.Address{}, errors.New("invalid params length")
	}
	validators, ok := ret[0].([]common.Address)
	if !ok {
		return []common.Address{}, errors.New("invalid validators format")
	}
	sort.Sort(validatorsAscending(validators))
	return validators, err
}

func (c *Npos) updateValidators(ctx *systemcontract.CallContext, vals []common.Address) error {
	// method
	method := "updateActiveValidatorSet"
	data, err := c.abi[systemcontract.ValidatorsContractName].Pack(method, vals, new(big.Int).SetUint64(c.config.Epoch))
	if err != nil {
		log.Error("Can't pack data for updateActiveValidatorSet", "error", err)
		return err
	}

	_, err = systemcontract.VmCall(ctx, systemcontract.ValidatorsContractAddr, data)
	if err != nil {
		log.Error("Can't update validators to contract", "err", err)
		return err
	}
	return nil
}

func (c *Npos) trySendBlockReward(ctx *systemcontract.CallContext) error {
	fee := ctx.Statedb.GetBalance(consensus.FeeRecoder)
	if fee.Cmp(common.Big0) <= 0 {
		return nil
	}

	// Caller will send tx to deposit block fees to contract, add to his balance first.
	ctx.Statedb.AddBalance(systemcontract.EngineCaller, fee)
	// reset fee
	ctx.Statedb.SetBalance(consensus.FeeRecoder, common.Big0)

	method := "distributeBlockReward"
	data, err := c.abi[systemcontract.ValidatorsContractName].Pack(method)
	if err != nil {
		log.Error("Can't pack data for distributeBlockReward", "err", err)
		return err
	}

	_, err = systemcontract.VmCallWithValue(ctx, systemcontract.EngineCaller, systemcontract.ValidatorsContractAddr, data, fee)

	if err != nil {
		return err
	}
	return nil
}

func (c *Npos) punishValidator(ctx *systemcontract.CallContext, val common.Address) error {
	// method
	method := "punish"
	data, err := c.abi[systemcontract.PunishContractName].Pack(method, val)
	if err != nil {
		log.Error("Can't pack data for punish", "error", err)
		return err
	}

	_, err = systemcontract.VmCall(ctx, systemcontract.PunishContractAddr, data)
	if err != nil {
		log.Error("can't punish validator", "err", err)
		return err
	}

	return nil
}

func (c *Npos) decreaseMissedBlocksCounter(ctx *systemcontract.CallContext) error {
	// method
	method := "decreaseMissedBlocksCounter"
	data, err := c.abi[systemcontract.PunishContractName].Pack(method, new(big.Int).SetUint64(c.config.Epoch))
	if err != nil {
		log.Error("Can't pack data for decreaseMissedBlocksCounter", "error", err)
		return err
	}

	_, err = systemcontract.VmCall(ctx, systemcontract.PunishContractAddr, data)
	if err != nil {
		log.Error("Can't decrease missed blocks counter for validator", "err", err)
		return err
	}

	return nil
}

func (c *Npos) tryPunishValidator(ctx *systemcontract.CallContext, chain consensus.ChainHeaderReader) error {
	number := ctx.Header.Number.Uint64()
	snap, err := c.snapshot(chain, number-1, ctx.Header.ParentHash, nil)
	if err != nil {
		return err
	}
	validators := snap.validators()
	outTurnValidator := validators[number%uint64(len(validators))]
	// check sigend recently or not
	signedRecently := false
	for _, recent := range snap.Recents {
		if recent == outTurnValidator {
			signedRecently = true
			break
		}
	}
	if !signedRecently {
		if err := c.punishValidator(ctx, outTurnValidator); err != nil {
			return err
		}
	}

	return nil
}

// syncWithSysContractAtEpoch: set current validators to system contract, decreaseMissedBlocksCounter, get and return next epoch validators
func (c *Npos) syncWithSysContractAtEpoch(ctx *systemcontract.CallContext, chain consensus.ChainHeaderReader) ([]common.Address, error) {
	// NPoS use a look-back validators set for safety(when supporting fast-sync).
	// the authorized validators come from the header.Extra at block currentNum - EPOCH.
	checkpointHeader := chain.GetHeaderByNumber(ctx.Header.Number.Uint64() - c.config.Epoch)
	if checkpointHeader == nil {
		return nil, consensus.ErrUnknownAncestor
	}
	// get validators from headers and use that for new validator set
	validators := make([]common.Address, (len(checkpointHeader.Extra)-extraVanity-extraSeal)/common.AddressLength)
	for i := 0; i < len(validators); i++ {
		copy(validators[i][:], checkpointHeader.Extra[extraVanity+i*common.AddressLength:])
	}
	if len(validators) < 1 {
		return []common.Address{}, errInvalidExtraValidators
	}
	if err := c.updateValidators(ctx, validators); err != nil {
		return []common.Address{}, err
	}

	//  decrease validator missed blocks counter at epoch
	if err := c.decreaseMissedBlocksCounter(ctx); err != nil {
		return []common.Address{}, err
	}

	nextEpochValidators, err := c.getTopValidators(ctx)
	if err != nil {
		return []common.Address{}, err
	}

	return nextEpochValidators, nil
}

// initializeSystemContracts initializes all genesis system contracts.
func (c *Npos) initializeSystemContracts(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB) error {
	snap, err := c.snapshot(chain, 0, header.ParentHash, nil)
	if err != nil {
		return err
	}

	genesisValidators := snap.validators()
	if len(genesisValidators) == 0 || len(genesisValidators) > maxValidators {
		return errInvalidValidatorsLength
	}

	method := "initialize"
	contracts := []struct {
		addr    common.Address
		packFun func() ([]byte, error)
	}{
		{systemcontract.ValidatorsContractAddr, func() ([]byte, error) {
			gvs := c.config.GenesisValidators
			if len(gvs) == 0 {
				return nil, errors.New("missing genesis validators")
			}
			vals, managers := make([]common.Address, 0, len(gvs)), make([]common.Address, 0, len(gvs))
			for _, v := range gvs {
				vals = append(vals, v.Validator)
				managers = append(managers, v.Manager)
			}
			stakingAdmin := c.config.StakingAdmin
			return c.abi[systemcontract.ValidatorsContractName].Pack(method, vals, managers, stakingAdmin)
		}},
		{systemcontract.PunishContractAddr, func() ([]byte, error) { return c.abi[systemcontract.PunishContractName].Pack(method) }},
		{systemcontract.AddressListContractAddr, func() ([]byte, error) {
			return c.abi[systemcontract.AddressListContractName].Pack(method, c.config.GovAdmin)
		}},
		{systemcontract.SysGovContractAddr, func() ([]byte, error) {
			return c.abi[systemcontract.SysGovContractName].Pack(method, c.config.GovAdmin)
		}},
	}

	ctx := &systemcontract.CallContext{
		Statedb:      state,
		Header:       header,
		ChainContext: newChainContext(chain, c),
		ChainConfig:  c.chainConfig,
	}

	for _, contract := range contracts {
		data, err := contract.packFun()
		if err != nil {
			log.Error("system contract initialize pack failed", "contract", contract.addr, "err", err)
			return err
		}

		_, err = systemcontract.VmCall(ctx, contract.addr, data)

		if err != nil {
			log.Error("system contract initialize failed", "contract", contract.addr, "err", err)
			return err
		}
	}

	return nil
}
