package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"
)

const (
	lubimyCzytacDomain = "lubimyczytac.pl"
	defaultTimeout     = 5 * time.Second
)

func main() {
	httpSrv := &http.Server{
		Addr:              ":9234",
		Handler:           newHandler(),
		ReadTimeout:       2 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      2 * time.Second,
		IdleTimeout:       2 * time.Second,
	}

	go func(httpSrv *http.Server) {
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln("httpSrv.ListenAndServe:", err)
		}
	}(httpSrv)

	log.Println("Server is listening on the port 9234")

	ctxSignal, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	<-ctxSignal.Done()

	ctxShutdown, cancelCtxShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelCtxShutdown()
	if err := httpSrv.Shutdown(ctxShutdown); err != nil {
		log.Fatalln("httpSrv.Shutdown:", err)
	}
}

type rssMain struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	XMLName     xml.Name  `xml:"channel"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Language    string    `xml:"language"`
	Items       []rssItem `xml:"items"`
}

type rssItem struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	GUID        string   `xml:"guid"`
	Description string   `xml:"description"`
}

func newHandler() http.Handler {
	client := &http.Client{
		Timeout: defaultTimeout,
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, r.URL.Path[1:], nil)
		if err != nil || req.URL.Host != lubimyCzytacDomain {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer func() {
			_ = resp.Body.Close()
		}()

		p, err := newAuthorPageParser(resp.Body)
		if err != nil {
			_, _ = io.Copy(io.Discard, resp.Body)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/xml; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if err := xml.NewEncoder(w).Encode(rssMain{
			Version: "2.0",
			Channel: rssChannel{
				Link:  req.URL.String(),
				Items: p.items(req.URL),
			},
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

type authorPageParser struct {
	doc *goquery.Document
}

func newAuthorPageParser(r io.Reader) (authorPageParser, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return authorPageParser{}, fmt.Errorf("goquery.NewDocumentFromReader: %w", err)
	}
	return authorPageParser{
		doc: doc,
	}, nil
}

func (p authorPageParser) items(baseURL *url.URL) []rssItem {
	booksSel := p.doc.Find("#authorBooks .authorAllBooks__single")
	items := make([]rssItem, booksSel.Length())
	booksSel.Each(func(i int, selection *goquery.Selection) {
		link := url.URL{
			Scheme: baseURL.Scheme,
			Host:   baseURL.Host,
			Path:   selection.Find(".authorAllBooks__singleTextTitle").AttrOr("href", ""),
		}
		linkStr := link.String()
		title := strings.TrimSpace(selection.Find(".authorAllBooks__singleTextTitle").Text())
		items[i] = rssItem{
			Title:       title,
			Link:        linkStr,
			GUID:        linkStr,
			Description: "",
		}
	})
	return items
}
