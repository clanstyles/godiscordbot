package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"

	"github.com/LetsLearnCommunity/godiscordbot"
	"github.com/LetsLearnCommunity/godiscordbot/modules/music"
)

var (
	token = flag.String("token", "", "discord bot token")
)

func init() {
	flag.Parse()

	// Require the application to have a token
	if *token == "" {
		flag.Usage()
		os.Exit(0)
	}
}

func main() {
	session, err := discordgo.New(fmt.Sprintf("Bot %s", *token))
	if err != nil {
		log.Fatalf("failed to create a discord session: %s", err)
	}
	defer session.Close()

	godiscordbot.AddModule(music.New())
	godiscordbot.Initalize(session)

	if err := session.Open(); err != nil {
		log.Fatalf("failed top open websocket: %s", err)
	}

	log.Printf("GoDiscordBot is now running. Press CTRL-C to quit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
