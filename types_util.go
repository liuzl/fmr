package fmr

import (
	"fmt"
	"math/big"
)

func metaEqual(m1, m2 interface{}) bool {
	if m1 == nil && m2 == nil {
		return true
	}
	if m1 != nil && m2 != nil {
		if Debug {
			fmt.Println("In Equal:", m1, m2)
		}
		switch m1.(type) {
		// meta for (any)
		case map[string]int:
			t1 := m1.(map[string]int)
			t2, ok2 := m2.(map[string]int)
			if ok2 && len(t1) == len(t2) {
				for k, v := range t1 {
					if Debug {
						fmt.Println(k, v)
					}
					if w, ok := t2[k]; !ok || v != w {
						if Debug {
							fmt.Println(v, w, ok)
						}
						return false
					}
				}
				return true
			}
			// meta for terminal text
		case string:
			s1 := m1.(string)
			s2, ok := m2.(string)
			if ok && s1 == s2 {
				return true
			}
		}
	}
	return false
}

// Equal func for Term
func (t *Term) Equal(t1 *Term) bool {
	if t == nil && t1 == nil {
		return true
	}
	if t == nil || t1 == nil {
		return false
	}
	if t.Value != t1.Value || t.Type != t1.Type {
		return false
	}
	return metaEqual(t.Meta, t1.Meta)
}

// Equal func for RuleBody
func (r *RuleBody) Equal(rb *RuleBody) bool {
	if rb == nil && r == nil {
		return true
	}
	if rb == nil || r == nil {
		return false
	}
	if len(rb.Terms) != len(r.Terms) {
		return false
	}
	for i, term := range rb.Terms {
		if !term.Equal(r.Terms[i]) {
			return false
		}
	}
	return r.F.Equal(rb.F)
}

// Equal func for FMR
func (f *FMR) Equal(fmr *FMR) bool {
	if f == nil && fmr == nil {
		return true
	}
	if !(f != nil && fmr != nil) {
		return false
	}
	if f.Fn != fmr.Fn {
		return false
	}
	if len(f.Args) != len(fmr.Args) {
		return false
	}
	for i, arg := range fmr.Args {
		if arg.Type != f.Args[i].Type {
			return false
		}
		switch arg.Type {
		case "string":
			s1, ok1 := arg.Value.(string)
			s2, ok2 := f.Args[i].Value.(string)
			if !ok1 || !ok2 || s1 != s2 {
				return false
			}
		case "int":
			s1, ok1 := arg.Value.(*big.Int)
			s2, ok2 := f.Args[i].Value.(*big.Int)
			if !ok1 || !ok2 || s1.Cmp(s2) != 0 {
				return false
			}
		case "float":
			s1, ok1 := arg.Value.(*big.Float)
			s2, ok2 := f.Args[i].Value.(*big.Float)
			if !ok1 || !ok2 || s1.Cmp(s2) != 0 {
				return false
			}
		case "index":
			s1, ok1 := arg.Value.(int)
			s2, ok2 := f.Args[i].Value.(int)
			if !ok1 || !ok2 || s1 != s2 {
				return false
			}
		case "func":
			s1, ok1 := arg.Value.(*FMR)
			s2, ok2 := f.Args[i].Value.(*FMR)
			if !ok1 || !ok2 || !s1.Equal(s2) {
				return false
			}
		}
	}
	return true
}
