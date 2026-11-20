package app

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// A book opens to its chapters, and a chapter to its verses: John 3
// reads its famous verse, with the words of Christ in the runs the verse
// row shows.
func TestReaderOpensChapter(t *testing.T) {
	a := newTestApp()
	tt := ui.NewTester(a.view, 1120, 760)
	a.openBook("JHN")
	settled(tt)
	if a.reader().code != "JHN" {
		t.Fatalf("the book reads %q", a.reader().code)
	}
	a.openChapter("JHN", 3)
	settled(tt)
	if len(a.reader().book.Chapters[2]) < 30 {
		t.Fatalf("John 3 loaded %d verses", len(a.reader().book.Chapters[2]))
	}
	// The verse reads as the translation has it.
	if got := a.verseText(16); len(got) < 40 {
		t.Errorf("John 3:16 reads %q", got)
	}
	// The words of Christ run red, and a verse beside them doesn't: the
	// runs of verse 16 wear the red, the runs of verse 33 wear the text's.
	runs := a.verseRuns(16, a.verseText(16))
	other := a.verseRuns(33, a.verseText(33))
	if len(runs) == 0 || len(other) == 0 {
		t.Fatal("a verse cut to no runs")
	}
	if runs[0].Color == other[0].Color {
		t.Error("Christ's words read the color of the rest")
	}
}

// The footer walks to the chapters before and after, round the canon's
// ends.
func TestReaderWalksChapters(t *testing.T) {
	a := newTestApp()
	a.openChapter("JHN", 3)
	a.nextChapter()
	if a.reader().chapter != 4 {
		t.Fatalf("next went to %d", a.reader().chapter)
	}
	a.prevChapter()
	a.prevChapter()
	if a.reader().chapter != 2 {
		t.Fatalf("prev went to %d", a.reader().chapter)
	}
	// Genesis 1's previous is Revelation's last.
	a.openChapter("GEN", 1)
	a.prevChapter()
	if a.reader().code != "REV" || a.reader().chapter != 22 {
		t.Errorf("prev from Genesis reads %s %d", a.reader().code, a.reader().chapter)
	}
	a.nextChapter()
	if a.reader().code != "GEN" || a.reader().chapter != 1 {
		t.Errorf("next from Revelation reads %s %d", a.reader().code, a.reader().chapter)
	}
}

// A search runs over the translation and a result opens its chapter.
func TestReaderSearch(t *testing.T) {
	a := newTestApp()
	a.search("In the beginning")
	if len(a.reader().results) < 2 {
		t.Fatalf("the search found %d verses", len(a.reader().results))
	}
	hit := a.reader().results[0]
	a.openSearchResult(hit)
	if a.reader().code != hit.Ref.Book || a.reader().chapter != hit.Ref.Chapter {
		t.Errorf("the result opened %s %d", a.reader().code, a.reader().chapter)
	}
}

// A highlight marks the verse on the page and stays for the next read.
func TestReaderHighlights(t *testing.T) {
	a := newTestApp()
	a.openChapter("PSA", 23)
	a.toggleHighlight(1)
	if !a.reader().highlights[1] {
		t.Fatal("the verse didn't mark")
	}
	// A fresh read of the chapter still has it.
	a.openChapter("PSA", 23)
	if !a.reader().highlights[1] {
		t.Error("the highlight didn't stay")
	}
	a.toggleHighlight(1)
	if a.reader().highlights[1] {
		t.Error("the highlight didn't go")
	}
}

// The translation switch reads the same chapter in the other
// translation.
func TestReaderSwitchesTranslation(t *testing.T) {
	a := newTestApp()
	a.openChapter("JHN", 3)
	if a.reader().tr.ID != "kjv" {
		t.Fatalf("the reader opens %s", a.reader().tr.ID)
	}
	a.setTranslation("cuv-simp")
	if a.reader().tr.ID != "cuv-simp" {
		t.Fatalf("the switch read %s", a.reader().tr.ID)
	}
	if got := a.verseText(16); got == "" || !containsRunes(got, '神') {
		t.Errorf("John 3:16 reads %q", got)
	}
}

// Listening queues the chapter to the device, from where it stood.
func TestReaderListenQueues(t *testing.T) {
	a := newTestApp()
	a.openChapter("PSA", 23)
	a.store.SetProgress(a.listenKey(), 2)
	a.listenChapter()
	items := a.bq.Items()
	if len(items) != 1 {
		t.Fatalf("the queue holds %d readings", len(items))
	}
	if items[0].At != 2 {
		t.Errorf("the reading stands at %d, want 2", items[0].At)
	}
	q := a.bq
	q.Stop()
}

func containsRunes(s string, r rune) bool {
	for _, c := range s {
		if c == r {
			return true
		}
	}
	return false
}
