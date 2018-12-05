package ethereum

import (
	"encoding/hex"
	"fmt"
	"math/big"
)

// EncodeHex tries to encode passed value into 0x prefixed hex data string.
func EncodeHex(value interface{}) (string, error) {
	switch value := value.(type) {
	case *Wei:
		return EncodeHex(&value.Int)

	case *big.Int:
		return "0x" + value.Text(16), nil

	case []byte:
		return "0x" + hex.EncodeToString(value), nil
	}

	return "", fmt.Errorf(
		"unsupported type for hex encode: %T",
		value,
	)
}
