package main

// Morning and Evening: genbible turns a transcription of Spurgeon's
// devotional (the text itself is public domain, Spurgeon died in 1892;
// this one is dwahrens/morning-and-evening's) into the app's own
// spurgeon.json — the morning and evening readings of every day of the
// year.

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const spurgeonURL = "https://raw.githubusercontent.com/dwahrens/morning-and-evening/main/data/m_e.json"

// spurgeonEntry is one reading of the source file.
type spurgeonEntry struct {
	Date     string `json:"date"`
	Time     string `json:"time"` // am or pm
	Month    int    `json:"month"`
	Day      int    `json:"day"`
	KeyVerse string `json:"keyverse"`
	Body     string `json:"body"`
}

// spurgeonOut is what the app reads: each day its morning and its evening,
// the verse and the text of each.
type spurgeonOut struct {
	Readings map[string]map[string]spurgeonReading `json:"readings"` // "1-1" → "am"/"pm"
}

type spurgeonReading struct {
	Verse string `json:"verse"`
	Text  string `json:"text"`
}

var (
	spurgeonHeader = regexp.MustCompile(`(?s)^[^\n]*—\s*(Morning|Evening) Reading\s*`)
	spurgeonSpaces = regexp.MustCompile(`[ \t]+`)
	spurgeonBreaks = regexp.MustCompile(`\n{3,}`)
	spurgeonMarks  = regexp.MustCompile(`[\x{2009}\x{200e}\x{200f}]`)
)

func genSpurgeon(cacheDir, outDir string) error {
	local := filepath.Join(cacheDir, "spurgeon.json")
	if err := fetch(spurgeonURL, local); err != nil {
		return err
	}
	b, err := os.ReadFile(local)
	if err != nil {
		return err
	}
	var src []spurgeonEntry
	if err := json.Unmarshal(b, &src); err != nil {
		return err
	}
	out := spurgeonOut{Readings: map[string]map[string]spurgeonReading{}}
	for _, e := range src {
		if e.Month < 1 || e.Month > 12 || e.Day < 1 || e.Day > 31 {
			continue
		}
		day := fmt.Sprintf("%d-%d", e.Month, e.Day)
		am := e.Time != "pm"
		verse := spurgeonMarks.ReplaceAllString(e.KeyVerse, "")
		verse = strings.TrimSpace(spurgeonSpaces.ReplaceAllString(verse, " "))
		// The verse field reads '"...the text..." — Joshua 5:12': the
		// reference is what follows the dash.
		if _, ref, has := strings.Cut(verse, "—"); has {
			verse = strings.TrimSpace(ref)
		} else if _, ref, has := strings.Cut(verse, "-"); has && strings.Count(verse, "-") == 1 && strings.HasSuffix(verse, `"`) == false {
			verse = strings.TrimSpace(ref)
		}
		text := spurgeonMarks.ReplaceAllString(e.Body, "")
		text = spurgeonHeader.ReplaceAllString(text, "")
		text = strings.ReplaceAll(text, "\r\n", "\n")
		text = spurgeonBreaks.ReplaceAllString(text, "\n\n")
		var lines []string
		for _, line := range strings.Split(text, "\n") {
			line = strings.TrimSpace(spurgeonSpaces.ReplaceAllString(line, " "))
			if line != "" {
				lines = append(lines, line)
			}
		}
		slot := "pm"
		if am {
			slot = "am"
		}
		readings, ok := out.Readings[day]
		if !ok {
			readings = map[string]spurgeonReading{}
			out.Readings[day] = readings
		}
		readings[slot] = spurgeonReading{Verse: verse, Text: strings.Join(lines, "\n\n")}
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	b, err = json.MarshalIndent(out, "", " ")
	if err != nil {
		return err
	}
	path := filepath.Join(outDir, "spurgeon.json")
	if err := os.WriteFile(path, b, 0o644); err != nil {
		return err
	}
	fmt.Printf("✓ spurgeon %d readings\n", len(out.Readings))
	return nil
}
