package geth

import (
	"fmt"
)

type InvalidAddressError struct {
	Address string
}

func (err InvalidAddressError) Error() string {
	return fmt.Sprintf(
		`"%s" is not valid address`,
		err.Address,
	)
}
