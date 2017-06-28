package geth

import (
	"fmt"
	"math"
	"math/big"
)

const (
	OneEtherInWei = 1000000000000000000
)

type Wei struct {
	big.Int
}

func (wei *Wei) SetFloat64(amount float64) *Wei {
	var (
		ether      big.Float
		multiplier big.Float
	)

	multiplier.SetInt64(OneEtherInWei)
	ether.SetFloat64(amount)

	ether.Mul(&ether, &multiplier)

	_, _ = ether.Int(&wei.Int)

	return wei
}

func (wei *Wei) Text(precision int) string {
	var (
		length = int(math.Log10(float64(OneEtherInWei)))

		integral, fractional = wei.Ether()
	)

	return fmt.Sprintf(
		"%s.%s",
		integral.String(),
		fmt.Sprintf(
			fmt.Sprintf("%%0%ds", length),
			fractional.String(),
		)[:precision],
	)
}

func (wei *Wei) Ether() (*big.Int, *big.Int) {
	var (
		amount  big.Int
		modulus big.Int
		divider = big.NewInt(OneEtherInWei)
	)

	amount.Set(&wei.Int)

	return amount.DivMod(&amount, divider, &modulus)
}
