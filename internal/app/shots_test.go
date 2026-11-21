package app

import (
	"image/png"
	"os"
	"testing"

	"github.com/egoist/mygo/ui"
)

// TestShots renders every page to /tmp/shots for a look at the design.
func TestShots(t *testing.T) {
	dir := "/tmp/shots"
	os.MkdirAll(dir, 0o755)
	a := newTestApp()
	a.page = pageRead
	tt := ui.NewTester(a.view, 1120, 760)
	save := func(name string) {
		settled(tt)
		img := tt.Image()
		f, _ := os.Create(dir + "/" + name + ".png")
		if f != nil {
			png.Encode(f, img)
			f.Close()
		}
	}
	save("read-books")
	a.openBook("GEN")
	save("read-chapters")
	a.openChapter("JHN", 3)
	save("read-chapter")
	a.open(pageDevotions)
	save("devotions")
	a.open(pagePray)
	save("pray")
	a.open(pageStories)
	save("stories")
	a.open(pageSeasons)
	save("seasons")
	a.open(pageMusic)
	save("music")
	a.open(pageDevice)
	save("device")
	a.open(pageSettings)
	save("settings")
	// the pairing guide
	a.open(pageDevice)
	a.wizard.open = true
	a.wizard.step = 0
	save("wizard")
}
