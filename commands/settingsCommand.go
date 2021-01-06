package commands

import (
	"github.com/andersfylling/disgord"
)

func init() {
	command := NewCommand("settings", false, false, onSettingsCommand)

	command.SetArgumentsRequirements(true, -1, -1)
	command.SetHelpMessage("Polecenie rasy bogin, tylko nieliczni maja dostep do niego. Użycie: /settings <setting-name> <value>")

	command.PermissionsLevel = disgord.PermissionAdministrator
	command.RequirePermissions = true

	command.Register()
}

func onSettingsCommand(session disgord.Session, event *disgord.MessageCreate, guild *disgord.Guild, args []string) (err error) {
	general := []byte("")
	general = append(general, "```"...)
	general = append(general, ".                ."...)
	general = append(general, ".                ."...)
	general = append(general, ".                ."...)
	general = append(general, ".                ."...)
	general = append(general, ".                ."...)
	general = append(general, ".                ."...)
	general = append(general, ".                ."...)
	general = append(general, "```"...)


	fields := []*disgord.EmbedField{
		{
			Name:   "Ogólne",
			Value:  string(general),
			Inline: true,
		},
		{
			Name:   "Wiadomości",
			Value:  string(general),
			Inline: true,
		},
		{
			Name: "Kanały",
			Value: string(general),
			Inline: true,
		},
	}

	_, err = session.Channel(event.Message.ChannelID).CreateMessage(&disgord.CreateMessageParams{
		Content:                  "",
		Tts:                      false,
		Embed:                    &disgord.Embed{
			Title:       "",
			Timestamp:   disgord.Time{},
			Fields:      fields,
		},
	})
	if err != nil {
		return err
	}

	return
}
