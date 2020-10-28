package commands

import (
	"github.com/bwmarrin/discordgo"
	"strconv"
	"strings"
	"time"
)

func init() {
	command := newCommand("server", false, false, onServerCommand)
	command.register()
}

func onServerCommand(session *discordgo.Session, event *discordgo.MessageCreate, guild *discordgo.Guild, args []string) error {
	//timestamp, err := guild.JoinedAt.Parse()
	//if err != nil {
	//	panic(err)
	//	return err
	//}
	emojiNames := make([]string, len(guild.Emojis))
	for i := 0; i < len(guild.Emojis); i++ {
		emojiNames[i] = ":" + guild.Emojis[i].Name + ":"
	}
	fields := []*discordgo.MessageEmbedField{
		{
			Name:   "ID Serwera",
			Value:  guild.ID,
			Inline: true,
		},
		{
			Name:   "ID Własciciela",
			Value:  guild.OwnerID,
			Inline: true,
		},
		{
			Name:   "Region",
			Value:  guild.Region,
			Inline: true,
		},
		{
			Name:   "Ilosc dostepnych kanalow",
			Value:  strconv.Itoa(len(guild.Channels)),
			Inline: true,
		},
		{
			Name:   "Ilość uzytkownikow",
			Value:  strconv.Itoa(guild.MemberCount),
			Inline: true,
		},
		{
			Name:   "Serwerowe emoji",
			Value:  strings.Join(emojiNames, " "),
			Inline: true,
		},
		{
			Name:   "Data powstania",
			Value:  "null",
			Inline: true,
		},
	}
	_, err := session.ChannelMessageSendEmbed(event.ChannelID, &discordgo.MessageEmbed{
		Title:       "Informacje dotyczące serwera",
		Description: "Oto najważniejsze informacje o tym serwerze!",
		Color:       30654,
		Fields:      fields,
		Timestamp:   time.Now().Format(time.RFC3339),
	})
	return err
}
