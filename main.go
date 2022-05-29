package main

import (
	"context"
	"encoding/xml"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/chi/v5"

	"github.com/Kichiyaki/lubimyczytacrss/internal/lubimyczytac"
)

const (
	defaultClientTimeout = 5 * time.Second
)

func main() {
	r := chi.NewRouter()
	r.Use(
		middleware.RealIP,
		middleware.RequestLogger(&middleware.DefaultLogFormatter{
			NoColor: true,
			Logger:  log.Default(),
		}),
		middleware.Recoverer,
		middleware.Heartbeat("/health"),
	)
	newHandler(lubimyczytac.NewClient(&http.Client{
		Timeout: defaultClientTimeout,
	})).register(r)

	httpSrv := &http.Server{
		Addr:              ":9234",
		Handler:           r,
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
		log.Println("httpSrv.Shutdown:", err)
	}
}

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

type handler struct {
	client *lubimyczytac.Client
}

func newHandler(client *lubimyczytac.Client) *handler {
	return &handler{client: client}
}

func (h *handler) register(r chi.Router) {
	r.Get("/api/v1/rss/author/{authorID}", h.getRSSAuthor)
}

func (h *handler) getRSSAuthor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	author, err := h.client.GetAuthor(ctx, chi.URLParamFromCtx(ctx, "authorID"))
	if err == lubimyczytac.ErrAuthorNotFound {
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
}
