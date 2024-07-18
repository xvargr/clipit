package main

import (
	"encoding/json"
	"net/http"

	"github.com/xvargr/clippit/internal/URLShortener"
)

func HomepageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.PathValue("resourceName")
	http.ServeFile(w, r, "./static/"+fileName)
}

func ShortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	originalURL := r.FormValue("url")
	shortenedURL := URLShortener.Instance().AddMapping(r, originalURL)
	json.NewEncoder(w).Encode(shortenedURL)
}

func ResolverHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	keyword := r.PathValue("keyword")
	originalURL, ok := URLShortener.Instance().ResolveShortKeyToOriginal(keyword)
	if !ok {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusSeeOther)
}
