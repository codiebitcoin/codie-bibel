// Package hymn is the background music the app ships with: hymns and
// chorales in the public domain, embedded as their upstreams recorded
// them, played quietly while the reader reads. The audio files are data,
// generated into data/ from Wikimedia Commons (see LICENSES.md).
package hymn

import (
	"embed"
	"io"
	"os"
	"path/filepath"
	"time"
)

//go:embed data
var files embed.FS

// Track is one hymn: what it is, and what it may be played under.
type Track struct {
	ID       string        `json:"id"`
	Title    string        `json:"title"`
	Artist   string        `json:"artist"` // the performer, for the system's Now Playing
	File     string        `json:"file"`   // in data/
	Duration time.Duration `json:"duration"`
	License  string        `json:"license"`
}

// Tracks are the hymns, in the order the Music page lists them: the
// chorales first, then the carols.
var Tracks = []Track{
	{ID: "bwv147", Title: "Jesu, Joy of Man's Desiring (BWV 147)", Artist: "Orchestra Gli Armonici", File: "bwv147.ogg", Duration: 202 * time.Second, License: "CC0 1.0"},
	{ID: "bwv140", Title: "Wachet auf, ruft uns die Stimme (BWV 140)", Artist: "MIT Concert Choir, dir. William Cutter", File: "bwv140.ogg", Duration: 347 * time.Second, License: "CC BY-SA 2.0"},
	{ID: "sleepers-wake", Title: "Sleepers, Wake (brass)", Artist: "Galaxy Brass, US Air Force Band of the Rockies", File: "sleepers-wake.mp3", Duration: 280 * time.Second, License: "Public domain (US Government work)"},
	{ID: "mighty-fortress", Title: "A Mighty Fortress Is Our God", Artist: "The Army and Navy Hymnal (1922)", File: "mighty-fortress.ogg", Duration: 223 * time.Second, License: "Public domain"},
	{ID: "amazing-grace", Title: "Amazing Grace", Artist: "Original Sacred Harp Choir (1922)", File: "amazing-grace.ogg", Duration: 120 * time.Second, License: "Public domain"},
	{ID: "silent-night", Title: "Silent Night", Artist: "The US Army Chorus", File: "silent-night.ogg", Duration: 150 * time.Second, License: "Public domain (US Government work)"},
	{ID: "joy-world-army", Title: "Joy to the World", Artist: "The US Army Band Chorus", File: "joy-to-world-army.ogg", Duration: 75 * time.Second, License: "Public domain (US Government work)"},
	{ID: "joy-world-diplomats", Title: "Joy to the World (The Diplomats)", Artist: "The US Air Force Band", File: "joy-to-world-diplomats.ogg", Duration: 268 * time.Second, License: "Public domain (US Government work)"},
}

// ByID is the hymn of that ID, or nil.
func ByID(id string) *Track {
	for i, t := range Tracks {
		if t.ID == id {
			return &Tracks[i]
		}
	}
	return nil
}

// Ensure writes the hymns' audio into dir, once each, and returns the
// path of every track's file, by ID. A file that is already there and the
// size it should be stays.
func Ensure(dir string) (map[string]string, error) {
	paths := map[string]string{}
	for _, t := range Tracks {
		path := filepath.Join(dir, t.File)
		if err := writeOnce(t.File, path); err != nil {
			return nil, err
		}
		paths[t.ID] = path
	}
	return paths, nil
}

// writeOnce writes the embedded file name to path, once.
func writeOnce(name, path string) error {
	if fi, err := os.Stat(path); err == nil && fi.Size() > 0 {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := files.Open("data/" + name)
	if err != nil {
		return err
	}
	defer f.Close()
	tmp := path + ".part"
	out, err := os.Create(tmp)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, f); err != nil {
		out.Close()
		return err
	}
	if err := out.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
