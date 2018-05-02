package funcs

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
)

func init() {
	builtinFuncs["nf.math.sum"] = sum
	builtinFuncs["nf.math.mul"] = mul
	builtinFuncs["nf.math.div"] = div
	builtinFuncs["nf.math.pow"] = pow
	builtinFuncs["nf.math.neg"] = neg
	builtinFuncs["nf.math.even"] = even
	builtinFuncs["nf.math.odd"] = odd
}

func sum(x, y string) string {
	return calc(x, y, "Add")
}

func mul(x, y string) string {
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

func pow(x, y string) string {
	fx, err := strconv.ParseFloat(x, 64)
	if err != nil {
		return fmt.Sprintf("%s^%s", x, y)
	}
	fy, err := strconv.ParseFloat(y, 64)
	if err != nil {
		return fmt.Sprintf("%s^%s", x, y)
	}
	return fmt.Sprintf("%f", math.Pow(fx, fy))
}

func neg(x string) string {
	xf := new(big.Float)
	if _, err := fmt.Sscan(x, xf); err != nil {
		return ""
	}
	return xf.Neg(xf).String()
}

func even(x string) string {
	xi := new(big.Int)
	if _, err := fmt.Sscan(x, xi); err == nil && xi.Bit(0) == 0 {
		return "true"
	}
	return "false"
}

func odd(x string) string {
	xi := new(big.Int)
	if _, err := fmt.Sscan(x, xi); err == nil && xi.Bit(0) == 1 {
		return "true"
	}
	return "false"
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
