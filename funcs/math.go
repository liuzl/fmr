package funcs

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

func init() {
	builtinFuncs["nf.math.sum"] = sum
	builtinFuncs["nf.math.product"] = product
	builtinFuncs["nf.math.div"] = div
}

func sum(x, y string) string {
	return calc(x, y, "Add")
}

func product(x, y string) string {
	return calc(x, y, "Mul")
}

func div(x, y string) string {
	fx, err := strconv.ParseFloat(x, 64)
	if err != nil {
		return fmt.Sprintf("%s/%s", x, y)
	}
	fy, err := strconv.ParseFloat(y, 64)
	if err != nil || fy == 0 {
		return fmt.Sprintf("%s/%s", x, y)
	}
	return fmt.Sprintf("%f", fx/fy)
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
