package main

import (
	"markov"
	"fmt"
	"os"
)

func main() {
	r, err := os.Open("turk.txt")
	if err != nil {
		panic(err)
	}

	t := markov.NewTokenizer(r)
	for {
		sentence, err := t.Sentence()
		if err != nil {
			break
		}
		fmt.Println(len(sentence), sentence)
	}
}