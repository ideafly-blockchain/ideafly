package systemcontract

import (
	"errors"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
)

type CallContext struct {
	Statedb      *state.StateDB
	Header       *types.Header
	ChainContext core.ChainContext
	ChainConfig  *params.ChainConfig
}

// VmCall is used for the consensus engine to interact with system contracts.
// It will just to use the EngineCaller as the `from` address and the transfer amount is 0.
// If need to use a different address, then you should call `VmCallWithValue`
func VmCall(ctx *CallContext, to common.Address, data []byte) (ret []byte, err error) {
	return VmCallWithValue(ctx, EngineCaller, to, data, big.NewInt(0))
}

func VmCallWithValue(ctx *CallContext, from common.Address, to common.Address, data []byte, value *big.Int) (ret []byte, err error) {
	if ctx == nil || ctx.Statedb == nil || ctx.Header == nil || ctx.ChainConfig == nil {
		return nil, errors.New("missing required call context")
	}
	blockContext := core.NewEVMBlockContext(ctx.Header, ctx.ChainContext, nil)
	vmenv := vm.NewEVM(blockContext, vm.TxContext{
		Origin:   from,
		GasPrice: big.NewInt(0),
	}, ctx.Statedb, ctx.ChainConfig, vm.Config{})

	ret, _, err = vmenv.Call(vm.AccountRef(from), to, data, math.MaxUint64, value)
	// Finalise the statedb so any changes can take effect,
	// and especially if the `from` account is empty, it can be finally deleted.
	ctx.Statedb.Finalise(true)
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
