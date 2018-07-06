package bnf

import (
	"testing"

	"github.com/liuzl/ling"
)

func TestLocalGrammar(t *testing.T) {
	tests := []string{
		`北京，liang@zliu.org是我的邮箱，https://crawler.club是爬虫主页`,
		`关于FMR的介绍在这里：https://zliu.org/project/fmr/,好的`,
	}
	l, err := ling.DefaultNLP()
	if err != nil {
		t.Error(err)
	}
	tagger, err := ling.NewDictTagger("dict")
	if err != nil {
		t.Error(err)
	}
	if err = l.AddTagger(tagger); err != nil {
		t.Error(err)
	}
	for _, c := range tests {
		_, err = localGrammar(c, l)
		if err != nil {
			t.Error(err)
		}
	}
}
