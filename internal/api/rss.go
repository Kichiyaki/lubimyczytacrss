package api

import (
	"encoding/xml"

	"github.com/Kichiyaki/lubimyczytacrss/internal/lubimyczytac"
)

type rssItem struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	GUID        string   `xml:"guid"`
	Description string   `xml:"description"`
}

func rssItemsFromBooks(books []lubimyczytac.Book) []rssItem {
	items := make([]rssItem, len(books))
	for i, b := range books {
		items[i] = rssItem{
			Title:       b.Title,
			Link:        b.URL,
			GUID:        b.URL,
			Description: "",
		}
	}
	return items
}

type rssChannel struct {
	XMLName     xml.Name  `xml:"channel"`
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Language    string    `xml:"language"`
	Items       []rssItem `xml:"items"`
}

func rssChannelFromAuthor(author lubimyczytac.Author) rssChannel {
	return rssChannel{
		Title:       author.Name,
		Description: author.ShortDescription,
		Link:        author.URL,
		Items:       rssItemsFromBooks(author.Books),
	}
}

type rssMain struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel rssChannel `xml:"channel"`
}

func rssMainFromAuthor(author lubimyczytac.Author) rssMain {
	return rssMain{
		Version: "2.0",
		Channel: rssChannelFromAuthor(author),
	}
}
