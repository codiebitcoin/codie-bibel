// Package media is the handful of types the player's streams speak: what
// a track is, what its audio is, and where a file of it lies. It stands
// in for the wider music package the player once lived with, keeping the
// decoding plumbing apart from the app.
package media

import (
	"errors"
	"net/http"
	"net/url"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// ErrUnavailable is what a track that has no audio returns.
var ErrUnavailable = errors.New("track unavailable")

// Artist is one who made the track.
type Artist struct {
	Name string
}

// Album is the track's album, with its cover.
type Album struct {
	Name     string
	ImageURL string
}

// Track is one piece of music to play: here, a hymn the app ships.
type Track struct {
	ID       string
	Source   string
	Name     string
	Artists  []Artist
	Album    Album
	Duration int64 // milliseconds
}

// ArtistNames reads the track's artists, as a system shows them.
func (t Track) ArtistNames() string {
	var out []string
	for _, a := range t.Artists {
		out = append(out, a.Name)
	}
	return strings.Join(out, ", ")
}

// Length is how long the track plays.
func (t Track) Length() time.Duration {
	return time.Duration(t.Duration) * time.Millisecond
}

// Stream is a track's audio: here, a file the app has written to disk.
type Stream struct {
	URL     string
	Type    string
	Size    int64
	Gain    float64
	Peak    float64
	Headers map[string]string
}

// Agent is the header the downloads go out under.
const Agent = "CodieBible/0.1"

// HeadersStay is a client that keeps its headers to the host they were
// set for: the streams the app plays come from home, and the headers of a
// hymn's file are none at all.
func HeadersStay(client *http.Client, host string, headers map[string]string) *http.Client {
	c := *client
	c.CheckRedirect = func(next *http.Request, via []*http.Request) error {
		if len(via) >= 10 {
			return errors.New("stopped after 10 redirects")
		}
		if to := strings.ToLower(next.URL.Hostname()); to != host && !strings.HasSuffix(to, "."+host) {
			for k := range headers {
				next.Header.Del(k)
			}
		}
		return nil
	}
	return &c
}

// FileURL is the file:// URL of a path.
func FileURL(path string) string {
	path = filepath.ToSlash(path)
	if !strings.HasPrefix(path, "/") {
		path = "/" + path // C:/Music
	}
	return (&url.URL{Scheme: "file", Path: path}).String()
}

// PathOf returns the path a file:// URL names, and whether u is one.
func PathOf(u string) (string, bool) {
	if !strings.HasPrefix(u, "file:") {
		return "", false
	}
	parsed, err := url.Parse(u)
	if err != nil || parsed.Path == "" {
		return "", false
	}
	path := parsed.Path
	if runtime.GOOS == "windows" {
		path = strings.TrimPrefix(path, "/")
	}
	return filepath.FromSlash(path), true
}
