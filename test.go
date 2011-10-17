package main

import (
	"markov"
	"fmt"
	"os"
	"rand"
	"time"
)

func main() {
	rand.Seed(time.Nanoseconds())

	m := markov.NewMarkov(2)

	r, err := os.Open("turk.txt")
	if err != nil {
		panic(err)
	}
	t := markov.NewTokenizer(r, 2)
	m.AnalyzeFully(t)

	l := markov.NewLimiter(m, 100)

	for i := 0; i < 100; i++ {
		fmt.Println(l.Get())
		<-time.After(5e8)
	}
}