package ethereum

import (
	"fmt"
	"math/big"
	"strings"
)

const (
	// HexPrefix is a ethereum data hex prefix.
	HexPrefix = "0x"
)

// DecodeHex tries to decode specified data as hex into specified container.
func DecodeHex(data string, result interface{}) error {
	if !strings.HasPrefix(data, HexPrefix) {
		return fmt.Errorf(
			"malformed hex without hex prefix: %q", data,
		)
	}

	data = data[2:]

	switch result := result.(type) {
	case *big.Int:
		_, success := result.SetString(data, 16)
		if !success {
			return fmt.Errorf(
				"unable to parse hex as big int: %q",
				data,
			)
		}

		return nil
	}

	return fmt.Errorf(
		"unsupported type for hex decode: %T",
		result,
	)
}
