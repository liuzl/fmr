package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/golang/glog"
	"github.com/liuzl/fmr"
	"github.com/robertkrimen/otto"
)

var (
	grammar = flag.String("g", "arithmetic.grammar", "grammar file")
	js      = flag.String("js", "math.js", "javascript file")
	input   = flag.String("i", "", "file of original text to read")
)

func main() {
	flag.Parse()
	b, err := ioutil.ReadFile(*grammar)
	if err != nil {
		glog.Fatal(err)
	}
	//bnf.Debug = true
	g, err := fmr.CFGrammar(string(b))
	if err != nil {
		glog.Fatal(err)
	}

	script, err := ioutil.ReadFile(*js)
	if err != nil {
		glog.Fatal(err)
	}
	vm := otto.New()
	if _, err = vm.Run(script); err != nil {
		glog.Fatal(err)
	}

	var in *os.File
	if *input == "" {
		in = os.Stdin
	} else {
		in, err = os.Open(*input)
		if err != nil {
			glog.Fatal(err)
		}
		defer in.Close()
	}
	br := bufio.NewReader(in)

	for {
		line, c := br.ReadString('\n')
		if c == io.EOF {
			break
		}
		if c != nil {
			glog.Fatal(c)
		}
		fmt.Println(line)
		//p, err := g.EarleyParse("number", line)
		ps, err := g.EarleyParseAll(line, "number")
		if err != nil {
			glog.Fatal(err)
		}
		for i, p := range ps {
			trees := p.GetTrees()
			//fmt.Printf("%+v\n", p)
			fmt.Printf("p%d tree number:%d\n", i, len(trees))
			for _, tree := range trees {
				//tree.Print(os.Stdout)
				sem, err := tree.Semantic()
				if err != nil {
					glog.Fatal(err)
				}
				result, err := vm.Run(sem)
				if err != nil {
					glog.Fatal(err)
				}
				fmt.Printf("%s = %v\n", sem, result)
			}
		}
		fmt.Println()
	}
}
