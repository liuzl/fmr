package fmr

import (
	"github.com/liuzl/goutil"
)

var builtinFuncs = make(map[string]interface{})

func init() {
	builtinFuncs["fmr.list"] = fmrList
	builtinFuncs["fmr.entity"] = fmrList
}

// Call funcs by name fn and args
func Call(fn string, args ...interface{}) (interface{}, error) {
	ret, err := goutil.Call(builtinFuncs, fn, args...)
	if err != nil {
		return nil, err
	}
	if len(ret) == 0 {
		return nil, nil
	}
	return ret[0].Interface(), nil
}

func fmrList(items ...interface{}) []interface{} {
	return items
}
