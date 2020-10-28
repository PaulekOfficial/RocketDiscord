package main

import (
	"RocketDiscord/commands"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func ReadyEvent(session *discordgo.Session, event *discordgo.Ready) {
	err := session.UpdateStatus(0, "Ready to lift off!")
	if err != nil {
		_ = fmt.Errorf("fail to update bot status message %s", err)
	}
}

func MemberAddGuildEvent(session *discordgo.Session, event *discordgo.GuildMemberAdd) {
	username := event.Member.User.Username
	_, err := session.ChannelMessageSend("770726177194770513", "Powitajmy nowego użytkownika, " + username)
	if err != nil {
		_ = fmt.Errorf("fail to send welcome message %s", err)
	}

	id := event.User.ID
	image := &discordgo.MessageEmbedImage{
		URL:      event.User.AvatarURL(""),
	}
	timestamp, err := event.Member.JoinedAt.Parse()
	if err != nil {
		_ = fmt.Errorf("fail to parse member join timestamp %s", err)
	}
	guild, err := session.Guild(event.GuildID)
	if err != nil {
		_ = fmt.Errorf("fail to get guild via id %s", err)
	}
	fields := []*discordgo.MessageEmbedField{
		{
			Name:   "Czas dołączenia",
			Value:  timestamp.Format(time.ANSIC),
			Inline: true,
		},
		{
			Name:   "ID",
			Value:  id,
			Inline: true,
		},
		{
			Name: "Nowa ilość użytkowników",
			Value: strconv.Itoa(guild.MemberCount),
			Inline: true,
		},
	}
	_, err = session.ChannelMessageSendEmbed("770728332773425203", &discordgo.MessageEmbed{
		Title:       username + " dołączył na nasz serwer!",
		Description: "Nowy użytkownik dołączył do naszej społeczności, witamy!",
		Image:       image,
		Color:       501767,
		Fields:      fields,
		Timestamp: time.Now().Format(time.RFC3339),
	})
	if err != nil {
		_ = fmt.Errorf("fail to send welcome embed on welcome channel %s", err)
	}
}

func MemberRemoveGuildEvent(session *discordgo.Session, event *discordgo.GuildMemberRemove) {
	username := event.Member.User.Username
	id := event.User.ID
	image := &discordgo.MessageEmbedImage{
		URL:      event.User.AvatarURL(""),
	}
	guild, err := session.Guild(event.GuildID)
	if err != nil {
		_ = fmt.Errorf("fail to get guild via id %s", err)
	}
	fields := []*discordgo.MessageEmbedField{
		{
			Name:   "Czas odłączenia",
			Value:  time.Now().Format(time.ANSIC),
			Inline: true,
		},
		{
			Name:   "ID",
			Value:  id,
			Inline: true,
		},
		{
			Name: "Nowa ilość użytkowników",
			Value: strconv.Itoa(guild.MemberCount),
			Inline: true,
		},
	}
	_, err = session.ChannelMessageSendEmbed("770728332773425203", &discordgo.MessageEmbed{
		Title:       username + " opuścił nasz serwerek! :(",
		Description: "Użytkownik opuścił serwer, będzie nam go brakowało!",
		Image:       image,
		Color:       14558244,
		Fields:      fields,
		Timestamp: time.Now().Format(time.RFC3339),
	})
	if err != nil {
		_ = fmt.Errorf("fail to send welcome embed on welcome channel %s", err)
	}
}

func MessageXD(session *discordgo.Session, event *discordgo.MessageCreate) {
	if !strings.Contains(strings.ToLower(event.Content), "xd")  {
		return
	}

	if rand.Intn(100) < 50 {
		return
	}

	_, err := session.ChannelMessageSend(event.ChannelID, "iks de")
	if err != nil {
		_ = fmt.Errorf("fail to send iks de message %s", err)
	}
}

func PleasePornGif(session *discordgo.Session, event *discordgo.MessageCreate) {
	if !strings.Contains(strings.ToLower(event.Content), "pls porngif")  {
		return
	}

	_, err := session.ChannelMessageSend(event.ChannelID, "( ͡° ͜ʖ ͡°)")
	if err != nil {
		_ = fmt.Errorf("fail to send porngif message %s", err)
	}
}

func CommandMessageCreate(session *discordgo.Session, event *discordgo.MessageCreate) {
	guild, err := session.Guild(event.GuildID)
	if err != nil {
		_ = fmt.Errorf("fail to get guild via id %s", err)
	}
	commands.ParseMessage(session, event, guild, event.Content)
}
