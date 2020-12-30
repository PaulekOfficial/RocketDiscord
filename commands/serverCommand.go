package commands

import (
	"github.com/andersfylling/disgord"
	"strconv"
	"strings"
)

func init() {
	command := newCommand("server", false, false, onServerCommand)
	command.register()
}

func onServerCommand(session disgord.Session, event *disgord.MessageCreate, guild *disgord.Guild, args []string) error {
	guildQuery := session.Guild(guild.ID)
	members, err := guildQuery.GetMembers(&disgord.GetMembersParams{})
	if err != nil {
		return err
	}
	channels, err := guildQuery.GetChannels()
	if err != nil {
		return err
	}

	emojiNames := make([]string, len(guild.Emojis))
	for i := 0; i < len(guild.Emojis); i++ {
		emojiNames[i] = ":" + guild.Emojis[i].Name + ":"
	}
	fields := []*disgord.EmbedField{
		{
			Name:   "ID Serwera",
			Value:  guild.ID.String(),
			Inline: true,
		},
		{
			Name:   "ID Własciciela",
			Value:  guild.OwnerID.String(),
			Inline: true,
		},
		{
			Name:   "Region",
			Value:  guild.Region,
			Inline: true,
		},
		{
			Name:   "Ilosc dostepnych kanalow",
			Value:  strconv.Itoa(len(channels)),
			Inline: true,
		},
		{
			Name:   "Ilość uzytkownikow",
			Value:  strconv.Itoa(len(members)),
			Inline: true,
		},
		{
			Name:   "Serwerowe emoji",
			Value:  strings.Join(emojiNames, " "),
			Inline: true,
		},
	}
	_, err = session.Channel(event.Message.ChannelID).CreateMessage(&disgord.CreateMessageParams{
		Content:                  "",
		Tts:                      false,
		Embed:                    &disgord.Embed{
			Title:       "Informacje dotyczące serwera",
			Description: "Oto najważniejsze informacje o tym serwerze!",
			Timestamp:   disgord.Time{},
			Color:       30654,
			//Image:       image,
			Fields:      fields,
		},
	})
	return err
}
