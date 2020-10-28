package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
	"time"
)

var (
	availableCommands = make(map[string]Command)
)

type Command struct {
	Name                 string
	HelpMessage          string

	GuildOwnerOnly       bool
	RequirePermissions   bool
	PermissionsLevel     int

	RequireArguments     bool
	MinimumArgumentsSize int
	MaximumArgumentsSize int

	NSFW                 bool

	Exec                 func(*discordgo.Session, *discordgo.MessageCreate, *discordgo.Guild, []string) error
}

func newCommand(name string, requirePermissions bool, nsfw bool, f func(*discordgo.Session, *discordgo.MessageCreate, *discordgo.Guild, []string) error) Command {
	if f == nil {
		panic("function is nil")
	}

	return Command{
		Name: name,
		RequirePermissions: requirePermissions,
		NSFW: nsfw,
		Exec: f,
	}
}

func (command *Command) register() Command {
	availableCommands[command.Name] = *command

	fmt.Printf("registered new command -> %s \n", command.Name)

	return *command
}

func ParseMessage(session *discordgo.Session, event *discordgo.MessageCreate, guild *discordgo.Guild, message string) {
	//Check if message contains command prefix
	if !strings.HasPrefix(message, "!") {
		return
	}

	//Parse to arguments
	args := strings.Fields(message)

	//Return if no args provided
	if len(args) <= 0 {
		return
	}

	//Parse command name
	commandName := strings.Replace(args[0], "!", "", 1)

	//Now we check if everything is clear to lift off
	command, ok := availableCommands[commandName]
	if !ok {
		return
	}

	//Load channel to memory
	channel, err := session.Channel(event.ChannelID)
	if err != nil || channel == nil {
		_ = fmt.Errorf("could not get channel for id %s, %s", event.GuildID, err)
	}

	//Check guild member permissions
	//TODO perms

	//Arguments check
	//TODO args check

	//Check if command is executed on nsfw channel
	if !channel.NSFW && command.NSFW {
		_, err := session.ChannelMessageSend(event.ChannelID, ":no_entry: Polecienie __**MOŻE**__ być tylko wykonane na kanale __**NSFW**__!")
		if err != nil {
			_ = fmt.Errorf("an error ocurrent while performing nfsw %s command warning: %s", commandName, err)
		}
		return
	}

	//Logs
	fmt.Printf("Executing user command. Command name: %s, user: %s, guild: %s, channel: %S, timestamp: %s \n", commandName, event.Author.Username, guild.Name, event.GuildID, time.Now().String())

	//Execute
	err = command.Exec(session, event, guild, args[1:])
	if err != nil {
		_ = fmt.Errorf("an error ocurrent while performing %s command: %s", commandName, err)
	}
}
