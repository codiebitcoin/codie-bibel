package app

import (
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"

	"codie-bibel/internal/hymn"
	"codie-bibel/internal/player"
	"codie-bibel/internal/player/media"
)

// The Music page: the hymns the app ships, playing quietly, with the
// bar's transport under the pages.

func (a *app) musicPage(c *ui.Context) {
	a.pageHeading(c, T("Music"), func() {
		if a.pl.Shuffle() {
			if ui.PrimaryButton(c.Key("shuffle"), T("Shuffle")).Clicked() {
				a.pl.ToggleShuffle()
			}
		} else if ui.Button(c.Key("shuffle"), T("Shuffle")).Clicked() {
			a.pl.ToggleShuffle()
		}
		if a.pl.Repeat() == player.RepeatOne {
			if ui.PrimaryButton(c.Key("repeat"), T("Repeat")).Clicked() {
				a.pl.CycleRepeat()
			}
		} else if ui.Button(c.Key("repeat"), T("Repeat")).Clicked() {
			a.pl.CycleRepeat()
		}
	})
	top, bottom := a.insets(c)
	ui.Scroll(c.Key("hymns")).Grow(1).Children(func() {
		ui.Column(c).Padding(top, 24, bottom, 24).Gap(2).Children(func() {
			for i, t := range hymn.Tracks {
				a.hymnRow(c, i, t)
			}
			ui.Text(c, T("Every hymn here is in the public domain, or free as the app is; the licenses stand in LICENSES.md.")).TextColor(a.pal.faint).FontSize(11.5)
		})
	})
}

// hymnRow is one hymn: its title, its performer, the time it takes, and
// what plays when it's tapped.
func (a *app) hymnRow(c *ui.Context, i int, t hymn.Track) {
	current := a.pl.Current()
	playing := current != nil && current.ID == t.ID
	row := ui.Row(c.Key("hymn-"+t.ID)).Padding(10, 12).Gap(10).Radius(10).AlignItems(ui.Center)
	bg := ui.Transparent
	switch {
	case playing && a.pl.Playing():
		bg = a.pal.accent.Alpha(0.12)
	case row.Hovered():
		bg = a.pal.hover
	}
	row.Background(bg)
	row.Children(func() {
		mark := strconv.Itoa(i + 1)
		if playing && a.pl.Playing() {
			mark = "▶"
		}
		ui.Text(c, mark).FontSize(12).TextColor(a.pal.accent)
		ui.Column(c).Gap(1).Children(func() {
			ui.Text(c, t.Title).FontSize(13.5).FontWeight(500)
			ui.Text(c, t.Artist).FontSize(11.5).TextColor(a.pal.faint)
		})
		ui.Box(c).Grow(1)
		ui.Text(c, clock(t.Duration)).FontSize(11.5).TextColor(a.pal.faint)
	})
	row.Label(t.Title)
	if row.Clicked() {
		a.playHymnAt(i)
	}
}

// playHymnAt queues the hymns from the one at i, playing it.
func (a *app) playHymnAt(i int) {
	var tracks []media.Track
	for _, t := range hymn.Tracks {
		tracks = append(tracks, media.Track{
			ID: t.ID, Source: "hymn", Name: t.Title,
			Artists:  []media.Artist{{Name: t.Artist}},
			Duration: int64(t.Duration / time.Millisecond),
		})
	}
	a.pl.PlayQueue(tracks, i)
}

// clock reads a duration as its minutes and seconds.
func clock(d time.Duration) string {
	s := int(d.Seconds())
	return strconv.Itoa(s/60) + ":" + leadingZero(s%60)
}

func leadingZero(n int) string {
	if n < 10 {
		return "0" + strconv.Itoa(n)
	}
	return strconv.Itoa(n)
}
