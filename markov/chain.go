/*
* Assignment:	Final Project, Part #3
* Author:		Nile Grice (nile@email.arizona.edu)
*
* Course:		CSC372
* Instructor:	L. McCann
* TA:			Tito Ferra
* Due Date:		December 9, 2019
*
* Description:	A markov-chain implementation
*
* History URL:	https://ngric.github.io/goborg/
 */

package markov

import (
	"math/rand"
	"sync"
)

// Overall markov-chain struct. Just a string:Word map with a mutex
type Chain struct {
	mut   sync.RWMutex
	Words map[string]Word
}

// Word struct contains a string:int map representing edges from this word
// to other words, and thier weights. Links keeps track of the total weight
// of all edges out of this Word.
type Word struct {
	Edges map[string]int
	Links int
}

// instantiates a new Word
func newWord() Word {
	var w Word
	w.Edges = make(map[string]int)
	return w
}

// instantiates a new Chain
func NewChain() Chain {
	var c Chain
	c.Words = make(map[string]Word)
	return c
}

// adds a new directed edge between two words. If a matching edge already exists
// the weight of the edge is incremented by one.
func (c *Chain) AddEdge(from, to string) {
	// potentially excessive mutex. Invalidates any potintial Chain-related
	// concurrency as currently implemented
	c.mut.Lock()
	defer c.mut.Unlock()

	// see if "from" is a word that exists in the chain
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

// Uses the Chain to generate a "sentence."
// generation is seeded on the passed string
func (c *Chain) GetLine(start string) string {
	s, w := "", c.nextFrom(start)
	for w != "" { // continue building the string until we get an empty word
		s += w + " "
		w = c.nextFrom(w)
	}
	return s
}

// retreives a word adjacent from passed string in the Chain.
// retreived word is influenced by weight of the edge connecting it
func (c *Chain) nextFrom(s string) string {
	c.mut.RLock() // lock chain
	defer c.mut.RUnlock()
	w := (*c).Words[s] // assume s exists in the chain, get its neighbors

	if l := w.Links; l > 0 { // there should be no words with 0 links, but this
		// may change
		// choose next word by rolling a random number n, and summing weights
		// until they surpass n
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
