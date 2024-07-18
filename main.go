package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/xvargr/clippit/internal/config"
	"github.com/xvargr/clippit/internal/scheduler"
)

func main() {
	port := config.GetConfig().Port
	purgeInterval := config.GetConfig().PruneIntervalHour

	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", HomepageHandler)
	mux.HandleFunc("GET /static/{resourceName}", StaticHandler)
	mux.HandleFunc("POST /shorten", ShortenHandler)
	mux.HandleFunc("GET /s/{keyword}", ResolverHandler)

	scheduler.Register(purgeInterval, PruneTask)

	log.Default().Println("Starting server on port: ", port)
	if error := http.ListenAndServe(":"+port, mux); errors.Is(error, http.ErrServerClosed) {
		log.Default().Printf("server closed\n")
	} else if error != nil {
		log.Fatal("Error starting server: ", error)
		os.Exit(1)
	}
}
