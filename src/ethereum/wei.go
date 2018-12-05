package ethereum

import (
	"fmt"
	"math"
	"math/big"
)

const (
	OneEtherInWei = 1000000000000000000
)

var OneEtherInWeiFloat *big.Float

type Wei struct {
	big.Int
}

func init() {
	OneEtherInWeiFloat = big.NewFloat(0).SetInt64(1000000000000000000)
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

		integral, fractional = wei.Parts()
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

func (wei *Wei) Parts() (*big.Int, *big.Int) {
	var (
		amount  big.Int
		modulus big.Int
		divider = big.NewInt(OneEtherInWei)
	)

	amount.Set(&wei.Int)

	return amount.DivMod(&amount, divider, &modulus)
}

func (wei *Wei) Ether() float64 {
	value := big.NewFloat(0).SetInt(&wei.Int)

	amount, _ := value.Quo(
		value,
		big.NewFloat(0).SetInt64(OneEtherInWei),
	).Float64()

	return amount
}

func (wei *Wei) EtherFloatText() string {
	up := big.NewFloat(0).SetInt(&wei.Int)

	v := big.NewFloat(0)
	v = v.Quo(up, OneEtherInWeiFloat)
	return v.Text('f', 18)
}
