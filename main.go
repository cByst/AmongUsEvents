package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/cbyst/AmongUsHelper/amongushandlers"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Set up run time flags
	botToken := flag.String("token", "", "Discord bot token.")
	logDebug := flag.Bool("debug", false, "Set logging level to debug.")

	flag.Parse()

	// Setup debug logging
	if *logDebug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	// Check for required token
	if *botToken == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	session, err := discordgo.New(fmt.Sprintf("Bot %s", *botToken))
	if err != nil {
		log.Fatal(errors.WithMessage(err, "Error starting and connecting discord bot session"))
		os.Exit(1)
	}

	defer session.Close()

	amongushandlers.AttachHandlers(session)

	err = session.Open()
	if err != nil {
		log.Fatal(errors.WithMessage(err, "Error opening Discord session"))
		os.Exit(1)
	}

	fmt.Println("Among Us Helper is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

}
