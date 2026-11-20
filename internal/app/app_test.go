package app

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// Clicking a page in the sidebar shows that page.
func TestSidebarNavigates(t *testing.T) {
	a := newTestApp()
	tt := ui.NewTester(a.view, 1120, 760)

	if _, ok := tt.Find("Read"); !ok {
		t.Fatal("no Read in the sidebar")
	}
	tt.Click("Pray")
	if a.page != pagePray {
		t.Errorf("after clicking Pray, the page is %d", a.page)
	}
	settled(tt)
	tt.Click("Seasons")
	if a.page != pageSeasons {
		t.Errorf("after clicking Seasons, the page is %d", a.page)
	}
}

// The app posts to the main thread through post, which a test runs itself.
func TestPostCollects(t *testing.T) {
	a := newTestApp()
	ch := posts(a)
	done := make(chan struct{})
	a.post(func() { close(done) })
	runPost(t, ch)
	select {
	case <-done:
	default:
		t.Error("the posted function hadn't run")
	}
}
