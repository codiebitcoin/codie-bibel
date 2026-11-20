package bible

import (
	"compress/gzip"
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"testing/fstest"
)

// testTranslation is a translation of the whole canon: Genesis and John
// read their own few verses, and every other book has one verse a chapter,
// so the registry's check of the canon holds.
func testTranslation() *Translation {
	texts := map[string][][]string{
		"GEN": {
			{"In the beginning God created the heaven and the earth.", "And the earth was without form, and void; and darkness was upon the face of the deep."},
			{"These are the generations of the heavens and of the earth when they were created."},
		},
		"JHN": {
			{"In the beginning was the Word, and the Word was with God, and the Word was God."},
			{"And the third day there was a marriage in Cana of Galilee."},
		},
	}
	t := &Translation{
		Info: Info{ID: "test", Name: "Test Version", Short: "TV", Language: "en", License: "Public domain"},
	}
	for _, code := range CanonCodes {
		var chapters [][]string
		custom, ok := texts[code]
		if ok {
			// The books with their own text keep it, and pad out to the
			// canon's count of chapters.
			chapters = custom
			for len(chapters) < Canon[code].Chapters {
				chapters = append(chapters, []string{"And it came to pass."})
			}
		} else {
			chapters = make([][]string, Canon[code].Chapters)
			for i := range chapters {
				chapters[i] = []string{"And it came to pass."}
			}
		}
		t.Books = append(t.Books, &Book{Code: code, Name: code, Chapters: chapters})
	}
	return t
}

// gzOf writes a translation as its data file reads.
func gzOf(t *testing.T, v Translation) []byte {
	t.Helper()
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if err := json.NewEncoder(gz).Encode(v); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

// testRegistry is a registry of one small translation, two books long, for
// the tests to read.
func testRegistry(t *testing.T) *Registry {
	t.Helper()
	fsys := fstest.MapFS{
		"index.json":   &fstest.MapFile{Data: []byte(`[{"id":"test","name":"Test Version","short":"TV","language":"en","license":"Public domain"}]`)},
		"test.json.gz": &fstest.MapFile{Data: gzOf(t, *testTranslation())},
	}
	r, err := NewRegistry(fsys)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// The registry lists what the index has and opens a translation whole: its
// books by code, its chapters by count.
func TestRegistryOpens(t *testing.T) {
	r := testRegistry(t)
	if got := len(r.List()); got != 1 {
		t.Fatalf("the list has %d translations", got)
	}
	tt, err := r.Get("test")
	if err != nil {
		t.Fatal(err)
	}
	if tt.Name != "Test Version" || tt.Short != "TV" {
		t.Errorf("the translation reads %q, %q", tt.Name, tt.Short)
	}
	if b := tt.Book("GEN"); b == nil || len(b.Chapters) != 50 || len(b.Chapters[0]) != 2 {
		t.Fatalf("Genesis reads %d chapters, %d verses in the first", len(b.Chapters), len(b.Chapters[0]))
	}
	if _, err := r.Get("nope"); err == nil {
		t.Error("a translation that isn't there opened")
	}
}

// A verse reads by book, chapter and verse; a verse that isn't there
// doesn't.
func TestVerse(t *testing.T) {
	r := testRegistry(t)
	tt, _ := r.Get("test")
	got, ok := tt.Verse("GEN", 1, 1)
	if !ok || !strings.HasPrefix(got, "In the beginning God created") {
		t.Errorf("Genesis 1:1: %q, %v", got, ok)
	}
	if got, ok := tt.Verse("GEN", 3, 1); !ok || got != "And it came to pass." {
		t.Errorf("Genesis 3:1: %q, %v", got, ok)
	}
	if got, ok := tt.Verse("GEN", 51, 1); ok {
		t.Errorf("Genesis 51 opened: %q", got)
	}
	if got, ok := tt.Verse("GEN", 1, 3); ok {
		t.Errorf("Genesis 1:3 opened: %q", got)
	}
}

// References read as people write them, in English and in Chinese, and a
// reference that isn't a place doesn't.
func TestParseRef(t *testing.T) {
	for _, c := range []struct {
		in   string
		want Ref
	}{
		{"John 3:16", Ref{"JHN", 3, 16, 0, 16}},
		{"john 3:16", Ref{"JHN", 3, 16, 0, 16}},
		{"Jn 3:16", Ref{"JHN", 3, 16, 0, 16}},
		{"1 Corinthians 13:4", Ref{"1CO", 13, 4, 0, 4}},
		{"1co 13:4-7", Ref{"1CO", 13, 4, 13, 7}},
		{"Psalm 23", Ref{"PSA", 23, 0, 0, 0}},
		{"Ps 23:1-6", Ref{"PSA", 23, 1, 23, 6}},
		{"Matthew 5-7", Ref{"MAT", 5, 0, 7, 0}},
		{"约3:16", Ref{"JHN", 3, 16, 0, 16}},
		{"约 3:16", Ref{"JHN", 3, 16, 0, 16}},
		{"约翰福音 3:16", Ref{"JHN", 3, 16, 0, 16}},
		{"诗篇 23", Ref{"PSA", 23, 0, 0, 0}},
		{"创世记 1", Ref{"GEN", 1, 0, 0, 0}},
		{"林前 13:4", Ref{"1CO", 13, 4, 0, 4}},
	} {
		got, ok := ParseRef(c.in)
		if !ok || got != c.want {
			t.Errorf("ParseRef(%q) = %+v, %v; want %+v", c.in, got, ok, c.want)
		}
	}
	for _, bad := range []string{"", "John", "John 0:3", "John 3:0", "John 3:16-14", "3:16", "Nospace 3:16"} {
		if got, ok := ParseRef(bad); ok {
			t.Errorf("ParseRef(%q) = %+v; want it refused", bad, got)
		}
	}
}

// A reference writes as people read it.
func TestRefString(t *testing.T) {
	for _, c := range []struct {
		ref  Ref
		want string
	}{
		{Ref{"JHN", 3, 16, 0, 16}, "John 3:16"},
		{Ref{"PSA", 23, 0, 0, 0}, "Psalms 23"},
		{Ref{"MAT", 5, 0, 7, 0}, "Matthew 5-7"},
		{Ref{"PSA", 23, 1, 23, 6}, "Psalms 23:1-6"},
		{Ref{"1CO", 13, 4, 13, 7}, "1 Corinthians 13:4-7"},
	} {
		if got := c.ref.String(); got != c.want {
			t.Errorf("%+v writes %q, want %q", c.ref, got, c.want)
		}
	}
}

// A search finds the verses that have the words, wherever the case, and
// stops at its limit.
func TestSearch(t *testing.T) {
	r := testRegistry(t)
	tt, _ := r.Get("test")
	hits := Search(tt, "IN THE BEGINNING", 10)
	if len(hits) != 2 {
		t.Fatalf("the search found %d verses, want 2", len(hits))
	}
	if hits[0].Ref.Book != "GEN" || hits[1].Ref.Book != "JHN" {
		t.Errorf("the hits read %+v", hits)
	}
	if hits[0].Start != 0 {
		t.Errorf("the match starts at %d", hits[0].Start)
	}
	if hits2 := Search(tt, "marriage in Cana", 10); len(hits2) != 1 || hits2[0].Ref.Chapter != 2 {
		t.Errorf("the second search found %+v", hits2)
	}
	if hits3 := Search(tt, "nowhere in this bible", 10); len(hits3) != 0 {
		t.Errorf("the third search found %+v", hits3)
	}
	if hits4 := Search(tt, "beginning", 1); len(hits4) != 1 {
		t.Errorf("a limited search found %d", len(hits4))
	}
}

// A translation that has a book the canon doesn't, or the wrong count of
// chapters, refuses to open.
func TestRegistryChecksCanon(t *testing.T) {
	bad := Translation{Books: []*Book{{Code: "XXX", Chapters: [][]string{{"x"}}}}}
	fsys := fstest.MapFS{
		"index.json":  &fstest.MapFile{Data: []byte(`[{"id":"bad","name":"Bad","short":"B","language":"en"}]`)},
		"bad.json.gz": &fstest.MapFile{Data: gzOf(t, bad)},
	}
	r, err := NewRegistry(fsys)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := r.Get("bad"); err == nil {
		t.Error("a translation of an unknown book opened")
	}
}
