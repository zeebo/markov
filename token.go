package markov

import (
	"fmt"
	"strings"
)

type tokenType int

const (
	tokenSentenceStartType = tokenType(iota)
	tokenSentenceEndType
	tokenWordType
)

func (tt tokenType) String() string {
	switch tt {
	case tokenSentenceStartType:
		return "SentenceStartType"
	case tokenSentenceEndType:
		return "SentenceEndType"
	case tokenWordType:
		return "WordType"
	}
	return "Unknown Type"
}

type Token struct {
	Value string
	Type  tokenType
}

func (t Token) Equals(ot Token) bool {
	return t.Value == ot.Value && t.Type == ot.Type
}

func (t Token) String() string {
	switch t.Type {
	case tokenWordType:
		return fmt.Sprintf("[%s] %q", t.Type, t.Value)
	}
	return fmt.Sprintf("[%s]", t.Type)
}

func (t *Token) Dup() Token {
	return Token{t.Value, t.Type}
}

type Sentence []Token

func (s Sentence) String() string {
	r := make([]string, len(s))
	for i := range s {
		r[i] = s[i].Value
	}
	return strings.TrimSpace(strings.Join(r, " "))
}

func (s Sentence) Hash() string {
	r := make([]string, len(s))
	for i := range s {
		r[i] = s[i].String()
	}
	return strings.Join(r, "|")
}
