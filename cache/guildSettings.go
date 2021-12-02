package cache

import "github.com/andersfylling/disgord"

var guilds = make(map[disgord.Snowflake]*GuildSettings)

type GuildSettings struct {
	Id                 disgord.Snowflake
	CommandPrefix      string

	EarnExperience        bool
	EarnFromMessage       int
	EarnFromVoiceTime     int
	EarnFromMeme          int

	TrustLevels           bool
	MinimumTrustLevel     int
	MaximumTrustLevel     int
	TrustRoles            map[int]disgord.Role
	TrustExp              map[int]int
	TrustNextLevelMessage string
	TrustNextLevelImage   string

	AutoChannels          bool
	AutoChannelsList      []disgord.Snowflake
	AutoChannelName       string
	AutoChannelDelete     int

	ModLogChannel         disgord.Snowflake
	WelcomeChannel        disgord.Snowflake
	CleverBotChannels     []disgord.Snowflake
	MemesChannels         []disgord.Snowflake

	WelcomeMessage        string
	WelcomeImage          string
	WelcomeEmbedStats     bool

	QuitMessage           string
	QuitImage             string
	QuitEmbedStats        bool

	ReactRoles            bool
	ReactMessages         map[disgord.Snowflake]ReactMessage

	SaveRoles             bool
	ChatLogs              bool
	AntiMention           bool
	MutedRole             disgord.Role
	TrackDeletes          bool
	TrackModifications    bool
	AdminRole             disgord.Role

	MuteWarningsAt        int
	KickWarningsAt        int
	BanWarningsAt         int

	VoiceBot              bool
	VoiceDjRole           disgord.Role
	VoiceMaxQueue         int
	VoiceLoopTrack        bool
	VoiceLoopPlaylist     bool
	VoiceTimeout          int
}

type ReactMessage struct {
	ChannelID     disgord.Snowflake
	ReactionRoles map[string]disgord.Role
}

func CreateOrGetGuildSettings(guild *disgord.Guild) *GuildSettings {
	settings := GetGuildSettings(guild.ID)
	if settings == nil {
		settings = createGuildSettings(guild)
		guilds[guild.ID] = settings
	}

	return settings
}

func createGuildSettings(guild *disgord.Guild) *GuildSettings {
	guildSettings := &GuildSettings{
		Id:                    guild.ID,
		CommandPrefix:         "rd!",
		EarnExperience:        false,
		EarnFromMessage:       0,
		EarnFromVoiceTime:     0,
		EarnFromMeme:          0,
		TrustLevels:           false,
		MinimumTrustLevel:     0,
		MaximumTrustLevel:     0,
		TrustRoles:            make(map[int]disgord.Role),
		TrustExp:              make(map[int]int),
		TrustNextLevelMessage: "",
		TrustNextLevelImage:   "",
		AutoChannels:          false,
		AutoChannelsList:      make([]disgord.Snowflake, 0),
		AutoChannelName:       "",
		AutoChannelDelete:     0,
		ModLogChannel:         0,
		WelcomeChannel:        0,
		CleverBotChannels:     make([]disgord.Snowflake, 0),
		MemesChannels:         make([]disgord.Snowflake, 0),
		WelcomeMessage:        "",
		WelcomeImage:          "",
		WelcomeEmbedStats:     false,
		QuitMessage:           "",
		QuitImage:             "",
		QuitEmbedStats:        false,
		ReactRoles:            false,
		ReactMessages:         make(map[disgord.Snowflake]ReactMessage),
		SaveRoles:             false,
		ChatLogs:              false,
		AntiMention:           false,
		MutedRole:             disgord.Role{},
		TrackDeletes:          false,
		TrackModifications:    false,
		AdminRole:             disgord.Role{},
		MuteWarningsAt:        0,
		KickWarningsAt:        0,
		BanWarningsAt:         0,
		VoiceBot:              false,
		VoiceDjRole:           disgord.Role{},
		VoiceMaxQueue:         0,
		VoiceLoopTrack:        false,
		VoiceLoopPlaylist:     false,
		VoiceTimeout:          0,
	}
	return guildSettings
}

func createGuildSettingsRecord(settings *GuildSettings) {

}

func GetGuildSettings(guildId disgord.Snowflake) *GuildSettings {
	return guilds[guildId]
}
