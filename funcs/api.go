package funcs

import (
	"github.com/liuzl/goutil"
	"reflect"
)

var builtinFuncs = make(map[string]interface{})

func Call(fn string, args ...interface{}) ([]reflect.Value, error) {
	return goutil.Call(builtinFuncs, fn, args)
}
