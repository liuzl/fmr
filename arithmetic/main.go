package main

import (
	"errors"
	"fmt"
)

type SemanticRepresentation struct {
	Operator string
	Operands []interface{}
}

var (
	inputs = []string{
		"one plus one",
		"one plus two",
		"one plus three",
		"two plus two",
		"two plus three",
		"three plus one",
		"three plus minus two",
		"two plus two",
		"three minus two",
		"minus three minus two",
		"two times two",
		"two times three",
		"three plus three minus two",
		"minus three",
		"three plus two",
		"two times two plus three",
		"minus four",
	}

	sems = []SemanticRepresentation{
		{"+", []interface{}{1, 2}},                                                //one plus two
		{"-", []interface{}{SemanticRepresentation{"~", []interface{}{3}}, 2}},    //minus three minus two
		{"-", []interface{}{SemanticRepresentation{"+", []interface{}{3, 3}}, 2}}, //three plus three minus two
		{"+", []interface{}{SemanticRepresentation{"*", []interface{}{2, 2}}, 3}}, //two times two plus three
	}

	ops = map[string]interface{}{
		"~": func(x int) int { return -x },
		"+": func(x, y int) int { return x + y },
		"-": func(x, y int) int { return x - y },
		"*": func(x, y int) int { return x * y },
	}
)

func execute(sem interface{}) (int, error) {
	switch sem.(type) {
	case int:
		return sem.(int), nil
	case SemanticRepresentation:
		sr := sem.(SemanticRepresentation)
		if op, ok := ops[sr.Operator]; ok {
			switch op.(type) {
			case func(x int) int:
				f := op.(func(x int) int)
				if len(sr.Operands) < 1 {
					return 0, errors.New("parameter error")
				}
				n, err := execute(sr.Operands[0])
				if err != nil {
					return 0, err
				}
				return f(n), nil
			case func(x, y int) int:
				f := op.(func(x, y int) int)
				if len(sr.Operands) < 2 {
					return 0, errors.New("parameter error2")
				}
				n1, err1 := execute(sr.Operands[0])
				if err1 != nil {
					return 0, err1
				}
				n2, err2 := execute(sr.Operands[1])
				if err2 != nil {
					return 0, err2
				}
				return f(n1, n2), nil
			default:
				return 0, errors.New("func signature not correct")
			}
			return 0, nil
		} else {
			return 0, errors.New(sr.Operator + " not found")
		}
	default:
		return 0, errors.New("must be int or SR")
	}
}

func main() {
	for _, input := range inputs {
		fmt.Println(input)
	}

	for _, sem := range sems {
		ret, err := execute(sem)
		fmt.Println(sem, "=", ret, err)
	}
}
