package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/eaok/ashe/config"
	"github.com/eaok/ashe/router"
)

func init() {
	// Read config from configuration
	config.ReadConfig()
}

func main() {
	// Create a new Discord session using the provided bot token.
	discord, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println("Error creating discord session.")
		return
	}

	router.Route(discord)

	// need to open the socket
	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()
}
