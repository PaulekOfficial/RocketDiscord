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

func (command *Command) setPermissions(requirePermissions bool, guildOwnerOnly bool, permissionsLevel int) {
	command.GuildOwnerOnly = guildOwnerOnly
	command.RequirePermissions = requirePermissions
	command.PermissionsLevel = permissionsLevel
}

//If MinimumArgumentsSize is set to -1 that means no limit
//If MinimumArgumentsSize is set to -1 that means no limit
func (command *Command) setArgumentsRequirements(requireArguments bool, minimumArgumentsSize int, maximumArgumentsSize int) {
	command.RequireArguments = requireArguments
	command.MinimumArgumentsSize = minimumArgumentsSize
	command.MaximumArgumentsSize = maximumArgumentsSize
}

func (command *Command) setHelpMessage(message string) {
	command.HelpMessage = message
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
	permissionLevel, err := channelPermissionLevel(session, event.ChannelID, event.Author.ID)
	if err != nil {
		_ = fmt.Errorf("could not get user permissionlevel for user %s, %s", event.Author.Username, err)
		return
	}
	if (command.GuildOwnerOnly && event.Author.ID != guild.OwnerID) || (command.RequirePermissions && permissionLevel&command.PermissionsLevel <= 0) {
		_, err := session.ChannelMessageSend(event.ChannelID, fmt.Sprintf(":rotating_light: Brak wystarczających uprawnień do użycia tego polecenia! Wymagany poziom dostępu %d.", command.PermissionsLevel))
		if err != nil {
			_ = fmt.Errorf("an error ocurrent while performing no permissions %s command warning: %s", commandName, err)
		}
		return
	}

	//Arguments check
	length := len(args[1:])
	if (length > command.MaximumArgumentsSize && command.MaximumArgumentsSize != -1) || (length < command.MinimumArgumentsSize && command.MinimumArgumentsSize != -1) {
		_, err := session.ChannelMessageSend(event.ChannelID, fmt.Sprintf(":boom: Niepoprawne użycie polecenia %s.", commandName))
		if err != nil {
			_ = fmt.Errorf("an error ocurrent while performing bad usage %s command warning: %s", commandName, err)
		}
		_, err = session.ChannelMessageSend(event.ChannelID, fmt.Sprintf(":star: Wyświetlam pomoc dla polecenia %s: %s", commandName, command.HelpMessage))
		if err != nil {
			_ = fmt.Errorf("an error ocurrent while performing help usage %s command warning: %s", commandName, err)
		}
		return
	}

	//Check if command is executed on nsfw channel
	if !channel.NSFW && command.NSFW {
		_, err := session.ChannelMessageSend(event.ChannelID, ":no_entry: Polecienie __**MOŻE**__ być tylko wykonane na kanale __**NSFW**__!")
		if err != nil {
			_ = fmt.Errorf("an error ocurrent while performing nfsw %s command warning: %s", commandName, err)
		}
		return
	}

	//Logs
	fmt.Printf("Executing user command. Command name: %s, user: %s, guild: %s, channel: %v, timestamp: %s \n", commandName, event.Author.Username, guild.Name, event.GuildID, time.Now().String())

	//Execute
	err = command.Exec(session, event, guild, args[1:])
	if err != nil {
		_ = fmt.Errorf("an error ocurrent while performing %s command: %s", commandName, err)
	}
}

func channelPermissionLevel(session *discordgo.Session, channelID string, memberID string) (permissionLevel int, err error) {
	permissionLevel, err = session.UserChannelPermissions(memberID, channelID)
	return
}
