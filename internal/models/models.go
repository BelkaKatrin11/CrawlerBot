package models

import "encoding/xml"

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel struct {
		XMLName xml.Name `xml:"channel"`
		Title string `xml:"title"`
		Description string `xml:"description"`
		Item []Item `xml:"item"`
	} `xml:"channel"`
}


type Item struct {
	XMLName xml.Name `xml:"item"`
	Guid string `xml:"guid"`
	Title string `xml:"title"`
	Link string `xml:"link"`
	Description string `xml:"description"`
	Pubdate string `xml:"pubDate"`
	Category string `xml:"category"`
}

