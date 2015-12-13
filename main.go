package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const Url = "http://rss.dw.com/xml/DKpodcast_dwn1_pt"
const Folder = "/Users/alberto/Documents/Projects/docker/gorss/audio/"

type EnclosureTag struct {
	XMLName xml.Name `xml:"enclosure"`
	Url     string   `xml:"url,attr"`
}

type Item struct {
	XMLName   xml.Name     `xml:"item"`
	Enclosure EnclosureTag `xml:"enclosure"`
}

type ChannelTag struct {
	XMLName xml.Name `xml:"channel"`
	Item    []Item   `xml:"item"`
}

type RssTag struct {
	XMLName  xml.Name     `xml:"rss"`
	Channels []ChannelTag `xml:"channel"`
}

func GetEpisodeList() []byte {
	response, err := http.Get(Url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return content
}

func GetFilename(url string) string {
	urlSplit := strings.Split(url, "/")
	return urlSplit[len(urlSplit)-1]
}

func GetEpisode(url string, downloadChannel chan string) {
	filename := GetFilename(url)

	fmt.Printf("Getting episode %s \n", filename)

	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	err = WriteFile(filename, content)
	if err != nil {
		fmt.Println(err)
	}

	downloadChannel <- fmt.Sprintf("Downloaded %s \n", filename)
}

func WriteFile(filename string, content []byte) error {
	err := ioutil.WriteFile(Folder + filename, content, 0444)
	if err != nil {
		fmt.Println(err)
	}

	return err
}

func HandleDownload() {
	r := GetEpisodeList()

	downloadChannel := make(chan string)

	var rssTag RssTag

	xml.Unmarshal(r, &rssTag)

	counterEpisodes := 0
	for _, channel := range rssTag.Channels {
		for _, item := range channel.Item {
			go GetEpisode(item.Enclosure.Url, downloadChannel)
			counterEpisodes++
		}
	}

	for i := 0; i < counterEpisodes; i++ {
		fmt.Printf("%s", <-downloadChannel)
	}
}

func main() {
	HandleDownload()
}
