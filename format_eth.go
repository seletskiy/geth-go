package geth

import (
	"fmt"
	"math"
	"math/big"
)

func FormatWei(amount *big.Int, precision int) string {
	eth, wei := ConvertWeiToEth(amount)

	length := int(math.Log10(float64(OneEthInWei)))

	fractional := fmt.Sprintf(fmt.Sprintf("%%0%ds", length), wei.String())

	return fmt.Sprintf("%s.%s", eth.String(), fractional[:precision])
}
