package fmr

import (
	"fmt"

	"zliu.org/goutil"
)

var builtinFuncs = make(map[string]interface{})

func init() {
	builtinFuncs["fmr.list"] = fmrList
	builtinFuncs["fmr.entity"] = fmrEntity
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

func fmrEntity(items ...interface{}) map[string]interface{} {
	l := len(items)
	if l == 0 {
		return nil
	}
	typ := fmt.Sprintf("%v", items[0])
	if typ == "" {
		return nil
	}
	if l == 1 {
		return map[string]interface{}{typ: nil}
	}
	return map[string]interface{}{typ: items[1:]}
}
