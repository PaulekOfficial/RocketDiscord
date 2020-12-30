package commands

import (
	"github.com/andersfylling/disgord"
)

func init() {
	command := newCommand("server", false, false, onServerCommand)
	command.register()
}

func onServerCommand(session disgord.Session, event *disgord.MessageCreate, guild disgord.GuildQueryBuilder, args []string) error {
	//emojiNames := "nil"
	//members, err := guild.GetMembers(&disgord.GetMembersParams{}, 0)
	//fields := []*discordgo.MessageEmbedField{
	//	{
	//		Name:   "ID Serwera",
	//		Value:  guild.ID.String(),
	//		Inline: true,
	//	},
	//	{
	//		Name:   "ID Własciciela",
	//		Value:  guild.OwnerID.String(),
	//		Inline: true,
	//	},
	//	{
	//		Name:   "Region",
	//		Value:  guild.Region,
	//		Inline: true,
	//	},
	//	//{
	//	//	Name:   "Ilosc dostepnych kanalow",
	//	//	Value:  strconv.Itoa(len(guild.Channels)),
	//	//	Inline: true,
	//	//},
	//	{
	//		Name:   "Ilość uzytkownikow",
	//		Value:  strconv.Itoa(len(members)),
	//		Inline: true,
	//	},
	//	{
	//		Name:   "Serwerowe emoji",
	//		Value:  emojiNames,
	//		Inline: true,
	//	},
	//	{
	//		Name:   "Data powstania",
	//		Value:  "null",
	//		Inline: true,
	//	},
	//}
	//_, err := session.ChannelMessageSendEmbed(event.ChannelID, &discordgo.MessageEmbed{
	//	Title:       "Informacje dotyczące serwera",
	//	Description: "Oto najważniejsze informacje o tym serwerze!",
	//	Color:       30654,
	//	Fields:      fields,
	//	Timestamp:   time.Now().Format(time.RFC3339),
	//})
	return nil
}
