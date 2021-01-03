package main

import (
	"RocketDiscord/cache"
	"context"
	"database/sql"
	"errors"
	"github.com/andersfylling/disgord"
	"github.com/sirupsen/logrus"
	_ "github.com/ziutek/mymysql/godrv"
	"gopkg.in/gorp.v2"
	"math/rand"
	"os"
	"time"
)

var dbMap *gorp.DbMap

var Log = &logrus.Logger{
	Out:       os.Stderr,
	Formatter: new(logrus.TextFormatter),
	Hooks:     make(logrus.LevelHooks),
	Level:     logrus.InfoLevel,
}

var botSettings = disgord.Config{
	BotToken:                     "xxx",
	Intents:                      disgord.AllIntents(),
	ProjectName:                  "RocketDiscord",
	CancelRequestWhenRateLimited: false,
	LoadMembersQuietly:           true,
	Logger:                       Log,
	DisableCache:                 false,
	Cache:                        &disgord.CacheNop{},
}

func main()  {
	Log.Info("Init database connection...")
	mysql, err := sql.Open("mymysql", "tcp:127.0.0.1:3306*rocketdiscord/rocketdiscorduser/Yrhzmudg")
	if err != nil {
		panic(err)
	}
	dbMap = &gorp.DbMap{
		Db:              mysql,
		Dialect:         gorp.MySQLDialect{"InnoDB", "UTF8"},
	}
	Log.Info("Database connected!")

	rand.Seed(time.Now().UTC().UnixNano())
	Log.Info("Set random time base")

	Log.Info("Bot starting...")

	RegisterModLogModule()

	botClient := disgord.New(botSettings)

	err = handle(botClient)
	if err != nil {
		Log.Panic("Handler error", err)
	}
}

//////////////////////////////////////////////////////
//
// BOT HANDLER
//
//////////////////////////////////////////////////////
func handle(client *disgord.Client) error {
	deadline, _ := context.WithTimeout(context.Background(), 5*time.Second)
	mdlw, err := NewMiddlewareHolder(client, deadline)
	if err != nil {
		panic(err)
	}

	//Register all listeners
	client.Gateway().WithMiddleware(mdlw.filterOutHumans, mdlw.filterOutOthersMsgs)

	client.Gateway().MessageCreate(CommandMessageCreate, PleasePornGif, MessageXD)
	client.Gateway().GuildMemberRemove(MemberRemoveGuildEvent)
	client.Gateway().GuildMemberAdd(MemberAddGuildEvent)
	client.Gateway().VoiceStateUpdate(cache.VoiceStateUpdate)
	client.Gateway().Ready(ReadyEvent)

	client.Gateway().MessageCreate(LogMessage)
	client.Gateway().MessageUpdate(LogMessageUpdate)
	client.Gateway().UserUpdate(LogUserUpdate)


	// connect now, and disconnect on system interrupt
	Log.Info("Connection estimated!")
	Log.Info("Bot running. Press CTRL-C to exit.")
	err = client.Gateway().StayConnectedUntilInterrupted()
	return err
}

//////////////////////////////////////////////////////
//
// MIDDLEWARES
//
//////////////////////////////////////////////////////
func NewMiddlewareHolder(s disgord.Session, ctx context.Context) (m *MiddlewareHolder, err error) {
	m = &MiddlewareHolder{session: s}
	if m.myself, err = s.CurrentUser().WithContext(ctx).Get(); err != nil {
		return nil, errors.New("unable to fetch info about the bot instance")
	}

	return m, nil
}

// instead of storing the instances in global variables. Middlewares can ofcourse be standalone functions too.
type MiddlewareHolder struct {
	session disgord.Session
	myself  *disgord.User
}

func (m *MiddlewareHolder) filterOutHumans(evt interface{}) interface{} {
	if e, ok := evt.(*disgord.MessageCreate); ok {
		// ignore humans
		if !e.Message.Author.Bot {
			return nil
		}
	}

	return evt
}

func (m *MiddlewareHolder) filterOutOthersMsgs(evt interface{}) interface{} {
	// ignore other bots
	// remove this check if you want to delete all bot messages after N seconds
	if e, ok := evt.(*disgord.MessageCreate); ok {
		if e.Message.Author.ID != m.myself.ID {
			return nil
		}
	}

	return evt
}
