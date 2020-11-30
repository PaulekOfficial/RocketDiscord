package main

import (
	"database/sql"
	_ "github.com/ziutek/mymysql/godrv"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"gopkg.in/gorp.v2"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var dbMap *gorp.DbMap

func main()  {
	fmt.Println("Init database connection...")
	mysql, err := sql.Open("mymysql", "xxx")
	if err != nil {
		panic(err)
	}
	dbMap = &gorp.DbMap{
		Db:              mysql,
		Dialect:         gorp.MySQLDialect{"InnoDB", "UTF8"},
	}
	fmt.Println("Database connected!")

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

	RegisterUpdateListeners(discordSession)

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
