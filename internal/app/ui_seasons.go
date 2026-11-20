package app

import (
	"context"
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"

	"codie-bibel/internal/bible"
	"codie-bibel/internal/seasons"
)

// The Seasons page: where the church year stands, its feasts, and the
// readings the Sunday keeps.

func (a *app) seasonsPage(c *ui.Context) {
	a.pageHeading(c, T("Seasons"), nil)
	top, bottom := a.insets(c)
	ui.Scroll(c.Key("seasons")).Grow(1).Children(func() {
		ui.Column(c).Padding(top, 24, bottom, 24).Gap(16).Children(func() {
			a.seasonCard(c)
			a.weekReadingsCard(c)
			a.feastsCard(c)
		})
	})
}

// seasonCard is where the year stands: the season's name and color, the
// feast of the day, and the days to the next feast.
func (a *app) seasonCard(c *ui.Context) {
	season, feast, next, days := seasons.Today(time.Now())
	ui.Column(c).Padding(16, 18).Gap(8).Radius(12).Background(a.pal.cardShade).Children(func() {
		ui.Row(c).Gap(10).AlignItems(ui.Center).Children(func() {
			ui.Box(c.Key("season-color")).Size(14, 14).Radius(7).Background(seasonColor(string(season.Color)))
			ui.Text(c, season.Name).FontSize(18).FontWeight(700)
			ui.Box(c).Grow(1)
			if days >= 0 {
				ui.Text(c, Tf("%d days to %s", days, next.Name)).FontSize(12.5).TextColor(a.pal.muted)
			}
		})
		if feast != nil {
			ui.Text(c, feast.Name).FontSize(14).FontWeight(600).TextColor(a.pal.accent)
			if ui.Button(c.Key("feast-listen"), T("Listen")).Clicked() {
				a.listenFeast(feast.Code, feast.Name)
			}
		}
	})
}

// seasonColor is a color for the season, as the swatch reads.
func seasonColor(color string) ui.Color {
	switch color {
	case "white":
		return ui.Hex("#f5f1e8")
	case "purple":
		return ui.Hex("#6a4c93")
	case "red":
		return ui.Hex("#c0392b")
	}
	return ui.Hex("#4a7c59")
}

// weekReadingsCard is the Sunday's readings, as the lectionary keeps
// them, each opening in Read and all of them queuing to the device.
func (a *app) weekReadingsCard(c *ui.Context) {
	l, err := seasons.NewLectionary()
	if err != nil {
		return
	}
	yearLetter, _, lessons, ok := l.For(time.Now())
	if !ok {
		return
	}
	ui.Column(c).Gap(8).Children(func() {
		ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
			ui.Text(c, T("This week's readings")).FontSize(12).FontWeight(600).TextColor(a.pal.muted)
			ui.Text(c, "Year "+yearLetter).FontSize(11.5).TextColor(a.pal.faint)
		})
		for _, lesson := range lessons {
			if lesson == "" {
				continue
			}
			ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
				ui.Text(c, lesson).FontSize(13)
				ui.Box(c).Grow(1)
				if ref, ok := bible.ParseRef(lesson); ok {
					if ui.Button(c.Key("lesson-open-"+lesson), T("Read")).Clicked() {
						a.openChapter(ref.Book, ref.Chapter)
					}
				}
			})
		}
	})
}

// feastsCard is the year's feasts in their order, with the days to each.
func (a *app) feastsCard(c *ui.Context) {
	year := time.Now().Year()
	feasts := seasons.Feasts(year)
	ui.Column(c).Gap(6).Children(func() {
		ui.Text(c, seasonYearName(year)).FontSize(12).FontWeight(600).TextColor(a.pal.muted)
		now := time.Now()
		for _, f := range feasts {
			days := daysTo(now, f.Date)
			ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
				ui.Text(c, f.Name).FontSize(13)
				ui.Box(c).Grow(1)
				mark := f.Date.Format("Jan 2")
				if days >= 0 {
					mark = Tf("in %d days", days) + " · " + mark
				}
				ui.Text(c, mark).FontSize(12).TextColor(a.pal.faint)
				if days >= 0 && days <= 7 {
					if ui.Button(c.Key("feast-soon-listen-"+f.Code), T("Listen")).Clicked() {
						a.listenFeast(f.Code, f.Name)
					}
					if ui.Button(c.Key("feast-soon-alarm-"+f.Code), T("Remind me")).Clicked() {
						a.setFeastAlarm(f)
					}
				}
			})
		}
	})
}

// seasonYearName names the liturgical year, as its Advent names it.
func seasonYearName(year int) string {
	return strconv.Itoa(year) + "–" + strconv.Itoa(year+1)
}

// daysTo is how many days ahead the feast lies, 0 for today, -1 past.
func daysTo(from, to time.Time) int {
	d := int(time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, to.Location()).
		Sub(time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, from.Location())).Hours() / 24)
	return d
}

// listenFeast reads a feast's own scriptures aloud, the verses the
// seasons package keeps by the feast's code.
func (a *app) listenFeast(code, name string) {
	for _, r := range seasons.FeastVerses(code) {
		a.listenRef(r, name + ": " + r.String())
	}
}

// setFeastAlarm asks the device to speak at the feast's hour.
func (a *app) setFeastAlarm(f seasons.Feast) {
	if a.dev == nil || !a.dev.Paired() {
		a.toast = T("Pair a device first")
		a.redraw()
		return
	}
	go func() {
		at := time.Date(f.Date.Year(), f.Date.Month(), f.Date.Day(), 8, 0, 0, 0, time.Local)
		err := a.dev.SetAlarm(context.Background(), f.Name, at)
		a.post(func() {
			switch {
			case err == nil:
				a.toast = Tf("The device will remember %s.", f.Name)
			default:
				a.toast = err.Error()
			}
			a.redraw()
		})
	}()
}
