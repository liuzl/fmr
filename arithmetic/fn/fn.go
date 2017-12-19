package main

import (
	"flag"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"log"
)

var (
	js = flag.String("js", "math.js", "javascript file")
)

func main() {
	flag.Parse()
	script, err := ioutil.ReadFile(*js)
	if err != nil {
		log.Fatal(err)
	}
	vm := otto.New()
	if _, err = vm.Run(script); err != nil {
		log.Fatal(err)
	}
	jsVal1, _ := vm.ToValue("abc")
	jsVal2, _ := vm.ToValue(3)
	result, err := vm.Call("nf.math.sum", nil, jsVal1, jsVal2)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("nf.math.sum of", jsVal1, "and", jsVal2, "is", result)
}
