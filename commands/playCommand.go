package commands

import (
	"RocketDiscord/cache"
	"fmt"
	"github.com/andersfylling/disgord"
	"github.com/jonas747/dca"
	"github.com/kkdai/youtube/v2"
	"github.com/sirupsen/logrus"
	"io"
	"strings"
	"time"
)

func init() {
	command := newCommand("play", false, false, onPlayCommand)
	command.setArgumentsRequirements(true, 1, -1)
	command.setHelpMessage("Polecenia pozwalajace na dolaczenie bota do kanalu glosowego przez ktorego mozna puscic muzykne na kanale. Użycie: /play <link z youtube lub innego serwisu z muzyką>")
	command.register()

	ticker := time.NewTicker(time.Second * 5)

	go func() {
		for {
			select {
				case <-ticker.C:
					for _, musicState := range cache.GetVoiceStates() {
						 if !musicState.Running && time.Since(musicState.LastPlay) >= time.Second * 10 {
						 	err := musicState.Voice.Close()
						 	if err != nil {
								logrus.WithFields(logrus.Fields{
									"guild-id": musicState.GuildId,
									"music-state": musicState,
								}).Error("Failed to play songs", err)
							}
						 	cache.DeleteGuildMusicState(musicState.GuildId)
						 }
					}
				}
			}
	}()
}

func onPlayCommand(session disgord.Session, event *disgord.MessageCreate, guild *disgord.Guild, args []string) (err error) {
	member := event.Message.Member
	musicState := cache.GetVoiceState(guild.ID)

	if musicState == nil {
		voiceState := cache.GetUserVoiceState(member.UserID)
		if voiceState == nil {
			_, err := session.SendMsg(event.Message.ChannelID, ":cloud_tornado: Abym mógł puścić dla ciebie muzykę, musisz znajdować się na kanale głosowym. (Jeżeli jesteś na kanale głosowym poprostu wyjdz i wejdz ponownie na kanał)")
			return err
		}
		voiceChannelId := voiceState.ChannelID

		if voiceChannelId.IsZero() {
			_, err := session.SendMsg(event.Message.ChannelID, ":cloud_tornado: Abym mógł puścić dla ciebie muzykę, musisz znajdować się na kanale głosowym. (Jeżeli jesteś na kanale głosowym poprostu wyjdz i wejdz ponownie na kanał)")
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

	track, err := getYoutubeFile(args)
	if err != nil || track == nil {
		_, err := session.SendMsg(event.Message.ChannelID, ":sob: Nie znaleziono takiego utworu video/audio.")
		return err
	}

	musicState.Tracks[len(musicState.Tracks)] = track

	if len(musicState.Tracks) == 1 {
		_, err = session.SendMsg(event.Message.ChannelID, ":beers: No i jest grane.")
		if err != nil {
			return err
		}

		err = startPlaying(session, musicState, event.Message.ChannelID)
		if err != nil {
			return err
		}
	} else {
		_, err = session.SendMsg(event.Message.ChannelID, ":tickets: Dodano do kolejki :D.")
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
		err = playDCA(musicState.GuildId, session, message)
		if err != nil {
			_, err = session.SendMsg(channelId, ":boom: Oj cos poszlo nie tak, nie mogę puścić muzyki :(.")
			logrus.WithFields(logrus.Fields{
				"guild-id": message.GuildID,
				"music-state": musicState,
			}).Error("Failed to play songs", err)
		}
		musicState.Running = false
	}()

	return nil
}

func getYoutubeFile(args []string) (track *cache.MusicBotTrack, err error) {
	ytClient := youtube.Client{}

	var video *youtube.Video
	if len(args) == 1 && (strings.Contains(strings.ToLower(args[0]), "youtube.com") || strings.Contains(strings.ToLower(args[0]), "youtu.be")) {
		video, err = ytClient.GetVideo(args[0])
		if err != nil {
			return
		}
	}
	if len(args) > 1 || video == nil {
		url := "https://www.youtube.com/results?q=" + strings.Join(args, "+") + "&page=1"
		video, err = ytClient.GetVideo(url)
		if err != nil {
			return
		}
	}

	if video == nil {
		return
	}

	var videoFormat *youtube.Format
	for _, format := range video.Formats {
		switch format.ItagNo {
		case 43:
		case 44:
		case 45:
		case 46:
		case 171:
		case 249:
		case 250:
		case 251:
			videoFormat = &format
			break
		}
		if videoFormat != nil {
			break
		}
	}

	if videoFormat == nil {
		return
	}

	result, err := ytClient.GetStream(video, videoFormat)
	if err != nil {
		return
	}

	track = &cache.MusicBotTrack{
		Stream:   result.Body,
		Playback: nil,
		Name:     video.Title,
		Duration: &video.Duration,
		YtFormat: videoFormat,
	}

	return
}

type SavedFrame struct {
	Payload  []byte
	Duration time.Duration
}

func playDCA(guildID disgord.Snowflake, session disgord.Session, message *disgord.Message) (err error) {
	opts := dca.StdEncodeOptions
	opts.RawOutput = true

	musicState := cache.GetVoiceState(guildID)
	if musicState == nil {
		return fmt.Errorf("music state is nil")
	}

	trackId := 0
	for {
		track := musicState.Tracks[trackId]
		if track == nil {
			if musicState.LoopPlayList {
				continue
			}
			break
		}

		_, err = session.SendMsg(message.ChannelID, fmt.Sprintf(":headphones: %v Teraz gram: %s", trackId, track.Name))
		if err != nil {
			return err
		}

		encodeSession, err := dca.EncodeMem(track.Stream, opts)
		if err != nil {
			return err
		}

		decoder := dca.NewDecoder(encodeSession)

		// inputReader is an io.Reader
		err = cache.GetVoiceState(guildID).Voice.StartSpeaking()
		if err != nil {
			return err
		}

		musicState.Running = true

		ticker := time.NewTicker(time.Second * 5)

		go func() {
			for {
				select {
				case <-ticker.C:
					stats := encodeSession.Stats()
					_, err = session.Channel(message.ChannelID).Message(message.ID).SetContent(fmt.Sprintf(":scientist: Playback: %10s, Size: %5dkB, Bitrate: %6.2fkB, Speed: %5.1fx\r", stats.Duration.String(), stats.Size, stats.Bitrate, stats.Speed));
					if err != nil {
						return
					}
				}
			}
		}()

		savedFrames := make([]SavedFrame, 0)

		for {
			musicState.LastPlay = time.Now()
			track.Playback = &encodeSession.Stats().Duration
			if track.Playback.Nanoseconds() >= track.Duration.Nanoseconds() || !musicState.Running {
				break
			}
			if musicState.Paused {
				continue
			}

			frame, err := decoder.OpusFrame()
			if err != nil {
				if err != io.EOF {
					return err
				}
				break
			}

			savedFrames = append(savedFrames, SavedFrame{
				Payload:  frame,
				Duration: decoder.FrameDuration(),
			})
			err = musicState.Voice.SendOpusFrame(frame)
			if err != nil {
				return err
			}
		}

		ticker.Stop()
		encodeSession.Truncate()
		musicState.Running = false
		time.Sleep(time.Second * 2)
		if musicState.LoopTrack {
			musicState.Running = true
			for {
				for _, savedFrame := range savedFrames {
					if !musicState.Running {
						break
					}

					if musicState.Paused {
						for {
							if !musicState.Paused {
								break
							}
						}
					}

					err = musicState.Voice.SendOpusFrame(savedFrame.Payload)
					if err != nil {
						return err
					}
				}
				if !musicState.LoopTrack {
					musicState.Running = false
					break
				}
			}
			continue
		}
		trackId++
	}

	return
}
