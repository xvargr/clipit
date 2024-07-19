package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/xvargr/clippit/internal/URLShortener"
	"github.com/xvargr/clippit/internal/config"
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
	if originalURL == "" {
		http.Error(w, "Incomplete form data, URL is required", http.StatusBadRequest)
		return
	}

	if strings.Contains(originalURL, r.Host) {
		http.Error(w, "Self referencing URLs are not allowed", http.StatusBadRequest)
		return
	}

	shortenedURL := URLShortener.Instance().AddMapping(r, originalURL)
	json.NewEncoder(w).Encode((map[string]string{"shortenedUrl": shortenedURL, "validity": config.GetConfig().PruneIntervalHour.String()}))
}

func ResolverHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	keyword := r.PathValue("keyword")
	originalURL, ok := URLShortener.Instance().ResolveShortKeyToOriginal(keyword)
	if !ok {
		http.Error(w, "Not a valid short URL", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusSeeOther)
}
