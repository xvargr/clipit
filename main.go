package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/xvargr/clippit/internal/URLShortener"
)

func main() {
	env, envErr := godotenv.Read()
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}

	port, portOk := env["PORT"]
	if !portOk {
		log.Fatal("Error loading port from .env file")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", homepageHandler)
	mux.HandleFunc("/static/{resourceName}", staticHandler)
	mux.HandleFunc("/shorten", shortenHandler)

	log.Default().Println("Starting server on port: ", port)
	if error := http.ListenAndServe(":"+port, mux); errors.Is(error, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if error != nil {
		log.Fatal("Error starting server: ", error)
		os.Exit(1)
	}
}

func homepageHandler(w http.ResponseWriter, r *http.Request) {
	handleGet := func() {
		http.ServeFile(w, r, "./static/index.html")
	}

	handlePost := func() {
		originalURL := r.FormValue("url")
		shortenedURL := URLShortener.GetInstance().AddURL(r, originalURL)
		fmt.Println("URL added successfully!", originalURL, shortenedURL)
	}

	switch r.Method {
	case http.MethodGet:
		handleGet()
	case http.MethodPost:
		handlePost()
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.PathValue("resourceName")
	http.ServeFile(w, r, "./static/"+fileName)
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	originalURL := r.FormValue("url")
	shortenedURL := URLShortener.GetInstance().AddURL(r, originalURL)
	json.NewEncoder(w).Encode(shortenedURL)
}
