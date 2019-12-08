package main

import (
	"bufio"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/ngric/goborg/markov"
)

var (
	Chain    markov.Chain
	token    string
	lastChan string
)

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()
	load()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	sarr := strings.Fields(m.Content)

	if len(sarr) > 0 {
		for i, v := range sarr {
			if i == len(sarr)-1 {
				go Chain.AddEdge(v, "")
			} else {
				go Chain.AddEdge(v, sarr[i+1])
			}
		}
		reply := Chain.GetLine(sarr[0])
		lastChan = m.ChannelID
		s.ChannelMessageSend(m.ChannelID, reply)
	}

}

func main() {
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	discord.AddHandler(messageCreate)

	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.ChannelMessageSend(lastChan, "Byebye")
	discord.Close()
	save()
}

func load() {
	f, err := os.Open("brain")
	defer f.Close()

	if err != nil {
		fmt.Println("Unable to open brain. Creating new...")
		Chain = markov.NewChain()
		return
	}

	r := bufio.NewReader(f)
	dec := gob.NewDecoder(r)
	err = dec.Decode(&Chain)

	if err != nil {
		log.Fatal("Error while reading brain, ", err)
	}
}

func save() {
	f, err := os.Create("brain")
	defer f.Close()

	if err != nil {
		fmt.Println("Unable to open brain for saving.")
		return
	}

	w := bufio.NewWriter(f)
	enc := gob.NewEncoder(w)
	err = enc.Encode(&Chain)

	if err != nil {
		log.Fatal("Error while writing brain, ", err)
	}

	w.Flush()
}
