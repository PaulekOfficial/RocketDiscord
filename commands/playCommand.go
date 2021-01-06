package commands

import (
	"RocketDiscord/cache"
	"RocketDiscord/utils"
	"github.com/andersfylling/disgord"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

func init() {
	command := NewCommand("play", false, false, onPlayCommand)
	command.SetArgumentsRequirements(true, 1, -1)
	command.SetHelpMessage("Polecenia pozwalajace na dolaczenie bota do kanalu glosowego przez ktorego mozna puscic muzykne na kanale. Użycie: /play <link z youtube lub innego serwisu z muzyką>")
	command.Register()
}

func onPlayCommand(session disgord.Session, event *disgord.MessageCreate, guild *disgord.Guild, args []string) (err error) {
	member := event.Message.Member
	musicState := cache.GetVoiceState(guild.ID)
	channelId := event.Message.ChannelID
	channel, err := session.Channel(channelId).Get()
	if err != nil {
		return err
	}

	if musicState == nil {
		voiceState := cache.GetUserVoiceState(member.UserID)
		if voiceState == nil {
			_, err := session.SendMsg(channelId, ":cloud_tornado: Abym mógł puścić dla ciebie muzykę, musisz znajdować się na kanale głosowym. (Jeżeli jesteś na kanale głosowym poprostu wyjdz i wejdz ponownie na kanał)")
			return err
		}
		voiceChannelId := voiceState.ChannelID

		if voiceChannelId.IsZero() {
			_, err := session.SendMsg(channelId, ":cloud_tornado: Abym mógł puścić dla ciebie muzykę, musisz znajdować się na kanale głosowym. (Jeżeli jesteś na kanale głosowym poprostu wyjdz i wejdz ponownie na kanał)")
			return err
		}

		musicState = &cache.MusicBotState{
			GuildId:        guild.ID,
			VoiceChannelId: voiceChannelId,
			Voice:          nil,
			Tracks:         make(map[int]*cache.MusicBotTrack),
			LoopTrack:      false,
			LoopPlayList:   false,
			Paused:         false,
			Running:        false,
			LastPlay:       time.Now(),
		}

		musicState = cache.PutGuildMusicState(guild.ID, musicState)
	}

	track, err := utils.GetYoutubeStreamMemory(args, session, channelId)
	if err != nil || track == nil {
		_, err := session.SendMsg(channelId, ":sob: Nie znaleziono takiego utworu video/audio.")
		return err
	}

	musicState.Tracks[len(musicState.Tracks)] = track

	if len(musicState.Tracks) == 1 {
		_, err = session.SendMsg(channelId, ":beers: No i jest grane.")
		if err != nil {
			return err
		}

		err = startPlaying(session, musicState, channelId)
		if err != nil {
			return err
		}
	} else {
		user, err := session.User(member.UserID).Get()
		if err != nil {
			return err
		}
		avatarUrl, err := user.AvatarURL(32, false)
		if err != nil {
			return err
		}

		author := &disgord.EmbedAuthor{
			Name:         "Dodano do kolejki",
			URL:          "https://paulek.pro/",
			IconURL:      avatarUrl,
			ProxyIconURL: "",
		}
		fields := []*disgord.EmbedField{
			{
				Name:   "Kanal",
				Value:  channel.Name,
				Inline: true,
			},
			{
				Name:   "Czas trwania utworu",
				Value:  track.Duration.String(),
				Inline: true,
			},
			{
				Name: "Pozycja w kolejce",
				Value: strconv.Itoa(len(musicState.Tracks)),
				Inline: false,
			},
		}
		_, err = session.Channel(channelId).CreateMessage(&disgord.CreateMessageParams{
			Content:                  "",
			Tts:                      false,
			Embed:                    &disgord.Embed{
				Title:       "",
				Author:      author,
				Description: track.URL,
				Timestamp:   disgord.Time{},
				Fields:      fields,
			},
		})
		if err != nil {
			return err
		}
	}
	return
}

func startPlaying(session disgord.Session, musicState *cache.MusicBotState, channelId disgord.Snowflake) error {
	voice, err := session.Guild(musicState.GuildId).VoiceChannel(musicState.VoiceChannelId).Connect(false, true)
	if err != nil {
		return err
	}

	musicState.Voice = voice

	message, err := session.SendMsg(channelId, ":scientist: Statystyki odtwarzanego utworu.")
	if err != nil {
		return err
	}

	go func() {
		err = utils.PlayDCAAudio(musicState.GuildId, session, message)
		if err != nil {
			_, err = session.SendMsg(channelId, ":boom: Oj cos poszlo nie tak, nie mogę puścić muzyki :(.")
			logrus.WithError(err).WithFields(logrus.Fields{
				"guild-id": message.GuildID,
				"music-state": musicState,
			}).Error("Failed to play songs")
		}
		musicState.Running = false
	}()

	return nil
}
