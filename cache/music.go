package cache

import (
	"github.com/andersfylling/disgord"
	"github.com/sirupsen/logrus"
	"io"
	"time"
)

var voiceStates = make(map[disgord.Snowflake]*disgord.VoiceState)

var musicStates = make(map[disgord.Snowflake]*MusicBotState)

type MusicBotState struct {
	GuildId        disgord.Snowflake
	VoiceChannelId disgord.Snowflake
	Voice          disgord.VoiceConnection
	Tracks         map[int]*MusicBotTrack
	LoopTrack      bool
	LoopPlayList   bool
	Paused         bool
	Running        bool
	LastPlay       time.Time
}

type RawSavedFrame struct {
	Payload  []byte
	Duration time.Duration
}

type MusicBotTrack struct {
	MusicBytes     *[]byte
	ReadCloser     io.ReadCloser
	BitRate        int
	URL            string
	Name           string
	Playback       *time.Duration
	Duration       time.Duration
	RawSavedFrames []RawSavedFrame
	ReaderRead     bool
}

func init()  {
	ticker := time.NewTicker(time.Second * 5)

	go func() {
		for {
			select {
				case <-ticker.C:
					for _, musicState := range GetVoiceStates() {
						if !musicState.Running && time.Since(musicState.LastPlay) >= time.Second * 10 {
							if musicState.Voice == nil {
								continue
							}

							err := musicState.Voice.Close()
							if err != nil {
								logrus.WithFields(logrus.Fields{
									"guild-id": musicState.GuildId,
									"music-state": musicState,
								}).Error("Failed to play songs", err)
							}
							DeleteGuildMusicState(musicState.GuildId)
						}
					}
			}
		}
	}()
}

func VoiceStateUpdate(session disgord.Session, event *disgord.VoiceStateUpdate) {
	voiceState := event.VoiceState

	if voiceState.ChannelID.IsZero() {
		for _, oldVoiceState := range voiceStates {
			if oldVoiceState.UserID != voiceState.UserID && oldVoiceState.GuildID != voiceState.GuildID {
				continue
			}

			delete(voiceStates, voiceState.UserID)
			break
		}
		return
	}

	voiceStates[voiceState.UserID] = voiceState
}

func GetUserVoiceState(userId disgord.Snowflake) *disgord.VoiceState {
	return voiceStates[userId]
}

func GetVoiceState(guildId disgord.Snowflake) *MusicBotState {
	return musicStates[guildId]
}

func PutGuildMusicState(guildId disgord.Snowflake, state *MusicBotState) *MusicBotState {
	musicStates[guildId] = state

	return musicStates[guildId]
}

func DeleteGuildMusicState(guildId disgord.Snowflake) {
	delete(musicStates, guildId)
}

func GetVoiceStates() []*MusicBotState {
	states := make([]*MusicBotState, 0)
	for _, value := range musicStates {
		states = append(states, value)
	}
	return states
}
