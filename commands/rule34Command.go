package commands

import (
	"github.com/andersfylling/disgord"
)

func init() {
	command := newCommand("rule34", false, true, onRule34Request)
	command.setArgumentsRequirements(true, 1, 1)
	command.setHelpMessage("Polecenie pozwalające pokazać obrazek ze strony rule34. Ostrzeżenie treść jest przeznaczona dla osób pełnoletnich!  Użycie: !rule34 <tag>")
	command.register()
}

func onRule34Request(session disgord.Session, event *disgord.MessageCreate, guild disgord.GuildQueryBuilder, args []string) error {
	_, err := session.SendMsg(event.Message.ChannelID, ":satellite_orbital: Not implemented yet.")
	if err != nil {
		return err
	}
	return nil
}
