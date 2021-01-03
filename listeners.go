package main

import (
	"RocketDiscord/cache"
	"RocketDiscord/commands"
	"fmt"
	"github.com/andersfylling/disgord"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

//////////////////////////////////////////////////////
//
// LISTENERS
//
//////////////////////////////////////////////////////
func ReadyEvent(session disgord.Session, event *disgord.Ready) {
	err := session.UpdateStatusString("Ready to lift off!")
	if err != nil {
		_ = fmt.Errorf("fail to update bot status message %s", err)
	}
}

func MemberAddGuildEvent(session disgord.Session, event *disgord.GuildMemberAdd) {
	username := event.Member.User.Username
	guildSettings := cache.CreateOrGetGuildSettings(event.Member.GuildID)
	channelID := guildSettings.WelcomeChannel

	if channelID.IsZero() {
		return
	}

	memberId := event.Member.User.ID
	guildID := event.Member.GuildID

	avatarUrl, err := event.Member.User.AvatarURL(2048, true)

	if err != nil {
		Log.WithFields(logrus.Fields{
			"member-memberId": memberId,
			"member-name": event.Member.Nick,
			"guild-memberId":  guildID.String(),
		}).WithError(err).Error("Guild member join event fail for avatar url")
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
			"member-name": event.Member.Nick,
			"guild-memberId":  guildID.String(),
		}).WithError(err).Error("Guild member join event fail for getting guild members")
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
			"member-name": event.Member.Nick,
			"guild-memberId":  guildID.String(),
		}).WithError(err).Error("Guild member join event fail for sending embed message")
	}
}

func MemberRemoveGuildEvent(session disgord.Session, event *disgord.GuildMemberRemove) {
	username := event.User.Username
	guildSettings := cache.CreateOrGetGuildSettings(event.GuildID)
	channelID := guildSettings.WelcomeChannel

	if channelID.IsZero() {
		return
	}

	memberId := event.User.ID
	guildID := event.GuildID

	avatarUrl, err := event.User.AvatarURL(2048, true)

	if err != nil {
		Log.WithFields(logrus.Fields{
			"member-memberId": memberId,
			"member-name": event.User.Username,
			"guild-memberId":  guildID.String(),
		}).WithError(err).Error("Guild member quit event fail for avatar url")
	}

	image := &disgord.EmbedImage {
		URL:      avatarUrl,
	}

	guild := session.Guild(guildID)
	members, err := guild.GetMembers(&disgord.GetMembersParams{})

	if err != nil {
		Log.WithFields(logrus.Fields{
			"member-memberId": memberId,
			"member-name": event.User.Username,
			"guild-memberId":  guildID.String(),
		}).WithError(err).Error("Guild member quit event fail for getting guild members")
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
			"member-name": event.User.Username,
			"guild-memberId":  guildID.String(),
		}).WithError(err).Error("Guild member quit event fail for sending embed message")
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
		Log.WithFields(logrus.Fields{
			"member-memberId": event.Message.Member.UserID,
			"member-name": event.Message.Member.Nick,
			"guild-id": event.Message.GuildID.String(),
			"content":  event.Message.Content,
		}).WithError(err).Error("Guild member xd event fail for sending message")
	}
}

func PleasePornGif(session disgord.Session, event *disgord.MessageCreate) {
	if !strings.Contains(strings.ToLower(event.Message.Content), "pls porngif")  {
		return
	}

	_, err := session.SendMsg(event.Message.ChannelID, "( ͡° ͜ʖ ͡°)")
	if err != nil {
		Log.WithFields(logrus.Fields{
			"member-memberId": event.Message.Member.UserID,
			"member-name": event.Message.Member.Nick,
			"guild-id": event.Message.GuildID.String(),
			"content":  event.Message.Content,
		}).WithError(err).Error("Guild member porngif event fail for sending message")
	}
}

func CommandMessageCreate(session disgord.Session, event *disgord.MessageCreate) {
	guild, err := session.Guild(event.Message.GuildID).Get()

	if err != nil {
		Log.WithFields(logrus.Fields{
			"member-memberId": event.Message.Member.UserID,
			"member-name": event.Message.Member.Nick,
			"guild-id": event.Message.GuildID.String(),
			"content":  event.Message.Content,
		}).WithError(err).Error("Guild command fail")
	}

	go func() {
		commands.ParseMessage(session, event, guild, event.Message.Content)
	}()
}
