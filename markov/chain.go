package markov

import (
	"math/rand"
	"sync"
)

type Chain struct {
	mut   sync.RWMutex
	Words map[string]Word
}

type Word struct {
	Edges map[string]int
	Links int
}

func newWord() Word {
	var w Word
	w.Edges = make(map[string]int)
	return w
}

func NewChain() Chain {
	var c Chain
	c.Words = make(map[string]Word)
	return c
}

func (c *Chain) AddEdge(from, to string) {
	c.mut.Lock()
	defer c.mut.Unlock()

	w, ok := (*c).Words[from]

	if ok { // if the starting word exists, add edge leading to dest word
		_, ok = w.Edges[to]
		if ok { // if a matching edge already exists, increment its weight
			w.Edges[to]++
		} else { // make the edge if it's new
			w.Edges[to] = 1
		}
		w.Links++
		return
	} else { // starting word is new, add to chain
		w = newWord()
		w.Edges[to] = 1
		w.Links++
		(*c).Words[from] = w
		return
	}
}

func (c *Chain) nextFrom(s string) string {
	c.mut.RLock()
	defer c.mut.RUnlock()
	w := (*c).Words[s]

	if l := w.Links; l > 0 {
		n := rand.Intn(l)
		sum := 0
		for k, v := range w.Edges {
			sum += v
			if sum >= n {
				return k
			}
		}
	}
	return ""
}

func (c *Chain) GetLine(start string) string {
	s, w := "", c.nextFrom(start)
	for w != "" {
		s += w + " "
		w = c.nextFrom(w)
	}
	return s
}
