package types

import "github.com/ethereum/go-ethereum/common"

// EvmExtraValidator contains some extra validations to a transaction,
// and the validator is used inside the evm.
type EvmExtraValidator interface {
	// IsAddressBanned returns whether an address is banned.
	IsAddressBanned(address common.Address, cType common.AddressCheckType) bool
	// IsAddressBannedFromLog returns whether a log (contract event) is banned.
	IsAddressBannedFromLog(log *Log) bool
}
