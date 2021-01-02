package utils

import (
	"RocketDiscord/cache"
	"fmt"
	"github.com/andersfylling/disgord"
	"github.com/jonas747/dca"
	"io"
	"time"
)

type RawSavedFrame struct {
	Payload  []byte
	Duration time.Duration
}

func PlayDCAAudio(guildID disgord.Snowflake, session disgord.Session, message *disgord.Message) (err error) {
	opts := dca.StdEncodeOptions
	opts.VBR = true
	//opts.RawOutput = true

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

		encodeSession.Stats()

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

		savedFrames := make([]RawSavedFrame, 0)

		for {
			musicState.LastPlay = time.Now()
			track.Playback = &encodeSession.Stats().Duration
			if !musicState.Running {
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

			savedFrames = append(savedFrames, RawSavedFrame{
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
