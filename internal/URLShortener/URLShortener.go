package URLShortener

import (
	"math/rand/v2"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/xvargr/clippit/internal/fileReader"
)

var instance *URLShortener
var once sync.Once

type Vocabulary struct {
	Verb []string `json:"verb"`
	Noun []string `json:"noun"`
}

type URLShortener struct {
	urls       map[string]ClippedURL
	vocabulary Vocabulary
}

type ClippedURL struct {
	originalURL  string
	shortenedURL string // not used currently, but could have its use in the future
	createdAt    time.Time
}

func GetInstance() *URLShortener {
	once.Do(func() {
		instance = &URLShortener{
			urls: make(map[string]ClippedURL),
		}
		fileReader.Read("verbs.json", &instance.vocabulary)
		fileReader.Read("nouns.json", &instance.vocabulary)
	})

	return instance
}

func (u *URLShortener) AddURL(r *http.Request, originalURL string) string {
	keyword := generateShortURL()

	_, ok := u.urls[keyword]
	for ok {
		keyword = generateShortURL()
		_, ok = u.urls[keyword]
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	fullyQualifiedShortURL := scheme + "://" + r.Host + "/s/" + keyword

	u.urls[keyword] = ClippedURL{
		originalURL:  originalURL,
		shortenedURL: fullyQualifiedShortURL,
		createdAt:    time.Now(),
	}

	return fullyQualifiedShortURL
}

func (u *URLShortener) GetURL(clippedURL string) (ClippedURL, bool) {
	url, ok := u.urls[clippedURL]
	return url, ok
}

func generateShortURL() string {
	instance := GetInstance()
	nounList := instance.vocabulary.Noun
	verbList := instance.vocabulary.Verb

	return verbList[rand.IntN(len(verbList))] + "-" + nounList[rand.IntN(len(nounList))] + "-" + strconv.Itoa(rand.IntN(99))
}
