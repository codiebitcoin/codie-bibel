// Package store is what the reader keeps on this machine: highlights,
// notes, bookmarks, the prayers, and where each reading stands. It is
// SQLite in the reader's data directory; nothing of it leaves the
// machine.
package store

import (
	"database/sql"
	"os"
	"sync"

	_ "modernc.org/sqlite"
)

// Store is the reader's own data.
type Store struct {
	db  *sql.DB
	mu  sync.Mutex
	Dir string
}

// New opens the store in dir, making it and its tables as they're
// needed.
func New(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", dir+"/bibel.db?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)")
	if err != nil {
		return nil, err
	}
	for _, q := range []string{
		`CREATE TABLE IF NOT EXISTS kv (k TEXT PRIMARY KEY, v TEXT NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS highlights (
			tr TEXT NOT NULL, code TEXT NOT NULL, chapter INTEGER NOT NULL, verse INTEGER NOT NULL,
			PRIMARY KEY (tr, code, chapter, verse))`,
		`CREATE TABLE IF NOT EXISTS notes (
			tr TEXT NOT NULL, code TEXT NOT NULL, chapter INTEGER NOT NULL, verse INTEGER NOT NULL,
			text TEXT NOT NULL, updated TEXT NOT NULL,
			PRIMARY KEY (tr, code, chapter, verse))`,
		`CREATE TABLE IF NOT EXISTS prayers (
			id TEXT PRIMARY KEY, text TEXT NOT NULL, added TEXT NOT NULL, answered INTEGER NOT NULL DEFAULT 0)`,
		`CREATE TABLE IF NOT EXISTS progress (k TEXT PRIMARY KEY, at INTEGER NOT NULL)`,
	} {
		if _, err := db.Exec(q); err != nil {
			db.Close()
			return nil, err
		}
	}
	return &Store{db: db, Dir: dir}, nil
}

// Close closes the store.
func (s *Store) Close() error { return s.db.Close() }

// Put keeps a setting or a small state by its key.
func (s *Store) Put(k, v string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.db.Exec(`INSERT INTO kv (k, v) VALUES (?, ?) ON CONFLICT(k) DO UPDATE SET v = excluded.v`, k, v)
	return err
}

// Get reads what Put kept, or "" with false for a key that isn't there.
func (s *Store) Get(k string) (string, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var v string
	err := s.db.QueryRow(`SELECT v FROM kv WHERE k = ?`, k).Scan(&v)
	if err == sql.ErrNoRows {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return v, true, nil
}

// Highlight marks a verse highlighted, or not.
func (s *Store) Highlight(tr, code string, chapter, verse int, on bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if on {
		_, err := s.db.Exec(`INSERT OR IGNORE INTO highlights (tr, code, chapter, verse) VALUES (?, ?, ?, ?)`, tr, code, chapter, verse)
		return err
	}
	_, err := s.db.Exec(`DELETE FROM highlights WHERE tr = ? AND code = ? AND chapter = ? AND verse = ?`, tr, code, chapter, verse)
	return err
}

// Highlights is the set of highlighted verses of a chapter.
func (s *Store) Highlights(tr, code string, chapter int) (map[int]bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	rows, err := s.db.Query(`SELECT verse FROM highlights WHERE tr = ? AND code = ? AND chapter = ?`, tr, code, chapter)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[int]bool{}
	for rows.Next() {
		var v int
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		out[v] = true
	}
	return out, rows.Err()
}

// SetNote writes a note for a verse, or drops it when the text is empty.
func (s *Store) SetNote(tr, code string, chapter, verse int, text string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if text == "" {
		_, err := s.db.Exec(`DELETE FROM notes WHERE tr = ? AND code = ? AND chapter = ? AND verse = ?`, tr, code, chapter, verse)
		return err
	}
	_, err := s.db.Exec(`INSERT INTO notes (tr, code, chapter, verse, text, updated) VALUES (?, ?, ?, ?, ?, datetime('now'))
		ON CONFLICT(tr, code, chapter, verse) DO UPDATE SET text = excluded.text, updated = excluded.updated`,
		tr, code, chapter, verse, text)
	return err
}

// Note reads a verse's note, or "" with false for one that isn't there.
func (s *Store) Note(tr, code string, chapter, verse int) (string, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var text string
	err := s.db.QueryRow(`SELECT text FROM notes WHERE tr = ? AND code = ? AND chapter = ? AND verse = ?`, tr, code, chapter, verse).Scan(&text)
	if err == sql.ErrNoRows {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return text, true, nil
}

// Prayer is one prayer of the reader's list.
type Prayer struct {
	ID       string
	Text     string
	Answered bool
}

// AddPrayer adds a prayer to the list.
func (s *Store) AddPrayer(id, text string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.db.Exec(`INSERT INTO prayers (id, text, added) VALUES (?, ?, datetime('now'))`, id, text)
	return err
}

// SetPrayerAnswered marks a prayer answered, or not.
func (s *Store) SetPrayerAnswered(id string, answered bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	on := 0
	if answered {
		on = 1
	}
	_, err := s.db.Exec(`UPDATE prayers SET answered = ? WHERE id = ?`, on, id)
	return err
}

// DropPrayer takes a prayer off the list.
func (s *Store) DropPrayer(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.db.Exec(`DELETE FROM prayers WHERE id = ?`, id)
	return err
}

// Prayers is the list, the unanswered first and the newest of them
// first.
func (s *Store) Prayers() ([]Prayer, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	rows, err := s.db.Query(`SELECT id, text, answered FROM prayers ORDER BY answered, rowid DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Prayer
	for rows.Next() {
		var p Prayer
		var answered int
		if err := rows.Scan(&p.ID, &p.Text, &answered); err != nil {
			return nil, err
		}
		p.Answered = answered != 0
		out = append(out, p)
	}
	return out, rows.Err()
}

// SetProgress keeps how far a reading stands, as the turn it reads next.
func (s *Store) SetProgress(k string, at int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.db.Exec(`INSERT INTO progress (k, at) VALUES (?, ?) ON CONFLICT(k) DO UPDATE SET at = excluded.at`, k, at)
	return err
}

// Progress reads how far a reading stands, or 0.
func (s *Store) Progress(k string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	var at int
	if err := s.db.QueryRow(`SELECT at FROM progress WHERE k = ?`, k).Scan(&at); err != nil {
		return 0
	}
	return at
}
