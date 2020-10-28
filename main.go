package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
)

func main()  {
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
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	err = discordSession.Close()
	if err != nil {
		_ = fmt.Errorf("cannot close bot connection %s", err)
	}
}
