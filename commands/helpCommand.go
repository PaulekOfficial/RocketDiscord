package commands

import (
	"github.com/andersfylling/disgord"
)

func init() {
	command := NewCommand("help", false, false, onHelpCommand)
	command.Register()
}

func onHelpCommand(session disgord.Session, event *disgord.MessageCreate, guild *disgord.Guild, args []string) error {
	_, err := session.SendMsg(event.Message.ChannelID, ":satellite_orbital: There is no help for u.")
	if err != nil {
		return err
	}
	return nil
}
