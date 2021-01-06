package commands

import (
	"github.com/andersfylling/disgord"
)

func init() {
	command := NewCommand("ban", true, false, onBanCommand)
	command.SetPermissions(true, false, disgord.PermissionAdministrator)
	command.SetArgumentsRequirements(true, 1, -1)
	command.SetHelpMessage("Polecenie pozwalające zbanowanie uzytkownika, wymaga conajmniej uprawnień administratora. Użycie: !ban @user <reason>")
	command.Register()
}

func onBanCommand(session disgord.Session, event *disgord.MessageCreate, guild *disgord.Guild, args []string) error {
	_, err := session.SendMsg(event.Message.ChannelID, ":satellite_orbital: Not implemented yet")
	if err != nil {
		return err
	}
	return nil
}
