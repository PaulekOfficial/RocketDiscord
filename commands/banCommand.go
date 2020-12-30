package commands

import (
	"github.com/andersfylling/disgord"
	"github.com/bwmarrin/discordgo"
)

func init() {
	command := newCommand("ban", true, false, onBanCommand)
	command.setPermissions(true, false, discordgo.PermissionAdministrator)
	command.setArgumentsRequirements(true, 1, -1)
	command.setHelpMessage("Polecenie pozwalające zbanowanie uzytkownika, wymaga conajmniej uprawnień administratora. Użycie: !ban @user <reason>")
	command.register()
}

func onBanCommand(session disgord.Session, event *disgord.MessageCreate, guild disgord.GuildQueryBuilder, args []string) error {
	_, err := session.SendMsg(event.Message.ChannelID, ":satellite_orbital: Not implemented yet")
	if err != nil {
		return err
	}
	return nil
}
