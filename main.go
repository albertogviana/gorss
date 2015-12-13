package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"flag"
)

const URL = "http://rss.dw.com/xml/DKpodcast_dwn%d_pt"

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

func GetEpisodeList(season int) []byte {
	url := fmt.Sprintf(URL, season)
	response, err := http.Get(url)
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
	return urlSplit[len(urlSplit) - 1]
}

func GetEpisode(url string, directory string, downloadChannel chan string) {
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

	err = WriteFile(directory + "/" + filename, content)
	if err != nil {
		fmt.Println(err)
	}

	downloadChannel <- fmt.Sprintf("Downloaded %s \n", filename)
}

func WriteFile(filename string, content []byte) error {
	err := ioutil.WriteFile(filename, content, 0444)
	if err != nil {
		fmt.Println(err)
	}

	return err
}

func HandleDownload(season int, directory string) {
	r := GetEpisodeList(season)

	downloadChannel := make(chan string)

	var rssTag RssTag

	xml.Unmarshal(r, &rssTag)

	counterEpisodes := 0
	for _, channel := range rssTag.Channels {
		for _, item := range channel.Item {
			go GetEpisode(item.Enclosure.Url, directory, downloadChannel)
			counterEpisodes++
		}
	}

	for i := 0; i < counterEpisodes; i++ {
		fmt.Printf("%s", <-downloadChannel)
	}
}

func Run() {
	var seasonNumber = flag.Int("season", 0, "Deutsch warum nicht serie. There are 4 season")
	var directory = flag.String("save", "", "Folder where the audio will be saved")
	flag.Parse()

	if *seasonNumber < 1 || *seasonNumber > 4 {
		fmt.Println("The season number should be between 1 and 4 \n")
		os.Exit(0)
	}

	_, err := os.Stat(*directory)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	HandleDownload(*seasonNumber, *directory)
}

func main() {
	Run()
}
