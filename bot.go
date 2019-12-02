package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func main() {
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
