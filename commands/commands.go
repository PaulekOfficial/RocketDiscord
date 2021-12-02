package commands

import (
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

var (
	availableCommands = make(map[string]Command)
)

var Log = &logrus.Logger{
	Out:       os.Stderr,
	Formatter: new(logrus.TextFormatter),
	Hooks:     make(logrus.LevelHooks),
	Level:     logrus.InfoLevel,
}

type Command struct {
	Name                 string
	HelpMessage          string

	GuildOwnerOnly       bool
	RequirePermissions   bool
	PermissionsLevel     disgord.PermissionBit

	RequireArguments     bool
	MinimumArgumentsSize int
	MaximumArgumentsSize int

	NSFW                 bool

	Exec                 func(disgord.Session, *disgord.MessageCreate, *disgord.Guild, []string) error
}

type GuildCaller struct {
	session disgord.Session
	guildID disgord.Snowflake
}

func (caller *GuildCaller) Guild(id disgord.Snowflake) disgord.GuildQueryBuilder {
	return caller.session.Guild(id)
}

func NewCommand(name string, requirePermissions bool, nsfw bool, f func(disgord.Session, *disgord.MessageCreate, *disgord.Guild, []string) error) Command {
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

func (command *Command) SetPermissions(requirePermissions bool, guildOwnerOnly bool, permissionsLevel disgord.PermissionBit) {
	command.GuildOwnerOnly = guildOwnerOnly
	command.RequirePermissions = requirePermissions
	command.PermissionsLevel = permissionsLevel
}

//If MinimumArgumentsSize is set to -1 that means no limit
//If MinimumArgumentsSize is set to -1 that means no limit
func (command *Command) SetArgumentsRequirements(requireArguments bool, minimumArgumentsSize int, maximumArgumentsSize int) {
	command.RequireArguments = requireArguments
	command.MinimumArgumentsSize = minimumArgumentsSize
	command.MaximumArgumentsSize = maximumArgumentsSize
}

func (command *Command) SetHelpMessage(message string) {
	command.HelpMessage = message
}

func (command *Command) Register() Command {
	availableCommands[command.Name] = *command

	Log.WithFields(logrus.Fields{
		"command-name": command.Name,
	}).Info("Registered new command")

	return *command
}

func ParseMessage(session disgord.Session, event *disgord.MessageCreate, guild *disgord.Guild, message string) {
	guildSettings := cache.CreateOrGetGuildSettings(guild)
	commandPrefix := guildSettings.CommandPrefix

	//Check if message contains command prefix
	if !strings.HasPrefix(message, commandPrefix) {
		return
	}

	//Parse to arguments
	args := strings.Fields(message)

	//Return if no args provided
	if len(args) <= 0 {
		return
	}

	//Parse command name
	commandName := strings.Replace(args[0], commandPrefix, "", 1)

	//Now we check if everything is clear to lift off
	command, ok := availableCommands[commandName]
	if !ok {
		return
	}

	//Load channel to memory
	channel := session.Channel(event.Message.ChannelID)

	//Check guild member permissions
	channelID := event.Message.ChannelID
	memberID := event.Message.Member.UserID
	permissionLevel, err := channelPermissionLevel(session, channelID, memberID, guild.ID)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"member-memberId": memberID,
			"member-name": event.Message.Member.Nick,
			"guild-id": event.Message.GuildID.String(),
			"guild-name": guild.Name,
			"content":  event.Message.Content,
			"command-name": commandName,
		}).WithError(err).Error("Could not get user permissions level")
		return
	}
	if (command.GuildOwnerOnly && memberID.String() != guild.OwnerID.String()) || (command.RequirePermissions && permissionLevel&command.PermissionsLevel <= 0) {
		_, err := session.SendMsg(channelID, fmt.Sprintf(":rotating_light: Brak wystarczających uprawnień do użycia tego polecenia! Wymagany poziom dostępu %d.", command.PermissionsLevel))
		if err != nil {
			Log.WithFields(logrus.Fields{
				"member-memberId": memberID,
				"member-name": event.Message.Member.Nick,
				"guild-id": event.Message.GuildID.String(),
				"guild-name": guild.Name,
				"content":  event.Message.Content,
				"command-name": commandName,
			}).WithError(err).Error("Could not send no permissions warning")
			return
		}
		return
	}

	//Arguments check
	length := len(args[1:])
	if (length > command.MaximumArgumentsSize && command.MaximumArgumentsSize != -1) || (length < command.MinimumArgumentsSize && command.MinimumArgumentsSize != -1) {
		_, err := session.SendMsg(channelID, fmt.Sprintf(":boom: Niepoprawne użycie polecenia %s.", commandName))
		if err != nil {
			Log.WithFields(logrus.Fields{
				"member-memberId": memberID,
				"member-name": event.Message.Member.Nick,
				"guild-id": event.Message.GuildID.String(),
				"guild-name": guild.Name,
				"content":  event.Message.Content,
				"command-name": commandName,
			}).WithError(err).Error("Could send help message usage help")
			return
		}
		_, err = session.SendMsg(channelID, fmt.Sprintf(":star: Wyświetlam pomoc dla polecenia %s: %s", commandName, command.HelpMessage))
		if err != nil {
			Log.WithFields(logrus.Fields{
				"member-memberId": memberID,
				"member-name": event.Message.Member.Nick,
				"guild-id": event.Message.GuildID.String(),
				"guild-name": guild.Name,
				"content":  event.Message.Content,
				"command-name": commandName,
			}).WithError(err).Error("Could send help message usage")
			return
		}
		return
	}

	//Check if command is executed on nsfw channel
	pureChannel, err := channel.Get()
	if err != nil {
		return
	}

	if !pureChannel.NSFW && command.NSFW {
		_, err := session.SendMsg(channelID, ":no_entry: Polecienie __**MOŻE**__ być tylko wykonane na kanale __**NSFW**__!")
		if err != nil {
			Log.WithFields(logrus.Fields{
				"member-memberId": memberID,
				"member-name": event.Message.Member.Nick,
				"guild-id": event.Message.GuildID.String(),
				"guild-name": guild.Name,
				"content":  event.Message.Content,
				"command-name": commandName,
			}).WithError(err).Error("Could not check nsfw channel status")
			return
		}
		return
	}

	//Logs
	Log.WithFields(logrus.Fields{
		"member-memberId": memberID,
		"member-name": event.Message.Member.Nick,
		"guild-id": event.Message.GuildID.String(),
		"guild-name": guild.Name,
		"content":  event.Message.Content,
		"command-name": commandName,
	}).Debug("Executing command")

	//Execute
	err = command.Exec(session, event, guild, args[1:])
	if err != nil {
		Log.WithFields(logrus.Fields{
			"member-memberId": memberID,
			"member-name": event.Message.Member.Nick,
			"guild-id": event.Message.GuildID.String(),
			"guild-name": guild.Name,
			"content":  event.Message.Content,
			"command-name": commandName,
		}).WithError(err).Error("Command general error")
		return
	}
}

func channelPermissionLevel(session disgord.Session, channelID disgord.Snowflake, memberID disgord.Snowflake, guildID disgord.Snowflake) (permissionLevel disgord.PermissionBit, err error) {
	//Get pure channel object
	channel, err := session.Channel(channelID).Get()
	if err != nil {
		return
	}

	guild := session.Guild(guildID)

	member, err := guild.Member(memberID).Get()
	if err != nil {
		return
	}

	permissionLevel, err = channel.GetPermissions(context.TODO(), &GuildCaller{
		session: session,
		guildID: guildID,
	}, member)

	return
}
