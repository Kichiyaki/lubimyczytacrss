package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	httpSrv := &http.Server{
		Addr:              ":9234",
		Handler:           nil,
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
