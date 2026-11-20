package app

import (
	"runtime"

	"github.com/egoist/mygo/ui"
)

// palette holds the few colors the app's views share. Its accent is a
// burgundy, the deep red of a bound Bible, where MyGo's own is blue.
type palette struct {
	dark      bool
	bg        ui.Color // the content's background
	sidebar   ui.Color // under the sidebar where no material shows
	text      ui.Color
	muted     ui.Color
	faint     ui.Color
	hairline  ui.Color
	hover     ui.Color
	press     ui.Color
	accent    ui.Color
	onAccent  ui.Color
	cardShade ui.Color
}

// useAccent makes the app's accent the theme's, with its hover and press: left
// alone, those stay MyGo's blue, which primary buttons turned under the pointer.
// They go darker in light mode; in dark mode hover goes lighter, as MyGo's do.
func useAccent(t *ui.Theme, p palette) {
	black, white := ui.Hex("#000000"), ui.Hex("#ffffff")
	t.Accent, t.AccentText = p.accent, p.onAccent
	t.AccentHover, t.AccentPressed = p.accent.MixOklab(black, 0.1), p.accent.MixOklab(black, 0.2)
	if p.dark {
		t.AccentHover, t.AccentPressed = p.accent.MixOklab(white, 0.15), p.accent.MixOklab(black, 0.15)
	}
}

func paletteFor(t *ui.Theme) palette {
	if t.Dark {
		return palette{
			dark:      true,
			bg:        ui.Hex("#141416"),
			sidebar:   ui.Hex("#1c1c1f"),
			text:      ui.Hex("#f5f5f7"),
			muted:     ui.RGBA(235, 235, 245, 0.56),
			faint:     ui.RGBA(235, 235, 245, 0.3),
			hairline:  ui.RGBA(255, 255, 255, 0.07),
			hover:     ui.RGBA(255, 255, 255, 0.06),
			press:     ui.RGBA(255, 255, 255, 0.1),
			accent:    ui.Oklch(0.66, 0.15, 22),
			onAccent:  ui.Hex("#ffffff"),
			cardShade: ui.RGBA(0, 0, 0, 0.45),
		}
	}
	return palette{
		bg:        ui.Hex("#ffffff"),
		sidebar:   ui.Hex("#f4f4f6"),
		text:      ui.Hex("#1d1d1f"),
		muted:     ui.RGBA(60, 60, 67, 0.62),
		faint:     ui.RGBA(60, 60, 67, 0.32),
		hairline:  ui.RGBA(0, 0, 0, 0.07),
		hover:     ui.RGBA(0, 0, 0, 0.04),
		press:     ui.RGBA(0, 0, 0, 0.08),
		accent:    ui.Oklch(0.52, 0.15, 22),
		onAccent:  ui.Hex("#ffffff"),
		cardShade: ui.RGBA(0, 0, 0, 0.16),
	}
}

// pageGap is the room between the bar over a page and its content.
const pageGap = 12

// insets returns the room the bar over a page takes at its top and the
// player's at its bottom: pages pad their content by it, and scroll under
// the bars.
func (a *app) insets(c *ui.Context) (top, bottom float32) {
	top = topBar(c)
	if runtime.GOOS == "darwin" {
		// The edge reaches a little below the traffic lights' strip.
		top += topEdgePad
	}
	return top, 0
}

// topEdgePad is the room the edge over a page leaves below the strip, on
// macOS.
const topEdgePad = 6

// buttonsHeight is the height of the window buttons off macOS, which fill
// the edge over a page.
const buttonsHeight = topBarHeight + topEdgePad

// topBarHeight is the least height of the strip along the top.
const topBarHeight = 40

// topBar returns the height of the strip along the top of the window, as
// tall as the window buttons and the room above and below them, which
// what sits in it centers on.
func topBar(c *ui.Context) float32 {
	if runtime.GOOS != "darwin" {
		// The whole edge, which the buttons fill: what sits in the strip
		// centers on their glyphs.
		return buttonsHeight
	}
	return max(c.TitleBar().Height, topBarHeight)
}

// cardBg is a card's background, filled when it's the one chosen.
func cardBg(p palette, chosen bool) ui.Color {
	if chosen {
		return p.accent.Alpha(0.12)
	}
	return p.cardShade
}

// centered shows a message, with a retry button if retry is set: what a
// page shows in place of what failed to load.
func (a *app) centered(c *ui.Context, err string, retry func()) {
	top, bottom := a.insets(c)
	ui.Column(c).Grow(1).Padding(top, 0, bottom, 0).Center().Gap(12).Children(func() {
		ui.Text(c, err).TextColor(a.pal.muted).FontSize(13).TextAlign(ui.Center).MaxWidth(520)
		if retry != nil && ui.Button(c.Key("retry"), T("Retry")).Clicked() {
			retry()
		}
	})
}

// pageHeading is a page's title, and what sits at its right.
func (a *app) pageHeading(c *ui.Context, title string, right func()) {
	top, _ := a.insets(c)
	ui.Row(c).Padding(top, 24, 0, 24).Height(top + 40).Gap(12).AlignItems(ui.Center).Children(func() {
		ui.Text(c, title).FontSize(22).FontWeight(700).TextColor(a.pal.text)
		ui.Box(c).Grow(1)
		if right != nil {
			right()
		}
	})
}
