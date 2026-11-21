package app

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/egoist/mygo"

	"codie-bibel/internal/device"
	"codie-bibel/internal/i18n"
)

// settings are the reader's choices, kept in settings.json in the user
// data directory.
type settings struct {
	// Language is what the app speaks: one of i18n.Languages' codes. New
	// installs read English; an old file may still say i18n.System, which
	// follows the system's.
	Language string `json:"language"`
	// Appearance is light, dark, or system to follow the system.
	Appearance string `json:"appearance"`
	// Volume is where the hymns' volume was left, from 0 to 1.
	Volume float64 `json:"volume"`
	// FontSize is how big the reader's text is, in points; 0 is its default.
	FontSize float64 `json:"fontSize"`
	// Serif sets the reader in a serif face, as print Bibles read.
	Serif bool `json:"serif"`
	// Paragraph lays the chapter out as one flowing paragraph, rather
	// than a row a verse.
	Paragraph bool `json:"paragraph"`
	// Translation is the ID of the Bible translation last read.
	Translation string `json:"translation,omitempty"`
	// Codie is the paired device, when there is one.
	Codie *device.Creds `json:"codie,omitempty"`
	// CodiePersona is the code of the character the device speaks as,
	// as "jesus" or "monk".
	CodiePersona string `json:"codiePersona,omitempty"`
	// CodieVoice is the cloned voice each character wears, by its code.
	CodieVoice map[string]string `json:"codieVoice,omitempty"`
	// CodieSeedVoice says the app's own blessing has been cloned as the
	// Jesus persona's voice — the app clones it once, at the first
	// pairing, and never asks again.
	CodieSeedVoice bool `json:"codieSeedVoice,omitempty"`
	// CodieSetupDone says the pairing guide has been walked (or skipped).
	CodieSetupDone bool `json:"codieSetupDone,omitempty"`
	// MusicAuto plays the hymns when the app opens, as the reader turns
	// it on.
	MusicAuto bool `json:"musicAuto"`
	// SeasonReminders marks the app's keeping of the church seasons: the
	// feast card shows, and the tray counts down to the next feast.
	SeasonReminders bool `json:"seasonReminders"`

	path string // where they were loaded from, and are saved to
}

func defaultSettings() settings {
	return settings{Language: i18n.English, Appearance: "system", Volume: 0.8, SeasonReminders: true}
}

func settingsPath() (string, error) {
	dir, err := mygo.App.Path(mygo.PathUserData)
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "settings.json"), nil
}

// loadSettings reads the settings, or returns the defaults.
func loadSettings() settings {
	s := defaultSettings()
	if path, err := settingsPath(); err == nil {
		if b, err := os.ReadFile(path); err == nil {
			json.Unmarshal(b, &s)
		}
		s.path = path
	}
	return s
}

// save keeps the settings where they were loaded from; defaults that were
// never loaded, as in tests, are not saved.
func (s settings) save() error {
	path := s.path
	if path == "" {
		return nil
	}
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o600)
}

// applyAppearance makes the windows light, dark, or follow the system. Only
// a new one: it repaints every window, waiting for the main thread, and any
// setting changed applies it.
func (s settings) applyAppearance() {
	src := mygo.ThemeSystem
	switch s.Appearance {
	case "light":
		src = mygo.ThemeLight
	case "dark":
		src = mygo.ThemeDark
	}
	if src != mygo.Theme.Source() {
		mygo.Theme.SetSource(src)
	}
}

// systemLanguages are the languages the system prefers; tests replace it.
var systemLanguages = func() []string { return mygo.App.PreferredLanguages() }
