package main

import (
	"github.com/andersfylling/disgord"
	"github.com/sirupsen/logrus"
	"time"
)

type LoggedUser struct {
	ID            int
	UserID        string
	Username      string
	Discriminator uint16
	Email         string
	Avatar        string
	Token         string
	Verified      bool
	MFAEnabled    bool
	Bot           bool
	PremiumType   int
	Locale        string
	Flags         uint64
	PublicFlags   uint64
}

type LoggedMessage struct {
	ID             int
	UserID         string
	GuildID        string
	ChannelID      string
	MessageID      string
	Content        string
	UpdatedContent *string
	CreateTime     time.Time
	UpdateTime     *time.Time
}

func RegisterModLogModule() {
	dbMap.AddTableWithName(LoggedMessage{}, "rocketDiscord_messages").SetKeys(true, "ID")
	dbMap.AddTableWithName(LoggedUser{}, "rocketDiscord_users").SetKeys(true, "ID")
	err := dbMap.CreateTablesIfNotExists()
	if err != nil {
		Log.WithError(err).Error("Cannot create modLog tables")
	}
}

func LogMessage(session disgord.Session, event *disgord.MessageCreate) {
	err := dbMap.Insert(&LoggedMessage{
		GuildID:        event.Message.GuildID.String(),
		UserID:         event.Message.Author.ID.String(),
		ChannelID:      event.Message.ChannelID.String(),
		MessageID:      event.Message.ID.String(),
		Content:        event.Message.Content,
		CreateTime:     event.Message.Timestamp.Time,
	})

	if err != nil {
		Log.WithFields(logrus.Fields{
			//"member-memberId": event.Message.Member.UserID,
			//"member-name": event.Message.Member.Nick,
			"guild-id": event.Message.GuildID.String(),
			"content":  event.Message.Content,
		}).Errorf("Guild command fail", err)
	}
}

func LogMessageUpdate(session disgord.Session, event *disgord.MessageUpdate) {
	var message *LoggedMessage
	err := dbMap.SelectOne(&message, "SELECT * FROM rocketDiscord_messages WHERE ChannelID=? AND MessageID=?", event.Message.ChannelID.String(), event.Message.ID.String())

	if message == nil {
		return
	}

	if err != nil {
		Log.WithFields(logrus.Fields{
			//"member-memberId": event.Message.Member.UserID,
			//"member-name": event.Message.Member.Nick,
			"guild-id": event.Message.GuildID.String(),
			"content":  event.Message.Content,
		}).Errorf("Guild command fail", err)
	}

	message.UpdatedContent = &event.Message.Content
	message.UpdateTime = &event.Message.Timestamp.Time

	_, err = dbMap.Update(message)

	if err != nil {
		Log.WithFields(logrus.Fields{
			//"member-memberId": event.Message.Member.UserID,
			//"member-name": event.Message.Member.Nick,
			"guild-id": event.Message.GuildID.String(),
			"content":  event.Message.Content,
		}).Errorf("Guild command fail", err)
	}
}

func LogUserUpdate(session disgord.Session, event *disgord.UserUpdate) {
	err := dbMap.Insert(&LoggedUser{
		UserID:        event.ID.String(),
		Username:      event.Username,
		Discriminator: uint16(event.Discriminator),
		Email:         event.Email,
		Avatar:        event.Avatar,
		Token:         event.Token,
		Verified:      event.Verified,
		MFAEnabled:    event.MFAEnabled,
		Bot:           event.Bot,
		PremiumType:   int(event.PremiumType),
		Locale:        event.Locale,
		Flags:         uint64(event.Flags),
		PublicFlags:   uint64(event.PublicFlags),
	})

	if err != nil {
		Log.WithFields(logrus.Fields{
			"member-memberId": event.ID,
			"member-name": event.Username,
		}).Errorf("Guild command fail", err)
	}
}
