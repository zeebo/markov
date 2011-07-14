package markov

import (
	"strings"
	"rand"
	"time"
)

func choice(w []string) string {
	return w[rand.Intn(len(w))]
}

func hash(s1, s2 string) string {
	return s1 + " " + s2
}

func secondWord(words string) string {
	return strings.Split(words, " ", 1)[1]
}

func grabKey(list []string) string {
	return hash(list[len(list) - 2], list[len(list) - 1])
}

type Markov struct {
	data map[string][]string
}

func New() *Markov {
	return &Markov{
		data: make(map[string][]string),
	}
}

func (m *Markov) Analyze(corpus string) {
	words := strings.Split(corpus, " ", -1)
	if len(words) == 1 {
		return
	}
	for i, first := range words[:len(words) - 2] {
		m.Add(first, words[i+1], words[i+2])
	}
	m.Add(words[len(words) - 2], words[len(words) - 1], "")
}

func (m *Markov) Add(w1, w2, word string) {
	key := hash(w1, w2)
	item, exists := m.data[key]
	if !exists {
		m.data[key] = make([]string, 0)
		item = m.data[key]
	}
	m.data[key] = append(item, word)
}

func (m *Markov) Generate() string {
	//generate a seed from the keys
	keys, i := make([]string, len(m.data)), 0

	for key := range m.data {
		keys[i] = key
		i++
	}
	
	seed := choice(keys)
	generated := strings.Split(seed, " ", -1)
	for {
		key := grabKey(generated)
		words, exists := m.data[key]
		if !exists {
			break
		}

		word := choice(words)
		if word == "" {
			break
		}

		generated = append(generated, word)
	}

	return strings.Join(generated, " ")
}

func init() {
	rand.Seed(time.Nanoseconds())
}