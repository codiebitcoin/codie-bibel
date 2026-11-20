package app

import (
	_ "embed"
	"log"
	"runtime"
	"strconv"
	"time"

	"github.com/egoist/mygo"

	"codie-bibel/internal/devotion"
	"codie-bibel/internal/seasons"
)

// The app's icon in the menu bar on macOS, and always in the
// notification area elsewhere: the window's way back, the verse of the
// day, the feast ahead, and the quitting.

//go:embed tray/tray-template.png
var trayTemplate []byte

//go:embed tray/tray.png
var trayColor []byte

// applyTray puts the icon in the menu bar or notification area, with its
// menu.
func (a *app) applyTray() {
	icon, template := trayColor, false
	if runtime.GOOS == "darwin" {
		icon, template = trayTemplate, true
	}
	t, err := mygo.NewTray(mygo.TrayOptions{
		ID: "codie-bibel", Icon: icon, IconIsTemplate: template,
		ToolTip: mygo.App.Name(), Menu: a.trayMenu(),
	})
	if err != nil {
		log.Println("tray:", err)
		return
	}
	// Windows' left click, where the right one opens the menu.
	t.OnClick(a.showMain)
	a.tray = t
}

// trayMenu is the icon's menu: the window, the verse of the day, the
// feast ahead, the hymns' play, and quitting.
func (a *app) trayMenu() *mygo.Menu {
	run := func(fn func()) func(*mygo.MenuItem, *mygo.Window) {
		return func(*mygo.MenuItem, *mygo.Window) {
			fn()
			a.redraw()
		}
	}
	code, ch, vs := devotion.VerseOfTheDay(time.Now())
	verse, _ := bibleRef(code, ch, vs)
	_, feast, _, days := seasons.Today(time.Now())
	feastLine := T("No feast ahead")
	if feast.Code != "" {
		feastLine = feast.Name + " — " + strconv.Itoa(days) + " " + T("days")
	}
	name := mygo.App.Name()
	return mygo.NewMenu([]*mygo.MenuItem{
		{Label: Tf("Open %s", name), Click: run(a.showMain)},
		mygo.Separator(),
		{Label: T("Verse of the Day") + ": " + verse, Disabled: true},
		{Label: feastLine, Disabled: true},
		mygo.Separator(),
		{Label: T("Play/Pause"), Click: run(func() { a.pl.TogglePlay() })},
		mygo.Separator(),
		{Label: T("Settings"), Click: run(func() { a.open(pageSettings) })},
		{Label: Tf("Quit %s", name), Click: func(*mygo.MenuItem, *mygo.Window) { mygo.App.Quit() }},
	})
}
