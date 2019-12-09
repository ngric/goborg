/*
* Assignment:	Final Project, Part #3
* Author:		Nile Grice (nile@email.arizona.edu)
*
* Course:		CSC372
* Instructor:	L. McCann
* TA:			Tito Ferra
* Due Date:		December 9, 2019
*
* Description:	A markov-chain based discord chatbot
*		Usage: goborg -t <discord api token>
*
* History URL:	https://ngric.github.io/goborg/
 */

package main

import (
	"bufio"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"math/rand"
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
	lastChan string  // store this so that we can give a parting message
	rate     float32 // probability that any given message will get a reply
)

// runs at program launch. reads in token from arguments and loads
// the bot's "brain" from disk
func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()
	load()
	rate = .8
}

func main() {
	// instantiate a discord session with passed token
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// register messageCreate method as handler for message events
	discord.AddHandler(messageCreate)

	// open connection to discord
	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// continue to run until we receive an interrupt of some variety
	// eg: Ctrl-C
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.ChannelMessageSend(lastChan, "Byebye")
	discord.Close()
	save()
}

// loads a Chain struct (the bot's 'brain') from disk.
// expects to find it in a file named 'brain', if it can't find or open
// such a file, the bot will start from empty
func load() {
	f, err := os.Open("brain") // attempt to open file
	defer f.Close()            // close file when this method returns

	if err != nil {
		fmt.Println("Unable to open brain. Creating new...")
		Chain = markov.NewChain()
		return
	}

	// read opened file from disk
	r := bufio.NewReader(f)  // file stream
	dec := gob.NewDecoder(r) // stream decoder
	err = dec.Decode(&Chain)

	if err != nil {
		log.Fatal("Error while reading brain, ", err)
	}
}

// saves a Chain struct (the bot's 'brain') to disk.
// will overwrite the file if it already exists
func save() {
	f, err := os.Create("brain") // truncates the file if it already exists
	defer f.Close()              // make sure the file is closed when this method returns

	if err != nil {
		fmt.Println("Unable to open brain for saving.")
		return
	}

	// write Chain to disk
	w := bufio.NewWriter(f)  // file output stream
	enc := gob.NewEncoder(w) // encodes whatever it's given, then passes it to w
	err = enc.Encode(&Chain)

	if err != nil {
		log.Fatal("Error while writing brain, ", err)
	}

	w.Flush() // flush output buffers to disk
}

// method that runs every time a message is received from discord
// trains the markov chain on the contents of the message, and potentially
// gets/sends a reply
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore messages that originate from goborg itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// break message into string array delimited on spaces
	sarr := strings.Fields(m.Content)

	if len(sarr) > 0 { // messages should always be non-empty, but just in case
		for i, v := range sarr {
			// last word in message has an edge to "", which
			// we're using for termination
			if i == len(sarr)-1 {
				go Chain.AddEdge(v, "")
			} else { // add edge between current and following words
				go Chain.AddEdge(v, sarr[i+1])
			}
		}

		// roll for a reply
		if rand.Float32() < rate {
			reply := Chain.GetLine(sarr[0])
			s.ChannelMessageSend(m.ChannelID, reply)
		}

		lastChan = m.ChannelID // store for "goodbye" message on quit
	}

}
