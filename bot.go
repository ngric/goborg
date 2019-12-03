package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ngric/goborg/markov"
)

func main() {
	chain := markov.NewChain()
	reader := bufio.NewReader(os.Stdin)
	for {
		s, _ := reader.ReadString('\n')
		sarr := strings.Fields(s)

		if len(sarr) > 0 {
			for i, v := range sarr {
				if i == len(sarr)-1 {
					chain.AddEdge(v, "")
				} else {
					chain.AddEdge(v, sarr[i+1])
				}
			}
			fmt.Printf("***Bot: %s\n", chain.GetLine(sarr[0]))
		}

	}
}

func dg() {
	discord, err := discordgo.New("Bot" + "authentication token")
	if err != nil {
		panic(err)
	}
	err = discord.Close()
	if err != nil {
		panic(err)
	}
	fmt.Println("did nothing!")
}
