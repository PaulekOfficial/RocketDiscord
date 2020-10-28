package commands

import "github.com/bwmarrin/discordgo"

func init() {
	command := newCommand("rule34", false, true, onRule34Request)
	command.register()
}

func onRule34Request(session *discordgo.Session, event *discordgo.MessageCreate, guild *discordgo.Guild, args []string) error {
	_, err := session.ChannelMessageSend(event.ChannelID, "Not implemented yet.")
	if err != nil {
		return err
	}
	return nil
}
