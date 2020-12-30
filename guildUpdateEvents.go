package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

type RocketGuild struct {
	ID                       int
	GuildID                  string
	Name                     string
	Icon                     string
	Region                   string
	AfkChannelID             string
	EmbedChannelID           string
	OwnerID                  string
	Owner                    bool
	JoinedAt                 time.Time
	AfkTimeout               int
	MemberCount              int
	EmbedEnabled             bool
	Large                    bool
	MaxMembers               int
	Unavailable              bool
	WidgetEnabled            bool
	WidgetChannelID          string
	SystemChannelID          string
	RulesChannelID           string
	VanityURLCode            string
	Description              string
	Banner                   string
	PremiumSubscriptionCount int
	PreferredLocale          string
	PublicUpdatesChannelID   string
	MaxVideoChannelUsers     int
	ApproximateMemberCount   int
	ApproximatePresenceCount int
	Permissions              int
	Timestamp                time.Time
}

func RegisterUpdateListeners(session *discordgo.Session) {
	dbMap.AddTableWithName(RocketGuild{}, "rocketdiscord_Guilds").SetKeys(true, "ID")
	err := dbMap.CreateTablesIfNotExists()
	if err != nil {
		panic(err)
	}

	session.AddHandler(guildJoin)
	session.AddHandler(guildUpdate)

	session.AddHandler(guildMemberUpdate)
	session.AddHandler(guildMemberAdd)
	session.AddHandler(guildMemberRemove)
}

func guildJoin(session *discordgo.Session, event *discordgo.GuildCreate) {
	//joinedAt, err := event.JoinedAt.Parse()
	//if err != nil {
	//	panic(err)
	//}

	err := dbMap.Insert(&RocketGuild{
		GuildID:                  event.ID,
		Name:                     event.Name,
		Icon:                     event.Icon,
		Region:                   event.Region,
		AfkChannelID:             event.AfkChannelID,
		EmbedChannelID:           event.EmbedChannelID,
		OwnerID:                  event.OwnerID,
		Owner:                    event.Owner,
		//JoinedAt:                 event.JoinedAt,
		JoinedAt:                 time.Now(),
		AfkTimeout:               event.AfkTimeout,
		EmbedEnabled:             event.EmbedEnabled,
		Large:                    event.Large,
		MaxMembers:               event.MaxMembers,
		Unavailable:              event.Unavailable,
		WidgetEnabled:            event.WidgetEnabled,
		WidgetChannelID:          event.WidgetChannelID,
		SystemChannelID:          event.SystemChannelID,
		RulesChannelID:           event.RulesChannelID,
		VanityURLCode:            event.VanityURLCode,
		Description:              event.Description,
		Banner:                   event.Banner,
		PremiumSubscriptionCount: event.PremiumSubscriptionCount,
		PreferredLocale:          event.PreferredLocale,
		PublicUpdatesChannelID:   event.PublicUpdatesChannelID,
		MaxVideoChannelUsers:     event.MaxVideoChannelUsers,
		Permissions:              event.Permissions,
		Timestamp:                time.Now(),
	})

	if err != nil {
		panic(err)
	}
}

func guildUpdate(session *discordgo.Session, event *discordgo.GuildUpdate) {
	//joinedAt, err := event.JoinedAt.Parse()
	//if err != nil {
	//	panic(err)
	//}

	count, err := dbMap.Update(&RocketGuild{
		GuildID:                  event.ID,
		Name:                     event.Name,
		Icon:                     event.Icon,
		Region:                   event.Region,
		AfkChannelID:             event.AfkChannelID,
		EmbedChannelID:           event.EmbedChannelID,
		OwnerID:                  event.OwnerID,
		Owner:                    event.Owner,
		//JoinedAt:                 event.JoinedAt,
		JoinedAt:                 time.Now(),
		AfkTimeout:               event.AfkTimeout,
		MemberCount:              event.MemberCount,
		EmbedEnabled:             event.EmbedEnabled,
		Large:                    event.Large,
		Unavailable:              event.Unavailable,
		WidgetEnabled:            event.WidgetEnabled,
		WidgetChannelID:          event.WidgetChannelID,
		SystemChannelID:          event.SystemChannelID,
		RulesChannelID:           event.RulesChannelID,
		VanityURLCode:            event.VanityURLCode,
		Description:              event.Description,
		Banner:                   event.Banner,
		PremiumSubscriptionCount: event.PremiumSubscriptionCount,
		PreferredLocale:          event.PreferredLocale,
		PublicUpdatesChannelID:   event.PublicUpdatesChannelID,
		MaxVideoChannelUsers:     event.MaxVideoChannelUsers,
		Permissions:              event.Permissions,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("Guild %s update, affected %d rows", event.Name, count))
}

func guildMemberUpdate(session *discordgo.Session, event *discordgo.GuildMemberUpdate) {
	guild, err := session.Guild(event.GuildID)
	if err != nil {
		_ = fmt.Errorf("error while getting guild pointer on member update %v", err)
		return
	}
	count, err := dbMap.Update(&RocketGuild{
		GuildID:                  guild.ID,
		MemberCount:              guild.MemberCount,
		ApproximateMemberCount:   guild.ApproximateMemberCount,
		ApproximatePresenceCount: guild.ApproximatePresenceCount,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("Guild members join %s update, affected %d rows", guild.Name, count))
}

func guildMemberAdd(session *discordgo.Session, event *discordgo.GuildMemberAdd) {
	guild, err := session.Guild(event.GuildID)
	if err != nil {
		_ = fmt.Errorf("error while getting guild pointer on member update %v", err)
		return
	}
	count, err := dbMap.Update(&RocketGuild{
		GuildID:                  guild.ID,
		MemberCount:              guild.MemberCount,
		ApproximateMemberCount:   guild.ApproximateMemberCount,
		ApproximatePresenceCount: guild.ApproximatePresenceCount,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("Guild members quit %s update, affected %d rows", guild.Name, count))
}

func guildMemberRemove(session *discordgo.Session, event *discordgo.GuildMemberRemove) {
	guild, err := session.Guild(event.GuildID)
	if err != nil {
		_ = fmt.Errorf("error while getting guild pointer on member update %v", err)
		return
	}
	count, err := dbMap.Update(&RocketGuild{
		GuildID:                  guild.ID,
		MemberCount:              guild.MemberCount,
		ApproximateMemberCount:   guild.ApproximateMemberCount,
		ApproximatePresenceCount: guild.ApproximatePresenceCount,
	}, "")
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("Guild members %s update, affected %d rows", guild.Name, count))
}
