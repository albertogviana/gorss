package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"encoding/xml"
)

const Url = "http://rss.dw.com/xml/DKpodcast_dwn1_pt"

type EnclosureTag struct  {
	XMLName xml.Name `xml:"enclosure"`
	Url string `xml:"url,attr"`
}

type Item struct {
	XMLName xml.Name `xml:"item"`
	Enclosure EnclosureTag `xml:"enclosure"`
}

type ChannelTag struct {
	XMLName xml.Name `xml:"channel"`
	Item []Item `xml:"item"`
}

type RssTag struct {
	XMLName xml.Name `xml:"rss"`
	Channels []ChannelTag `xml:"channel"`
}


func GetSessionList() []byte {
	response, err := http.Get(Url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(response.StatusCode)

	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return content
}

func main() {
	r := GetSessionList()

	var rssTag RssTag

	xml.Unmarshal(r, &rssTag)

	for _, channel := range rssTag.Channels {
		for _, item := range channel.Item {
			fmt.Printf("%s \n", item.Enclosure.Url)
		}
	}
}
