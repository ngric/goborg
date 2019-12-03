package markov

import (
	"math/rand"
)

type Chain struct {
	words map[string]word
}

type word struct {
	edges map[string]int
	links int
}

func newWord() word {
	var w word
	w.edges = make(map[string]int)
	return w
}

func NewChain() Chain {
	var c Chain
	c.words = make(map[string]word)
	return c
}

func (c *Chain) AddEdge(from, to string) {
	w, ok := (*c).words[from]

	if ok { // if the starting word exists, add edge leading to dest word
		_, ok = w.edges[to]
		if ok { // if a matching edge already exists, increment its weight
			w.edges[to]++
		} else { // make the edge if it's new
			w.edges[to] = 1
		}
		w.links++
		return
	} else { // starting word is new, add to chain
		w = newWord()
		w.edges[to] = 1
		w.links++
		(*c).words[from] = w
		return
	}
}

func (c *Chain) nextFrom(s string) string {
	w := (*c).words[s]

	if l := w.links; l > 0 {
		n := rand.Intn(l)
		sum := 0
		for k, v := range w.edges {
			sum += v
			if sum >= n {
				return k
			}
		}
	}
	return ""
}

func (c *Chain) GetLine(start string) string {
	s := ""
	for w := c.nextFrom(start); w != ""; w = c.nextFrom(w) {
		s += w + " "
	}
	return s
}
