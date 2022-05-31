package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Kichiyaki/lubimyczytacrss/internal/api"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/chi/v5"

	"github.com/Kichiyaki/lubimyczytacrss/internal/lubimyczytac"
)

const (
	clientTimeout = 5 * time.Second
)

func main() {
	httpSrv := newServer(newRouter(lubimyczytac.NewClient(&http.Client{
		Timeout: clientTimeout,
	})))

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

func newServer(h http.Handler) *http.Server {
	return &http.Server{
		Addr:              ":9234",
		Handler:           h,
		ReadTimeout:       2 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      2 * time.Second,
		IdleTimeout:       2 * time.Second,
	}
}

func newRouter(client *lubimyczytac.Client) *chi.Mux {
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
	api.NewHandler(client).Register(r)
	return r
}
