package app

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"codie-bibel/internal/i18n"
)

// Every language the app speaks shows its own words in the sidebar: the
// reader picks one, and the pages follow.
func TestLanguageSwitches(t *testing.T) {
	a := newTestApp()
	tt := ui.NewTester(a.view, 1120, 760)
	for _, l := range i18n.Languages {
		a.setLanguage(l.Code)
		settled(tt)
		want := i18n.T("Read")
		if !hasText(tt, want) {
			t.Errorf("%s: the sidebar doesn't read %q", l.Code, want)
		}
	}
	a.setLanguage(i18n.English)
	settled(tt)
	if !hasText(tt, "Read") {
		t.Error("English didn't come back")
	}
}

// The days read as the language does: the heading of a page holds the
// weekday of its own words.
func TestTodayLocalized(t *testing.T) {
	defer i18n.SetLanguage(i18n.English)
	day := time.Date(2026, 5, 17, 12, 0, 0, 0, time.UTC) // a Sunday
	i18n.SetLanguage(i18n.ChineseSimplified)
	if got := i18n.Today(day); got != "5月17日 星期日" {
		t.Errorf("Chinese Sunday: %q", got)
	}
	i18n.SetLanguage(i18n.Japanese)
	if got := i18n.Today(day); got != "5月17日 日" {
		t.Errorf("Japanese Sunday: %q", got)
	}
	i18n.SetLanguage(i18n.Russian)
	if got := i18n.Today(day); got != "воскресенье, 17 мая" {
		t.Errorf("Russian Sunday: %q", got)
	}
	i18n.SetLanguage(i18n.Arabic)
	if got := i18n.Today(day); got != "الأحد, 17 مايو" {
		t.Errorf("Arabic Sunday: %q", got)
	}
}
