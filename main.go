package main

import (
	"context"
	"encoding/xml"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Kichiyaki/lubimyczytacrss/internal"
)

const (
	defaultLubimyCzytacClientTimeout = 5 * time.Second
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

type rssItem struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	GUID        string   `xml:"guid"`
	Description string   `xml:"description"`
}

type rssChannel struct {
	XMLName     xml.Name  `xml:"channel"`
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Language    string    `xml:"language"`
	Items       []rssItem `xml:"items"`
}

type rssMain struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel rssChannel `xml:"channel"`
}

func rssMainFromAuthor(author internal.Author) rssMain {
	items := make([]rssItem, len(author.Books))
	for i, b := range author.Books {
		items[i] = rssItem{
			Title:       b.Title,
			Link:        b.URL,
			GUID:        b.URL,
			Description: "",
		}
	}
	return rssMain{
		Version: "2.0",
		Channel: rssChannel{
			Title:       author.Name,
			Description: author.ShortDescription,
			Link:        author.URL,
			Items:       items,
		},
	}
}

func newHandler() http.Handler {
	client := internal.NewLubimyCzytacClient(&http.Client{
		Timeout: defaultLubimyCzytacClientTimeout,
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		author, err := client.GetAuthor(r.Context(), r.URL.Path[1:])
		if err == internal.ErrAuthorNotFound {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`author not found`))
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`something went wrong while getting author info: ` + err.Error()))
			return
		}

		w.Header().Set("Content-Type", "text/xml; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = xml.NewEncoder(w).Encode(rssMainFromAuthor(author))
	})
}
