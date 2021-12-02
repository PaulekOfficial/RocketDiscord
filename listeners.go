package main

import (
	"fmt"
	"github.com/andersfylling/disgord"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strings"
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
