package fmr

import (
	"testing"

	"github.com/liuzl/ling"
	"zliu.org/goutil"
)

func TestLocalGrammar(t *testing.T) {
	tests := []string{
		`天津，liang@zliu.org是我的邮箱，https://crawler.club是爬虫主页`,
		`关于FMR的介绍在这里：https://zliu.org/project/fmr/,好的`,
		`柏乡县是一个历史悠久的小城，高邑县也是，南开区呢，海淀区，思明区在哪里`,
	}
	for _, c := range tests {
		d := ling.NewDocument(c)
		if err := NLP().Annotate(d); err != nil {
			t.Error(err)
		}
		g := Grammar{}
		l, err := g.localGrammar(d)
		if err != nil {
			t.Error(err)
		}
		b, err := goutil.JSONMarshalIndent(l, "", "  ")
		if err != nil {
			t.Error(err)
		}
		t.Log(string(b))
	}
}
