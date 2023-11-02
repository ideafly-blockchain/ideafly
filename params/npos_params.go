package params

import "github.com/ethereum/go-ethereum/crypto"

// NPoS proof-of-stake-authority protocol constants.
const (
	NposExtraVanity = 32                     // Fixed number of extra-data prefix bytes reserved for validator vanity
	NposExtraSeal   = crypto.SignatureLength // Fixed number of extra-data suffix bytes reserved for validator seal
)
