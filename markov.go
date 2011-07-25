package markov

import (
    "strings"
    "rand"
    "time"
    "gob"
    "os"
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
    return hash(list[len(list)-2], list[len(list)-1])
}

type Markov struct {
    data  map[string][]string
    seeds []string
}

func New() *Markov {
    return &Markov{
        data:  make(map[string][]string),
        seeds: make([]string, 0),
    }
}

func (m *Markov) Analyze(corpus string) {
    words := strings.Split(corpus, " ", -1)
    if len(words) == 1 {
        return
    }
    m.AddSeed(words[0], words[1])
    for i, first := range words[:len(words)-2] {
        m.Add(first, words[i+1], words[i+2])
    }
    m.Add(words[len(words)-2], words[len(words)-1], "")
}

func (m *Markov) AddSeed(w1, w2 string) {
    m.seeds = append(m.seeds, hash(w1, w2))
}

func (m *Markov) RandSeed() string {
    return choice(m.seeds)
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
    return m.GenerateFrom(m.RandSeed())
}

func (m *Markov) GenerateFrom(seed string) string {
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

func (m *Markov) Save(filename string) os.Error {
    hnd, err := os.Create(filename)
    defer hnd.Close()
    if err != nil {
        return err
    }

    encoder := gob.NewEncoder(hnd)
    if err := encoder.Encode(m.data); err != nil {
        return err
    }

    if err := encoder.Encode(m.seeds); err != nil {
        return err
    }

    return nil
}

func (m *Markov) Load(filename string) os.Error {
    hnd, err := os.Open(filename)
    defer hnd.Close()
    if err != nil {
        return err
    }

    decoder := gob.NewDecoder(hnd)
    if err := decoder.Decode(&m.data); err != nil {
        return err
    }

    if err := decoder.Decode(&m.seeds); err != nil {
        return err
    }

    return nil
}

func init() {
    rand.Seed(time.Nanoseconds())
}
