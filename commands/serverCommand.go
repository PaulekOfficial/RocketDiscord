package commands

import (
	"github.com/andersfylling/disgord"
	"strconv"
	"time"
)

func init() {
	command := NewCommand("server", false, false, onServerCommand)
	command.Register()
}

func onServerCommand(session disgord.Session, event *disgord.MessageCreate, guild *disgord.Guild, args []string) error {
	guildQuery := session.Guild(guild.ID)
	members, err := guildQuery.GetMembers(&disgord.GetMembersParams{})
	if err != nil {
		return err
	}
	//channels, err := guildQuery.GetChannels()
	//if err != nil {
	//	return err
	//}
	owner, err := session.User(guild.OwnerID).Get()
	if err != nil {
		return err
	}

	emojiNames := make([]string, len(guild.Emojis))
	for i := 0; i < len(guild.Emojis); i++ {
		emojiNames[i] = ":" + guild.Emojis[i].Name + ":"
	}
	avatarUrl, err := event.Message.Author.AvatarURL(32, false)
	
	footer := &disgord.EmbedFooter{
		Text:         "Wygenerowal " + event.Message.Author.Username,
		IconURL:      avatarUrl,
	}
	
	fields := []*disgord.EmbedField{
		{
			Name:   "Nazwa",
			Value:  guild.Name,
			Inline: true,
		},
		{
			Name:   "ID",
			Value:  guild.ID.String(),
			Inline: true,
		},
		{
			Name:   "Uzytkownikow",
			Value:  strconv.Itoa(len(members)),
			Inline: true,
		},
		{
			Name:   "Region",
			Value:  guild.Region,
			Inline: true,
		},
		{
			Name:   "Wlasciciel",
			Value:  owner.Username,
			Inline: true,
		},
		{
			Name:   "Rangi",
			Value:  strconv.Itoa(len(guild.Roles)),
			Inline: true,
		},
		{
			Name:   "Poziom MFA",
			Value:  strconv.Itoa(int(guild.MFALevel)),
			Inline: true,
		},
		{
			Name:   "Poziom Weryfikacji",
			Value:  strconv.Itoa(int(guild.VerificationLevel)),
			Inline: true,
		},
		{
			Name:   "Emotikonki",
			Value:  strconv.Itoa(len(guild.Emojis)),
			Inline: true,
		},
	}
	_, err = session.Channel(event.Message.ChannelID).CreateMessage(&disgord.CreateMessageParams{
		Content:                  "",
		Tts:                      false,
		Embed:                    &disgord.Embed{
			Title:       "Informacje dotyczące serwera",
			Timestamp:   disgord.Time{Time: time.Now()},
			Color:       30654,
			Fields:      fields,
			Footer:      footer,
		},
	})
	return err
}
