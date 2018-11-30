package fmr

import (
	"fmt"

	"github.com/liuzl/goutil"
)

var builtinFuncs = make(map[string]interface{})

// Call funcs by name fn and args
func Call(fn string, args ...interface{}) (string, error) {
	ret, err := goutil.Call(builtinFuncs, fn, args...)
	if err != nil {
		return "", err
	}
	if len(ret) == 0 {
		return "", nil
	}
	return fmt.Sprintf("%+v", ret[0]), nil
}
