package commands

import "github.com/bwmarrin/discordgo"

func init() {
	command := newCommand("help", false, false, onHelpCommand)
	command.register()
}

func onHelpCommand(session *discordgo.Session, event *discordgo.MessageCreate, guild *discordgo.Guild, args []string) error {
	_, err := session.ChannelMessageSend(event.ChannelID, ":satellite_orbital: There is no help for u.")
	if err != nil {
		return err
	}
	return nil
}
