package bnf

import (
	"fmt"
	"strings"

	"github.com/liuzl/ling"
)

func localGrammar(text string, lnlp *ling.Pipeline) (*Grammar, error) {
	if text = strings.TrimSpace(text); text == "" {
		return nil, fmt.Errorf("text is empty")
	}
	if lnlp == nil {
		lnlp = nlp
	}
	d := ling.NewDocument(text)
	if err := lnlp.Annotate(d); err != nil {
		return nil, err
	}
	for _, span := range d.Spans {
		fmt.Println(span, span.Annotations)
	}
	return nil, nil
}
