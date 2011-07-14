package main

import (
	"markov"
	"io/ioutil"
	"fmt"
	"strings"
)

func load(m *markov.Markov, file string) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	chunks := strings.Split(string(data), "\n", -1)
	for _, line := range chunks {
		m.Analyze(line)
	}
}

func main() {
	m := markov.New()
	load(m, "jeeves2.txt")
	load(m, "hamlet2.txt")
	fmt.Println(m.Generate())
}