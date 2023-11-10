package vmcaller

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
)

// ExecuteMsg executes transaction sent to system contracts.
func ExecuteMsg(msg core.Message, state *state.StateDB, header *types.Header, chainContext core.ChainContext, chainConfig *params.ChainConfig) (ret []byte, err error) {
	blockContext := core.NewEVMBlockContext(header, chainContext, nil)
	vmenv := vm.NewEVM(blockContext, core.NewEVMTxContext(msg), state, chainConfig, vm.Config{})

	ret, _, err = vmenv.Call(vm.AccountRef(msg.From()), *msg.To(), msg.Data(), msg.Gas(), msg.Value())
	// Finalise the statedb so any changes can take effect,
	// and especially if the `from` account is empty, it can be finally deleted.
	state.Finalise(true)
	return ret, WrapVMError(err, ret)
}

// WrapVMError wraps vm error with readable reason
func WrapVMError(err error, ret []byte) error {
	if err == vm.ErrExecutionReverted {
		reason, errUnpack := abi.UnpackRevert(common.CopyBytes(ret))
		if errUnpack != nil {
			reason = "internal error"
		}
		return fmt.Errorf("%s: %s", err.Error(), reason)
	}
	return err
}

// NewLegacyMessage builds a message for consensus and system governance actions, it will not consumes any fee.
func NewLegacyMessage(from common.Address, to *common.Address, nonce uint64, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte, checkNonce bool) types.Message {
	return types.NewMessage(from, to, nonce, amount, gasLimit, gasPrice, gasPrice, gasPrice, data, nil, checkNonce)
}
