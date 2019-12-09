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
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ngric/goborg/markov"
)

var (
	bot      markov.Chain
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

	go func() { // start go routine that periodically calls save
		delay, _ := time.ParseDuration("10m")
		for { // loop forever
			time.Sleep(delay)
			bot.Save()
		}
	}()

	go consoleInput()

	// continue to run until we receive an interrupt of some variety
	// eg: Ctrl-C
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.ChannelMessageSend(lastChan, "Byebye")
	discord.Close()
	bot.Save()
}

// loads a Chain struct (the bot's 'brain') from disk.
// expects to find it in a file named 'brain', if it can't find or open
// such a file, the bot will start from empty
func load() {
	fmt.Println("Loading......")

	f, err := os.Open("brain") // attempt to open file
	defer f.Close()            // close file when this method returns

	if err != nil {
		fmt.Println("Unable to open brain. Creating new...")
		bot = markov.NewChain()
		return
	}

	// read opened file from disk
	r := bufio.NewReader(f)  // file stream
	dec := gob.NewDecoder(r) // stream decoder
	err = dec.Decode(&bot)

	if err != nil {
		log.Fatal("Error while reading brain, ", err)
	}
}

// method that runs every time a message is received from discord
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore messages that originate from goborg itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// print the message to console
	fmt.Printf("MSG: %s\n", m.Content)

	// have the message placed into the chain
	reply := parseMessage(m.Content)
	if reply != "" { // send reply if we got one
		s.ChannelMessageSend(m.ChannelID, reply)
		fmt.Printf("REPLYING: %s\n", reply)
	}

	lastChan = m.ChannelID // store for "goodbye" message on quit
}

// takes user input from console
func consoleInput() {
	r := bufio.NewReader(os.Stdin) // buffered stdin reader

	for { // endlessly loop, taking input from stdin
		in, _ := r.ReadString('\n')
		reply := parseMessage(in)
		if reply != "" {
			fmt.Printf("CONSOLE REPLY: %s\n", reply)
		}
	}
}

// trains the markov chain on the contents of the message, and potentially
// gets a reply. If no reply is generated, returns empty staing
func parseMessage(msg string) string {
	// break message into string array delimited on spaces
	sarr := strings.Fields(msg)

	if len(sarr) > 0 { // messages should always be non-empty, but just in case
		for i, v := range sarr {
			// last word in message has an edge to "", which
			// we're using for termination
			if i == len(sarr)-1 {
				go bot.AddEdge(v, "")
			} else { // add edge between current and following words
				go bot.AddEdge(v, sarr[i+1])
			}
		}

		// roll for a reply
		if rand.Float32() < rate {
			n := rand.Intn(len(sarr)) // choose reply seed at random
			reply := bot.GetLine(sarr[n])
			return reply
		}
	}

	return ""
}
