package commands

import "github.com/bwmarrin/discordgo"

func init() {
	command := newCommand("rule34", false, true, onRule34Request)
	command.setArgumentsRequirements(true, 1, 1)
	command.setHelpMessage("Polecenie pozwalające pokazać obrazek ze strony rule34. Ostrzeżenie treść jest przeznaczona dla osób pełnoletnich!  Użycie: !rule34 <tag>")
	command.register()
}

func onRule34Request(session *discordgo.Session, event *discordgo.MessageCreate, guild *discordgo.Guild, args []string) error {
	_, err := session.ChannelMessageSend(event.ChannelID, ":satellite_orbital: Not implemented yet.")
	if err != nil {
		return err
	}
	return nil
}
