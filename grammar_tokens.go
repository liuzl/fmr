package fmr

import (
	"sync"

	"github.com/liuzl/ling"
)

type cMap struct {
	tokens map[string]*ling.Token
	sync.RWMutex
}

func (m *cMap) get(k string) *ling.Token {
	m.RLock()
	defer m.RUnlock()
	return m.tokens[k]
}

func (m *cMap) put(k string, token *ling.Token) {
	m.Lock()
	defer m.Unlock()
	m.tokens[k] = token
}

var gTokens = &cMap{tokens: make(map[string]*ling.Token)}
