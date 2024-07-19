package URLShortener

import (
	"errors"
	"log"
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
	Adjective []string `json:"adjective"`
	Noun      []string `json:"noun"`
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
		fileReader.Read("adjectives.json", &instance.vocabulary)
		fileReader.Read("nouns.json", &instance.vocabulary)
	})

	return instance
}

func (u *URLShortener) AddMapping(r *http.Request, originalURL string) string {
	resolved, exists := u.ResolveOriginalToShortKey(originalURL)
	if exists {
		u.renewMapping(resolved)
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

	log.Default().Printf("Added mapping: %s -> %s\n", shortKey, originalURL)

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

func (u *URLShortener) Prune(interval time.Duration) int {
	purged := 0

	for shortKey, short := range u.shortToLong {
		if time.Since(short.createdAt) > interval {
			u.RemoveMapping(shortKey)
			purged++
		}
	}

	return purged
}

func (u *URLShortener) renewMapping(shortKey string) error {
	short, sOk := u.shortToLong[shortKey]
	if !sOk {
		return errors.New("short key not found in mapping")
	}
	long, lOk := u.longToShort[short.content]
	if !lOk {
		return errors.New("long key not found in mapping")
	}

	short.createdAt = time.Now()
	long.createdAt = short.createdAt

	u.shortToLong[shortKey] = short
	u.longToShort[short.content] = long
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
	adjectiveList := u.vocabulary.Adjective
	return adjectiveList[rand.IntN(len(adjectiveList))] + "-" + nounList[rand.IntN(len(nounList))] + "-" + strconv.Itoa(rand.IntN(99))
}

func generateShortUrl(r *http.Request, k string) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + r.Host + "/s/" + k
}
