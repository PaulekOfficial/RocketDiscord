package utils

import (
	"RocketDiscord/cache"
	"bytes"
	"fmt"
	"github.com/andersfylling/disgord"
	"github.com/jonas747/dca"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"sync"
	"time"
)

func PlayDCAAudio(guildID disgord.Snowflake, channelID disgord.Snowflake,session disgord.Session) (err error) {
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

		opts := dca.StdEncodeOptions
		opts.Bitrate = track.BitRate
		opts.RawOutput = true

		trueReader := false
		stream := track.ReadCloser
		if track.MusicBytes != nil {
			stream = ioutil.NopCloser(bytes.NewReader(*track.MusicBytes))
			trueReader = true
		}

		encodeSession, err := dca.EncodeMem(stream, opts)
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

		_, err = session.SendMsg(channelID, fmt.Sprintf(":headphones: Teraz gram: %s", track.Name))
		if err != nil {
			return err
		}

		ticker := time.NewTicker(time.Second * 5)

		playedBytes := 0

		//go func() {
		//	err := CreateNewDecoderIfEnviable(track.MusicBytes, playedBytes)
		//	if err != nil {
		//
		//	}
		//}()

		var mutex sync.Mutex

		for {
			musicState.LastPlay = time.Now()
			track.Playback = &encodeSession.Stats().Duration
			if !musicState.Running {
				break
			}
			if musicState.Paused {
				continue
			}

			if !trueReader && track.MusicBytes != nil {
				trueReader = true
				go func() {
					dec, index, err := CreateNewDecoderIfEnviable(track.MusicBytes, playedBytes, opts)
					if err != nil {
						logrus.WithFields(logrus.Fields{
							"server-id": guildID,
						}).Error("Fail to inject proper decoder to playing", err)
						return
					}

					mutex.Lock()
					if index < playedBytes {
						for {
							if index > playedBytes {
								break
							}
							_, err = dec.OpusFrame()
							if err != nil {
								logrus.WithFields(logrus.Fields{
									"server-id": guildID,
								}).Error("Fail to opus frame", err)
								return
							}

							index++
						}
					}
					decoder = dec
					mutex.Unlock()
				}()
			}

			mutex.Lock()
			frame, err := decoder.OpusFrame()
			if err != nil {
				if err != io.EOF {
					return err
				}

				if !trueReader {
					decoder, _, err = CreateNewDecoderIfEnviable(track.MusicBytes, playedBytes, opts)
					if err != nil {
						return err
					}

					trueReader = true
					continue
				}

				break
			}
			playedBytes++

			if !track.ReaderRead {
				track.RawSavedFrames = append(track.RawSavedFrames, cache.RawSavedFrame{
					Payload:  frame,
					Duration: decoder.FrameDuration(),
				})
			}

			err = musicState.Voice.SendOpusFrame(frame)
			if err != nil {
				return err
			}
			mutex.Unlock()
		}

		ticker.Stop()
		encodeSession.Truncate()
		musicState.Running = false
		time.Sleep(time.Second)

		if musicState.LoopTrack {
			track.ReaderRead = true
			continue
		}

		trackId++
	}

	return
}

func CreateNewDecoderIfEnviable(dcaBytes *[]byte, playedBytes int, options *dca.EncodeOptions) (decoder *dca.Decoder, readerIndex int, err error) {
	for {
		if dcaBytes != nil {
			break
		}
	}

	if len(*dcaBytes) <= 0 {
		return nil, 0, nil
	}

	stream := ioutil.NopCloser(bytes.NewReader(*dcaBytes))
	encodeSession, err := dca.EncodeMem(stream, options)
	if err != nil {
		return nil, 0, err
	}
	decoder = dca.NewDecoder(encodeSession)
	for i := 0; i < playedBytes; i ++ {
		readerIndex = playedBytes
		_, err := decoder.OpusFrame()
		if err != nil {
			if err != io.EOF {
				return nil, 0, err
			}
			break
		}
	}

	return
}
