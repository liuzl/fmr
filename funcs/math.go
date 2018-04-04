package funcs

import (
	"fmt"
	"math/big"
	"strings"
)

func init() {
	builtinFuncs["nf.math.sum"] = sum
	builtinFuncs["nf.math.product"] = product
}

func sum(x, y string) string {
	if strings.Contains(x, ".") || strings.Contains(y, ".") {
		xf := new(big.Float)
		yf := new(big.Float)
		_, err := fmt.Sscan(x, xf)
		if err != nil {
			return ""
		}
		_, err = fmt.Sscan(y, yf)
		if err != nil {
			return ""
		}
		return xf.Add(xf, yf).String()
	}
	xi := new(big.Int)
	yi := new(big.Int)
	_, err := fmt.Sscan(x, xi)
	if err != nil {
		return ""
	}
	_, err = fmt.Sscan(y, yi)
	if err != nil {
		return ""
	}
	return xi.Add(xi, yi).String()
}

func product(x, y string) string {
	if strings.Contains(x, ".") || strings.Contains(y, ".") {
		xf := new(big.Float)
		yf := new(big.Float)
		_, err := fmt.Sscan(x, xf)
		if err != nil {
			return ""
		}
		_, err = fmt.Sscan(y, yf)
		if err != nil {
			return ""
		}
		return xf.Mul(xf, yf).String()
	}
	xi := new(big.Int)
	yi := new(big.Int)
	_, err := fmt.Sscan(x, xi)
	if err != nil {
		return ""
	}
	_, err = fmt.Sscan(y, yi)
	if err != nil {
		return ""
	}
	return xi.Mul(xi, yi).String()
}
