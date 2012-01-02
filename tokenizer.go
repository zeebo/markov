package markov

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

type stateFn func(t *Tokenizer, to *Token) (stateFn, error)

type Tokenizer struct {
	r         *bufio.Reader
	state     stateFn
	startEmit int
	order     int
	readerErr bool
	sentence  []byte
	p         int
}

func NewTokenizer(r io.Reader, order int) *Tokenizer {
	return &Tokenizer{
		state: tokenizerStartOfSentence,
		r:     bufio.NewReader(r),
		order: order,
	}
}

func NewTokenizerString(data string, order int) *Tokenizer {
	return &Tokenizer{
		state: tokenizerStartOfSentence,
		r:     bufio.NewReader(bytes.NewBufferString(data)),
		order: order,
	}
}

func tokenizerEof(t *Tokenizer, to *Token) (stateFn, error) {
	return tokenizerEof, io.EOF
}

func tokenizerStartOfSentence(t *Tokenizer, to *Token) (stateFn, error) {
	t.startEmit++
	to.Value = ""
	to.Type = tokenSentenceStartType

	if t.startEmit == t.order {
		t.startEmit = 0
		return tokenizerGrabSentence, nil
	}
	return tokenizerStartOfSentence, nil
}

func tokenizerEndOfSentence(t *Tokenizer, to *Token) (stateFn, error) {
	to.Value = ""
	to.Type = tokenSentenceEndType
	if t.readerErr {
		return tokenizerEof, nil
	}
	return tokenizerStartOfSentence, nil
}

func tokenizerGrabSentence(t *Tokenizer, to *Token) (stateFn, error) {
	var err error
	t.sentence, err = t.r.ReadSlice('\n')
	if err != nil {
		t.readerErr = true
	}
	t.p = 0
	return tokenizerHasSentence(t, to)
}

func tokenizerHasSentence(t *Tokenizer, to *Token) (stateFn, error) {
	next := bytes.IndexByte(t.sentence[t.p:], ' ')
	if next == -1 {
		to.Value = strings.TrimSpace(string(t.sentence[t.p:]))
		to.Type = tokenWordType
		return tokenizerEndOfSentence, nil
	}
	var prev int
	prev, t.p = t.p, t.p+next+1
	to.Value = string(t.sentence[prev : prev+next])
	to.Type = tokenWordType
	return tokenizerHasSentence, nil
}

func (t *Tokenizer) Next(to *Token) error {
	var err error
	t.state, err = t.state(t, to)
	return err
}

func (t *Tokenizer) Sentence() (Sentence, error) {
	var (
		err  error
		next Token
	)
	sentence := make([]Token, 0, t.order+2) //begin + 1 word + end
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
