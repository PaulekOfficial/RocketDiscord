package commands

import "github.com/bwmarrin/discordgo"

func init() {
	command := newCommand("ban", true, false, onBanCommand)
	command.setPermissions(true, false, discordgo.PermissionAdministrator)
	command.setArgumentsRequirements(true, 1, -1)
	command.setHelpMessage("Polecenie pozwalające zbanowanie uzytkownika, wymaga conajmniej uprawnień administratora. Użycie: !ban @user <reason>")
	command.register()
}

func onBanCommand(session *discordgo.Session, event *discordgo.MessageCreate, guild *discordgo.Guild, args []string) error {
	_, err := session.ChannelMessageSend(event.ChannelID, ":satellite_orbital: Not implemented yet")
	if err != nil {
		return err
	}
	return nil
}
