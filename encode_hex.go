package geth

import (
	"encoding/hex"
	"fmt"
)

// EncodeHex tries to encode passed value into 0x prefixed hex data string.
func EncodeHex(value interface{}) (string, error) {
	switch value := value.(type) {
	case []byte:
		return "0x" + hex.EncodeToString(value), nil
	}

	return "", fmt.Errorf(
		"unsupported type for hex encode: %T",
		value,
	)
}
