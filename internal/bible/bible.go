// Package bible is the Bibles the app ships with: their text, their books,
// and what a passage of them is. The text is data, generated into data/ by
// cmd/genbible and read here: one gzipped JSON per translation, opened only
// when the reader turns to it.
package bible

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"sync"
)

// indexName is the registry's index of the translations, at the FS's root.
const indexName = "index.json"

// ErrNotFound is what a translation, book, chapter or verse that isn't there
// returns.
var ErrNotFound = errors.New("not found")

// Info is what the app knows of a translation without opening it: what the
// translation list shows.
type Info struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Short    string `json:"short"`
	Language string `json:"language"`
	License  string `json:"license"`
	Source   string `json:"source"`
}

// Book is one book of a translation: its chapters, each its verses, in the
// book's own numbering.
type Book struct {
	Code     string     `json:"code"`
	Name     string     `json:"name"`
	Chapters [][]string `json:"chapters"`
}

// Verse reads one verse.
func (b *Book) Verse(chapter, verse int) (string, bool) {
	if chapter < 1 || chapter > len(b.Chapters) {
		return "", false
	}
	ch := b.Chapters[chapter-1]
	if verse < 1 || verse > len(ch) {
		return "", false
	}
	return ch[verse-1], true
}

// Translation is a whole Bible: its books, by its own names and numbering.
type Translation struct {
	Info
	Books []*Book `json:"books"`

	byCode map[string]*Book
}

// Book returns the book of that USFM code, as GEN is Genesis, or nil for
// one the translation doesn't have.
func (t *Translation) Book(code string) *Book {
	if t.byCode == nil {
		return nil
	}
	return t.byCode[code]
}

// book indexes the books by code, in the canon's order. Their counts of
// chapters were checked against the canon as they were generated, and a
// translation the canon doesn't know (luther1912's four chapters of Joel)
// read as its own: what genbible wrote, this takes.
func (t *Translation) book() error {
	t.byCode = make(map[string]*Book, len(t.Books))
	for i, b := range t.Books {
		if _, ok := Canon[b.Code]; !ok {
			return fmt.Errorf("%s: unknown book %q", t.ID, b.Code)
		}
		if b.Code != CanonCodes[i] {
			return fmt.Errorf("%s: book %d is %s, want %s", t.ID, i, b.Code, CanonCodes[i])
		}
		if len(b.Chapters) == 0 {
			return fmt.Errorf("%s: %s has no chapters", t.ID, b.Code)
		}
		t.byCode[b.Code] = b
	}
	return nil
}

// Verse reads one verse of a translation.
func (t *Translation) Verse(code string, chapter, verse int) (string, bool) {
	b := t.Book(code)
	if b == nil {
		return "", false
	}
	return b.Verse(chapter, verse)
}

// Registry is the Bibles the app has, by ID, opened only as they're read.
type Registry struct {
	fsys   fs.FS
	infos  []Info
	mu     sync.Mutex
	loaded map[string]*Translation
}

// NewRegistry reads the index of the translations in the FS. The FS is the
// app's embed of bible/data, its root the data directory itself, or a
// test's own.
func NewRegistry(fsys fs.FS) (*Registry, error) {
	r := &Registry{fsys: fsys, loaded: map[string]*Translation{}}
	f, err := fsys.Open(indexName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&r.infos); err != nil {
		return nil, err
	}
	return r, nil
}

// List is every translation the app ships with, in the index's order, which
// genbible makes: English first, then the rest by language.
func (r *Registry) List() []Info {
	return r.infos
}

// Get opens a translation, the first time from its gzipped JSON, and keeps
// it for the next.
func (r *Registry) Get(id string) (*Translation, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if t, ok := r.loaded[id]; ok {
		return t, nil
	}
	f, err := r.fsys.Open(id + ".json.gz")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	var t Translation
	if err := json.NewDecoder(gz).Decode(&t); err != nil {
		return nil, err
	}
	if err := t.book(); err != nil {
		return nil, err
	}
	for _, info := range r.infos {
		if info.ID == id {
			t.Info = info
			break
		}
	}
	r.loaded[id] = &t
	return &t, nil
}
