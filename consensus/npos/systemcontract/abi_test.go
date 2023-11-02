package systemcontract

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/stretchr/testify/require"
)

func TestJsonUnmarshalABI(t *testing.T) {
	for _, abiStr := range []string{ValidatorsInteractiveABI, PunishInteractiveABI, SysGovInteractiveABI, AddrListInteractiveABI} {
		_, err := abi.JSON(strings.NewReader(ValidatorsInteractiveABI))
		require.NoError(t, err, abiStr)
	}
}
