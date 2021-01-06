package utils

import (
	"context"
	"encoding/xml"
	"net/http"
	"strconv"
	"strings"
)

var rule34Url string = "http://rule34.paheal.net/api/danbooru/find_posts/index.xml"

var httpClient *http.Client = http.DefaultClient

type Rule34Posts struct {
	XMLName   xml.Name `xml:"posts"`
	Count     string   `xml:"count,attr"`
	Offset    string   `xml:"offset,attr"`
	Tags      []Rule34Tag `xml:"tag"`
}

type Rule34Tag struct {
	XMLName xml.Name `xml:"tag"`
	ID    string   `xml:"id,attr"`
	MD5    string   `xml:"md5,attr"`
	FileName    string   `xml:"file_name,attr"`
	FileUrl    string   `xml:"file_url,attr"`
	Height    string   `xml:"height,attr"`
	Width    string   `xml:"width,attr"`
	PreviewHeight    string   `xml:"preview_height,attr"`
	PreviewWidth    string   `xml:"preview_width,attr"`
	Rating    string   `xml:"rating,attr"`
	Date    string   `xml:"date,attr"`
	Tags    string   `xml:"tags,attr"`
	Source    string   `xml:"source,attr"`
	Score    string   `xml:"score,attr"`
	Author    string   `xml:"author,attr"`
}

func GetRule34UrlImageFromTags(tags []string, limit int) (urls []string, err error) {
	pureUrl := pureRule34Url(tags, limit)

	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, pureUrl, nil)
	if err != nil {
		return
	}

	response, err := httpClient.Do(request)
	if err != nil {
		return
	}

	defer response.Body.Close()

	decoder := xml.NewDecoder(response.Body)

	var posts *Rule34Posts

	err = decoder.Decode(&posts)
	if err != nil {
		return
	}

	if posts == nil || len(posts.Tags) == 0 {
		return
	}

	urls = make([]string, 0)
	for _, tag := range posts.Tags {
		urls = append(urls, tag.FileUrl)
	}

	return
}

func pureRule34Url(tags []string, limit int) string {
	pureUrl := []byte(rule34Url)

	pureUrl = append(pureUrl, "?tags="...)
	pureUrl = append(pureUrl, strings.Join(tags, "%20")...)
	pureUrl = append(pureUrl, "&limit="...)
	pureUrl = append(pureUrl, strconv.Itoa(limit)...)

	return string(pureUrl)
}
