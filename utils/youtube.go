package utils

import (
	"RocketDiscord/cache"
	"bytes"
	"encoding/csv"
	"github.com/andersfylling/disgord"
	"github.com/kkdai/youtube"
	"github.com/lithdew/bytesutil"
	"github.com/valyala/fastjson"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Video struct {
	ID string `json:"encrypted_id"`

	Title       string `json:"title"`
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnail"`

	Added       string    `json:"added"`
	TimeCreated time.Time `json:"time_created"`

	Rating   float64 `json:"rating"`
	Likes    uint    `json:"likes"`
	Dislikes uint    `json:"dislikes"`

	Views    string `json:"views"`
	Comments string `json:"comments"`

	Duration      string        `json:"duration"`
	LengthSeconds time.Duration `json:"length_seconds"`

	Author  string `json:"author"`
	UserID  string `json:"user_id"`
	Privacy string `json:"privacy"`

	CategoryID uint `json:"category_id"`

	IsHD bool `json:"is_hd"`
	IsCC bool `json:"is_cc"`

	CCLicense bool `json:"cc_license"`

	Keywords []string `json:"keywords"`
}

type YoutubeSearchQuery struct {
	Hits  uint       `json:"hits"`
	Videos []Video `json:"video"`
}

func SearchYoutubeVideos(keywords []string, page string) (query *YoutubeSearchQuery, err error){
	uri := []byte("https://www.youtube.com/search_ajax?style=json")

	uri = append(uri, "&search_query="...)
	uri = append(uri, strings.Join(keywords, "+")...)

	uri = append(uri, "&page="...)
	uri = append(uri, page...)

	uri = append(uri, "&hl="...)
	uri = append(uri, "en"...)

	httpClient := http.Client{}

	request, err := http.NewRequest("GET", string(uri), nil)
	if err != nil {
		return
	}
	request.Header.Set("x-youtube-client-name", "56")
	request.Header.Set("x-youtube-client-version", "20200911")

	response, err := httpClient.Do(request)
	if err != nil {
		return
	}

	buffer, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	value, err := fastjson.ParseBytes(buffer)
	if err != nil {
		return
	}

	rawVideos := value.GetArray("video")
	
	query = &YoutubeSearchQuery{
		Hits:   value.GetUint("hits"),
		Videos: make([]Video, 0, len(rawVideos)),
	}

	for _, rawVideo := range rawVideos {
		query.Videos = append(query.Videos, ParseRawVideo(rawVideo))
	}

	return
}

func GetYoutubeStreamMemory(args []string, session disgord.Session, channelID disgord.Snowflake) (track *cache.MusicBotTrack, err error) {
	ytClient := youtube.Client{
		Debug:      true,
		HTTPClient: http.DefaultClient,
	}

	var video *youtube.Video
	if len(args) == 1 && (strings.Contains(strings.ToLower(args[0]), "youtube.com") || strings.Contains(strings.ToLower(args[0]), "youtu.be")) {
		//videoId := GetVideoIdFromLink(args[0])
		video, err = ytClient.GetVideo(args[0])
		if err != nil {
			return nil, err
		}

		_, err = session.SendMsg(channelID, ":telescope: Ładuje podany utwór...")
		if err != nil {
			return nil, err
		}
	}
	if len(args) > 1 || video != nil {
		_, err = session.SendMsg(channelID, ":hourglass_flowing_sand: Wyszukuje dopasowań w serwisie youtube...")
		if err != nil {
			return nil, err
		}

		videos, err := SearchYoutubeVideos(args, "0")
		if err != nil {
			return nil, err
		}

		if videos == nil || videos.Hits == 0 && video != nil {
			return nil, err
		}

		_, err = session.SendMsg(channelID, ":telescope: Znaleziono " + strconv.Itoa(int(videos.Hits)) + " dopasowań, wybieram najlepsze...")
		if err != nil {
			return nil, err
		}

		video, err =  ytClient.GetVideo("https://www.youtube.com/watch?v=" + videos.Videos[0].ID)
		if err != nil {
			return nil, err
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
	if err != nil || result == nil {
		return
	}

	buffer, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return
	}

	//reader := ioutil.NopCloser(bytes.NewReader(buffer))

	track = &cache.MusicBotTrack{
		Stream:     bytes.NewReader(buffer),
		MusicBytes: buffer,
		URL:        "https://www.youtube.com/watch?v=" + video.ID,
		Playback:   nil,
		Name:       video.Title,
		Duration:   video.Duration,
	}

	return
}

func GetVideoIdFromLink(link string) string {
	rawIds := strings.Split(link, "?v=")
	if len(rawIds) != 2 {
		return ""
	}

	if strings.Contains(rawIds[1], "&") {
		idAndList := strings.Split(rawIds[1], "&")
		if len(idAndList) == 2 {
			return idAndList[0]
		}
	}

	return rawIds[1]
}

func ParseRawVideo(v *fastjson.Value) Video {
	var video Video

	video.ID = string(v.GetStringBytes("encrypted_id"))

	video.Title = bytesutil.String(v.GetStringBytes("title"))
	video.Description = bytesutil.String(v.GetStringBytes("description"))
	video.Thumbnail = bytesutil.String(v.GetStringBytes("thumbnail"))

	video.Added = bytesutil.String(v.GetStringBytes("added"))
	video.TimeCreated = time.Unix(v.GetInt64("time_created"), 0)

	video.Rating = v.GetFloat64("rating")
	video.Likes = v.GetUint("likes")
	video.Dislikes = v.GetUint("dislikes")

	video.Views = bytesutil.String(v.GetStringBytes("views"))
	video.Comments = bytesutil.String(v.GetStringBytes("comments"))

	video.Duration = bytesutil.String(v.GetStringBytes("duration"))
	video.LengthSeconds = time.Duration(v.GetInt64("length_seconds")) * time.Second

	video.Author = bytesutil.String(v.GetStringBytes("author"))
	video.UserID = bytesutil.String(v.GetStringBytes("user_id"))
	video.Privacy = bytesutil.String(v.GetStringBytes("privacy"))

	video.CategoryID = v.GetUint("category_id")

	video.IsHD = v.GetBool("is_hd")
	video.IsCC = v.GetBool("is_cc")

	video.CCLicense = v.GetBool("cc_license")

	fr := csv.NewReader(bytes.NewReader(v.GetStringBytes("keywords")))
	fr.Comma = ' '

	video.Keywords, _ = fr.Read()

	return video
}
