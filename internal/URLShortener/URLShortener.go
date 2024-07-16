package URLShortener

import (
	"errors"
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
	shortToLong map[string]UrlStore
	longToShort map[string]UrlStore
	vocabulary  Vocabulary
}

type UrlStore struct {
	content   string
	createdAt time.Time
}

func Instance() *URLShortener {
	once.Do(func() {
		instance = &URLShortener{
			shortToLong: make(map[string]UrlStore),
			longToShort: make(map[string]UrlStore),
		}
		fileReader.Read("verbs.json", &instance.vocabulary)
		fileReader.Read("nouns.json", &instance.vocabulary)
	})

	return instance
}

func (u *URLShortener) AddMapping(r *http.Request, originalURL string) string {
	resolved, exists := u.ResolveOriginalToShortKey(originalURL)
	if exists {
		return resolved
	}

	// regenerate short key if there are duplicates
	shortKey := u.generateShortKey()
	_, ok := u.shortToLong[shortKey]
	for ok {
		shortKey = u.generateShortKey()
		_, ok = u.shortToLong[shortKey]
	}

	long := UrlStore{
		content:   generateShortUrl(r, shortKey),
		createdAt: time.Now(),
	}
	short := UrlStore{
		content:   originalURL,
		createdAt: long.createdAt,
	}

	u.shortToLong[shortKey] = short
	u.longToShort[originalURL] = long

	return long.content
}

func (u *URLShortener) RemoveMapping(shortKey string) error {
	short, ok := u.shortToLong[shortKey]
	if !ok {
		return errors.New("short key not found in mapping")
	}

	delete(u.shortToLong, shortKey)
	delete(u.longToShort, short.content)
	return nil
}

func (u *URLShortener) ResolveShortKeyToOriginal(shortKey string) (string, bool) {
	short, ok := u.shortToLong[shortKey]
	return short.content, ok
}

func (u *URLShortener) ResolveOriginalToShortKey(originalURL string) (string, bool) {
	long, ok := u.longToShort[originalURL]
	return long.content, ok
}

func (u *URLShortener) generateShortKey() string {
	nounList := u.vocabulary.Noun
	verbList := u.vocabulary.Verb
	return verbList[rand.IntN(len(verbList))] + "-" + nounList[rand.IntN(len(nounList))] + "-" + strconv.Itoa(rand.IntN(99))
}

func generateShortUrl(r *http.Request, k string) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + r.Host + "/s/" + k
}
