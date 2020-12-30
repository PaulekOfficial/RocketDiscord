package main

import (
	"RocketDiscord/commands"
	"fmt"
	"github.com/andersfylling/disgord"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func ReadyEvent(session disgord.Session, event *disgord.Ready) {
	err := session.UpdateStatusString("Ready to lift off!")
	if err != nil {
		_ = fmt.Errorf("fail to update bot status message %s", err)
	}
}

func MemberAddGuildEvent(session disgord.Session, event *disgord.GuildMemberAdd) {
	username := event.Member.User.Username
	channelID := disgord.Snowflake(770728332773425203)

	memberId := event.Member.User.ID
	guildID := event.Member.GuildID

	avatarUrl, err := event.Member.User.AvatarURL(2048, true)

	if err != nil {
		Log.WithFields(logrus.Fields{
			"member-memberId": memberId,
			"guild-memberId":  guildID.String(),
		}).Errorf("Guild member join event fail for avatar url", err)
	}

	image := &disgord.EmbedImage {
		URL:      avatarUrl,
	}

	timestamp := event.Member.JoinedAt
	guild := session.Guild(guildID)
	members, err := guild.GetMembers(&disgord.GetMembersParams{}, 0)

	if err != nil {
		Log.WithFields(logrus.Fields{
			"member-memberId": memberId,
			"guild-memberId":  guildID.String(),
		}).Errorf("Guild member join event fail for getting guild members", err)
	}

	fields := []*disgord.EmbedField{
		{
			Name:   "Czas dołączenia",
			Value:  timestamp.Format(time.ANSIC),
			Inline: true,
		},
		{
			Name:   "ID",
			Value:  memberId.String(),
			Inline: true,
		},
		{
			Name: "Nowa ilość użytkowników",
			Value: strconv.Itoa(len(members)),
			Inline: true,
		},
	}
	_, err = session.Channel(channelID).CreateMessage(&disgord.CreateMessageParams{
		Content:                  "",
		Tts:                      false,
		Embed:                    &disgord.Embed{
			Title:       username + " dołączył na nasz serwer!",
			Description: "Nowy użytkownik dołączył do naszej społeczności, witamy!",
			Timestamp:   disgord.Time{},
			Color:       501767,
			Image:       image,
			Fields:      fields,
		},
	})

	if err != nil {
		Log.WithFields(logrus.Fields{
			"member-memberId": memberId,
			"guild-memberId":  guildID.String(),
		}).Errorf("Guild member join event fail for sending embed message", err)
	}
}

func MemberRemoveGuildEvent(session disgord.Session, event *disgord.GuildMemberRemove) {
	username := event.User.Username
	channelID := disgord.Snowflake(770728332773425203)

	memberId := event.User.ID
	guildID := event.GuildID

	avatarUrl, err := event.User.AvatarURL(2048, true)

	if err != nil {
		Log.WithFields(logrus.Fields{
			"member-memberId": memberId,
			"guild-memberId":  guildID.String(),
		}).Errorf("Guild member quit event fail for avatar url", err)
	}

	image := &disgord.EmbedImage {
		URL:      avatarUrl,
	}

	guild := session.Guild(guildID)
	members, err := guild.GetMembers(&disgord.GetMembersParams{}, 0)

	if err != nil {
		Log.WithFields(logrus.Fields{
			"member-memberId": memberId,
			"guild-memberId":  guildID.String(),
		}).Errorf("Guild member quit event fail for getting guild members", err)
	}

	fields := []*disgord.EmbedField{
		{
			Name:   "Czas odłączenia",
			Value:  time.Now().Format(time.ANSIC),
			Inline: true,
		},
		{
			Name:   "ID",
			Value:  memberId.String(),
			Inline: true,
		},
		{
			Name: "Nowa ilość użytkowników",
			Value: strconv.Itoa(len(members)),
			Inline: true,
		},
	}
	_, err = session.Channel(channelID).CreateMessage(&disgord.CreateMessageParams{
		Content:                  "",
		Tts:                      false,
		Embed:                    &disgord.Embed{
			Title:       username + " opuścił nasz serwerek! :(",
			Description: "Użytkownik opuścił serwer, będzie nam go brakowało!",
			Timestamp:   disgord.Time{},
			Color:       14558244,
			Image:       image,
			Fields:      fields,
		},
	})

	if err != nil {
		Log.WithFields(logrus.Fields{
			"member-memberId": memberId,
			"guild-memberId":  guildID.String(),
		}).Errorf("Guild member quit event fail for sending embed message", err)
	}
}

func MessageXD(session disgord.Session, event *disgord.MessageCreate) {
	if !strings.Contains(strings.ToLower(event.Message.Content), "xd")  {
		return
	}

	if rand.Intn(100) < 50 {
		return
	}

	_, err := session.SendMsg(event.Message.ChannelID, "iks de")
	if err != nil {
		_ = fmt.Errorf("fail to send iks de message %s", err)
	}
}

func PleasePornGif(session disgord.Session, event *disgord.MessageCreate) {
	if !strings.Contains(strings.ToLower(event.Message.Content), "pls porngif")  {
		return
	}

	_, err := session.SendMsg(event.Message.ChannelID, "( ͡° ͜ʖ ͡°)")
	if err != nil {
		_ = fmt.Errorf("fail to send porngif message %s", err)
	}
}

func CommandMessageCreate(session disgord.Session, event *disgord.MessageCreate) {
	guild := session.Guild(event.Message.GuildID)
	commands.ParseMessage(session, event, guild, event.Message.Content)
}
