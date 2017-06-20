package geth

import "math/big"

const (
	OneEthInWei = 1000000000000000000
)

func ConvertWeiToEth(amount *big.Int) (*big.Int, *big.Int) {
	divider := big.NewInt(OneEthInWei)

	var wei big.Int

	return amount.DivMod(amount, divider, &wei)
}
