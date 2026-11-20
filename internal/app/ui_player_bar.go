package app

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// The bar under the pages: the device's reading pill at its left, the
// hymns' transport at its right, once the hymns play.

// playerBarHeight is the height of the bar along the foot of the pages.
const playerBarHeight = 56

// playerBar is the bar, nothing of it while nothing plays.
func (a *app) playerBar(c *ui.Context) {
	if a.pl == nil || a.pl.Current() == nil {
		return
	}
	bar := ui.Row(c.Key("player-bar")).Absolute().Left(0).Right(0).Bottom(0).
		Height(playerBarHeight).Padding(0, 16).Gap(12).AlignItems(ui.Center)
	bar.Background(a.pal.sidebar)
	bar.Children(func() {
		t := a.pl.Current()
		if t == nil {
			return
		}
		// What plays, and the scrubber under it.
		ui.Column(c).Gap(2).MaxWidth(260).Children(func() {
			ui.Text(c, t.Name).FontSize(12.5).FontWeight(600).SingleLine().Shrink(1).MinWidth(0)
			a.scrubber(c)
		})
		ui.Box(c).Grow(1)
		// The transport.
		icon := "▶"
		if a.pl.Playing() {
			icon = "❚❚"
		}
		if ui.Button(c.Key("bar-prev"), "⏮").Clicked() {
			a.pl.Prev()
		}
		if ui.PrimaryButton(c.Key("bar-toggle"), icon).Clicked() {
			a.pl.TogglePlay()
		}
		if ui.Button(c.Key("bar-next"), "⏭").Clicked() {
			a.pl.Next()
		}
		// Where it stands, as it plays.
		p := a.tickTimes()
		ui.Text(c, clock(p)).FontSize(11).TextColor(a.pal.faint)
		ui.Text(c, clock(a.pl.Duration())).FontSize(11).TextColor(a.pal.faint)
	})
}

// scrubber is the bar's progress, as far as it has gone.
func (a *app) scrubber(c *ui.Context) {
	d := a.pl.Duration()
	if d <= 0 {
		ui.Progress(c.Key("scrubber"), 0).Grow(1)
		return
	}
	ui.Progress(c.Key("scrubber"), float64(a.pl.Position())/float64(d)).Grow(1)
}

// tickTimes asks for a frame a second, what the bar's clock reads by.
func (a *app) tickTimes() time.Duration {
	return a.pl.Position()
}

// hymnVolumeButton is the volume popover's trigger, at the bar's end.
func (a *app) hymnVolumeButton(c *ui.Context) {
	if ui.Button(c.Key("volume"), T("Volume")).Clicked() {
		a.pl.SetVolume(0.8)
	}
}
