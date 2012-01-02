package markov

import (
	"encoding/gob"
	"io"
	"math"
	"math/rand"
)

type Markov struct {
	order int
	data  map[string][]Token
}

func NewMarkov(order int) *Markov {
	return &Markov{
		order: order,
		data:  make(map[string][]Token),
	}
}

func NewMarkovFrom(r io.Reader) (*Markov, error) {
	dec := gob.NewDecoder(r)
	m := &Markov{}

	if err := dec.Decode(&m.order); err != nil {
		return nil, err
	}
	if err := dec.Decode(&m.data); err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Markov) Save(w io.Writer) error {
	enc := gob.NewEncoder(w)

	if err := enc.Encode(m.order); err != nil {
		return err
	}
	if err := enc.Encode(m.data); err != nil {
		return err
	}

	return nil
}

func (m *Markov) next(prefix Sentence) Token {
	key := prefix[len(prefix)-m.order:].Hash()
	if _, ex := m.data[key]; !ex {
		return Token{"", tokenSentenceEndType}
	}
	//choose a random token
	return m.data[key][rand.Intn(len(m.data[key]))]
}

func (m *Markov) score(s Sentence) float64 {
	var score float64
	for i := 0; i < len(s)-m.order; i++ {
		score += m.scoreIndivdual(s[i:i+m.order], s[i+m.order])
	}
	score /= float64(len(s) + 7) //add 7 to penalize short sentences
	return score
}

func (m *Markov) scoreIndivdual(prefix Sentence, suffix Token) float64 {
	var num, denom float64
	key := prefix.Hash()

	if _, ex := m.data[key]; !ex {
		return 0
	}

	denom = float64(len(m.data[key]))
	for _, v := range m.data[key] {
		if suffix.Equals(v) {
			num++
		}
	}

	return -1 * math.Log2(num/denom)
}

func (m *Markov) Analyze(s Sentence) {
	for i := 0; i < len(s)-m.order; i++ {
		pre, post := s[i:i+m.order], s[i+m.order]
		m.insert(pre.Hash(), post)
	}
}

func (m *Markov) AnalyzeFully(t *Tokenizer) {
	for {
		sentence, err := t.Sentence()
		if err != nil {
			break
		}
		m.Analyze(sentence)
	}
}

func (m *Markov) insert(key string, val Token) {
	if _, ex := m.data[key]; !ex {
		m.data[key] = make([]Token, 0)
	}
	m.data[key] = append(m.data[key], val)
}

func (m *Markov) seed() Sentence {
	seed := make([]Token, m.order)
	for i := range seed {
		seed[i] = Token{Type: tokenSentenceStartType}
	}
	return seed
}

func (m *Markov) Generate() Sentence {
	seed := m.seed()
	for {
		next := m.next(seed)
		seed = append(seed, next)
		if next.Type == tokenSentenceEndType {
			break
		}
	}
	return seed
}

func (m *Markov) GenerateN(n int) Sentence {
	var (
		max       float64
		temp, val Sentence
	)
	for i := 0; i < n; i++ {
		temp = m.Generate()
		if s := m.score(temp); s >= max {
			max = s
			val = temp
		}
	}
	return val
}
