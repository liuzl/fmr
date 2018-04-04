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
	return calc(x, y, "Add")
}

func product(x, y string) string {
	return calc(x, y, "Mul")
}

func calc(x, y, method string) string {
	if strings.Contains(x, ".") || strings.Contains(y, ".") {
		xf, yf := new(big.Float), new(big.Float)
		if _, err := fmt.Sscan(x, xf); err != nil {
			return ""
		}
		if _, err := fmt.Sscan(y, yf); err != nil {
			return ""
		}
		switch method {
		case "Add":
			return xf.Add(xf, yf).String()
		case "Mul":
			return xf.Mul(xf, yf).String()
		default:
			return ""
		}
	}
	xi, yi := new(big.Int), new(big.Int)
	if _, err := fmt.Sscan(x, xi); err != nil {
		return ""
	}
	if _, err := fmt.Sscan(y, yi); err != nil {
		return ""
	}
	switch method {
	case "Add":
		return xi.Add(xi, yi).String()
	case "Mul":
		return xi.Mul(xi, yi).String()
	default:
		return ""
	}
}
