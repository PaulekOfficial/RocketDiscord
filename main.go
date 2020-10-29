package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main()  {
	rand.Seed(time.Now().UTC().UnixNano())
	fmt.Println("Set random time base")

	fmt.Println("Bot starting...")

	discordSession, err := discordgo.New("Bot " + "xxx")
	if err != nil {
		panic(err)
	}

	discordSession.AddHandler(ReadyEvent)
	discordSession.AddHandler(CommandMessageCreate)
	discordSession.AddHandler(MemberAddGuildEvent)
	discordSession.AddHandler(MemberRemoveGuildEvent)
	discordSession.AddHandler(MessageXD)
	discordSession.AddHandler(PleasePornGif)

	discordSession.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	err = discordSession.Open()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connection estimated!")

	fmt.Println("Bot running. Press CTRL-C to exit.")
	cc := make(chan os.Signal, 1)
	signal.Notify(cc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-cc

	err = discordSession.Close()
	if err != nil {
		_ = fmt.Errorf("cannot close bot connection %s", err)
	}
}
