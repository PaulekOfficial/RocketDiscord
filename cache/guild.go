package cache

import "github.com/andersfylling/disgord"

var guilds = make(map[disgord.Snowflake]*GuildSettings)

type GuildSettings struct {
	Id               disgord.Snowflake
	ModLogChannel    disgord.Snowflake
	WelcomeChannel   disgord.Snowflake
	CommandPrefix    string
	AutoChannels     bool
	TrustLevels      bool
	DmHelp           bool
	DisableEveryone  bool
	DisableHere      bool
}

func CreateOrGetGuildSettings(guildId disgord.Snowflake) *GuildSettings {
	settings := GetGuildSettings(guildId)
	if settings != nil {
		return settings
	}

	settings = &GuildSettings{
		Id:              guildId,
		ModLogChannel:   0,
		WelcomeChannel:  0,
		AutoChannels:    false,
		TrustLevels:     false,
		DmHelp:          false,
		CommandPrefix:   "?",
		DisableEveryone: false,
		DisableHere:     false,
	}
	PutGuildSettings(settings)

	return settings
}

func PutGuildSettings(settings *GuildSettings) {
	if settings == nil {
		return
	}

	guilds[settings.Id] = settings
}

func GetGuildSettings(guildId disgord.Snowflake) *GuildSettings {
	return guilds[guildId]
}
