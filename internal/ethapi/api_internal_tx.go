package ethapi

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

// ActionFilter is the filter config for internal txs API. It holds one more
// field to filters the result.
type ActionFilter struct {
	From     *common.Address
	To       *common.Address
	OpCode   *string // CALL/CALLCODE/DELEGATECALL/STATICCALL/CREATE/CREATE2/SELFDESTRUCT
	MinValue *big.Int
}

// GetInternalTxs return internal txs by block number or hash
func (api *BlockChainAPI) GetInternalTxs(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash, filter *ActionFilter) (types.InternalTxs, error) {
	block, err := api.b.BlockByNumberOrHash(ctx, blockNrOrHash)
	if err != nil {
		return nil, err
	}

	iTx, err := api.getInnerTx(block)
	if err != nil {
		return nil, err
	}
	if filter == nil {
		return iTx, nil
	}

	res := make([]*types.InternalTx, 0)
	for _, tx := range iTx {
		tx.Actions = api.filterAction(tx.Actions, filter)
		if len(tx.Actions) > 0 {
			res = append(res, tx)
		}
	}

	return res, nil
}

// GetInternalTxsByTxHash return internal txs by tx hash
func (api *BlockChainAPI) GetInternalTxsByTxHash(ctx context.Context, hash common.Hash, filter *ActionFilter) (*types.InternalTx, error) {
	tx, blkHash, _, _, err := api.b.GetTransaction(ctx, hash)
	if err != nil {
		return nil, err
	}

	if tx == nil {
		return nil, fmt.Errorf("tx #%s not found", hash)
	}

	block, err := api.b.BlockByHash(ctx, blkHash)
	if err != nil {
		return nil, err
	}

	// Trace the block if it was found
	if block == nil {
		return nil, fmt.Errorf("block #%s not found", hash)
	}

	txs, err := api.getInnerTx(block)
	if err != nil {
		return nil, err
	}

	for _, t := range txs {
		if t.TxHash == hash {
			t.Actions = api.filterAction(t.Actions, filter)
			return t, nil
		}
	}

	return nil, nil
}

func (api *BlockChainAPI) filterAction(actions []*types.Action, filter *ActionFilter) []*types.Action {
	if filter == nil {
		return actions
	}

	res := make([]*types.Action, 0, len(actions))

	for _, act := range actions {
		if filter.OpCode != nil && *filter.OpCode != act.OpCode {
			continue
		}

		if filter.MinValue != nil && (act.Value == nil || filter.MinValue.Cmp(act.Value) > 0) {
			continue
		}

		if filter.From != nil && *filter.From != act.From {
			continue
		}

		if filter.To != nil && *filter.To != act.To {
			continue
		}

		res = append(res, act)
	}

	return res
}

// getInnerTx returns internal txs
func (api *BlockChainAPI) getInnerTx(block *types.Block) (types.InternalTxs, error) {
	txs := rawdb.ReadInternalTxs(api.b.ChainDb(), block.Hash(), block.NumberU64())

	for _, tx := range txs {
		tx.BlockHash = block.Hash()
		tx.BlockNumber = block.Number()
	}

	return txs, nil
}
