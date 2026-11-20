// Package devotion is the devotions the app keeps: Spurgeon's Morning and
// Evening, the verse of the day, the reading plans, and the scriptures a
// prayer walks through. Spurgeon's text is data, generated into data/ by
// cmd/genbible; the rest are tables of references into the Bible.
package devotion

import (
	_ "embed"
	"encoding/json"
	"strconv"
	"time"
)

//go:embed data/spurgeon.json
var spurgeonData []byte

// Reading is one of Morning and Evening's readings: its verse, and what
// Spurgeon wrote of it.
type Reading struct {
	Verse string
	Text  string
}

// MorningEvening is the day's two readings.
type MorningEvening struct {
	AM Reading
	PM Reading
}

// Spurgeon is Morning and Evening, a reading for each part of the day.
type Spurgeon struct {
	readings map[string]map[string]Reading
}

// NewSpurgeon reads the generated table.
func NewSpurgeon() (*Spurgeon, error) {
	var raw struct {
		Readings map[string]map[string]struct {
			Verse string `json:"verse"`
			Text  string `json:"text"`
		} `json:"readings"`
	}
	if err := json.Unmarshal(spurgeonData, &raw); err != nil {
		return nil, err
	}
	s := &Spurgeon{readings: map[string]map[string]Reading{}}
	for day, slots := range raw.Readings {
		m := map[string]Reading{}
		for slot, r := range slots {
			m[slot] = Reading{Verse: r.Verse, Text: r.Text}
		}
		s.readings[day] = m
	}
	return s, nil
}

// Read is the day's morning and evening. A day the table doesn't keep
// (a February 30th) reads as empty.
func (s *Spurgeon) Read(t time.Time) MorningEvening {
	key := strconv.Itoa(int(t.Month())) + "-" + strconv.Itoa(t.Day())
	slots := s.readings[key]
	return MorningEvening{AM: slots["am"], PM: slots["pm"]}
}

// Half says which of the day's readings the hour asks for: the morning's
// before noon, the evening's after.
func Half(t time.Time) string {
	if t.Hour() < 12 {
		return "am"
	}
	return "pm"
}
