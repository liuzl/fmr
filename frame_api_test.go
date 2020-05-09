package fmr

import (
	"testing"
)

func TestMatchFrames(t *testing.T) {
	cases := []string{
		`从北京飞上海`,
		`飞上海，从北京，后天`,
		`我要从北京走`,
	}
	g, err := GrammarFromFile("sf.grammar")
	if err != nil {
		t.Error(err)
	}
	for _, c := range cases {
		fmrs, err := g.FrameFMR(c)
		if err != nil {
			t.Error(err)
		}
		t.Log(c, fmrs)
	}
}

func TestMatchFrames2(t *testing.T) {
	cases := []string{
		`获得亚军次数降序排前5的都是哪些羽毛球运动员？`,
		`注册资本大于1亿的品牌中，哪5个品牌收入最少？并给出它们的法定代表人`,
	}
	g, err := GrammarFromFile("grammars/sql.grammar")
	if err != nil {
		t.Error(err)
	}
	for _, c := range cases {
		fmrs, err := g.FrameFMR(c)
		if err != nil {
			t.Error(err)
		}
		t.Log(c, fmrs)
	}
}
