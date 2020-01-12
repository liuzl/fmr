package fmr

import (
	"testing"
)

func TestListFunc(t *testing.T) {
	t.Log(Call("fmr.list", "100000227", 78, "abc"))
}

func TestEntityFunc(t *testing.T) {
	t.Log(Call("fmr.entity", "PER", map[string]string{"name": "冯诺依曼"}))
}
