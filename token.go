package markov

import (
	"io"
	"bufio"
	"os"
	"bytes"
	"fmt"
	"strings"
)

type Tokenizer struct {
	r           *bufio.Reader
	state       tokenizerState
	readerState readerState
	sentence    []byte
	p           int
}

type tokenizerState int

const (
	tokenizerStartOfSentence = tokenizerState(iota)
	tokenizerAfterStart
	tokenizerGrabSentence
	tokenizerHasSentence
	tokenizerEndOfSentence
	tokenizerEof
)

func (t tokenizerState) String() string {
	switch t {
	case tokenizerStartOfSentence:
		return "StartOfSentence"
	case tokenizerAfterStart:
		return "AfterStart"
	case tokenizerGrabSentence:
		return "GrabSentence"
	case tokenizerHasSentence:
		return "HasSentence"
	case tokenizerEndOfSentence:
		return "EndOfSentence"
	case tokenizerEof:
		return "EOF"
	}
	return "Unknown State"
}

type readerState int

const (
	readerHasData = readerState(iota)
	readerEof
)

func (rs readerState) String() string {
	switch rs {
	case readerHasData:
		return "HasData"
	case readerEof:
		return "EOF"
	}
	return "Unknown State"
}

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

func NewTokenizer(r io.Reader) *Tokenizer {
	return &Tokenizer{r: bufio.NewReader(r)}
}

func (t *Tokenizer) Next(to *Token) os.Error {
	switch t.state {
	case tokenizerEof:
		return os.EOF

	case tokenizerStartOfSentence:
		t.state = tokenizerAfterStart
		to.Value = ""
		to.Type = tokenSentenceStartType

	case tokenizerAfterStart:
		t.state = tokenizerGrabSentence
		to.Value = ""
		to.Type = tokenSentenceStartType

	case tokenizerEndOfSentence:
		if t.readerState == readerEof {
			t.state = tokenizerEof
		} else {
			t.state = tokenizerStartOfSentence
		}
		to.Value = ""
		to.Type = tokenSentenceEndType

	case tokenizerGrabSentence:
		var err os.Error
		t.sentence, err = t.r.ReadSlice('\n')
		if err != nil {
			t.readerState = readerEof
		}
		t.p = 0
		t.state = tokenizerHasSentence
		return t.Next(to)

	case tokenizerHasSentence:
		next := bytes.IndexByte(t.sentence[t.p:], ' ')
		if next == -1 {
			t.state = tokenizerEndOfSentence
			to.Value = strings.TrimSpace(string(t.sentence[t.p:]))
			to.Type = tokenWordType
			return nil
		}
		var prev int
		prev, t.p = t.p, t.p+next+1
		to.Value = string(t.sentence[prev : prev+next])
		to.Type = tokenWordType
	}

	return nil
}

func (t *Tokenizer) Sentence() (Sentence, os.Error) {
	var (
		err  os.Error
		next Token
	)
	sentence := make([]Token, 0, 3)
	for {
		err = t.Next(&next)
		if err != nil {
			break
		}
		sentence = append(sentence, next.Dup())
		if next.Type == tokenSentenceEndType {
			break
		}
	}

	return Sentence(sentence), err
}
