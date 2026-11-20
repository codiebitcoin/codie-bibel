package hymn

import (
	"os"
	"path/filepath"
	"testing"
)

// The hymns' metadata reads whole: every track names a file the package
// embeds, of the length it says.
func TestTracksRead(t *testing.T) {
	if len(Tracks) < 8 {
		t.Fatalf("only %d hymns", len(Tracks))
	}
	for _, tr := range Tracks {
		if tr.ID == "" || tr.Title == "" || tr.File == "" || tr.Duration <= 0 || tr.License == "" {
			t.Errorf("hymn %+v is thin", tr)
		}
		if _, err := files.Open("data/" + tr.File); err != nil {
			t.Errorf("%s: %v", tr.ID, err)
		}
		if ByID(tr.ID) == nil {
			t.Errorf("ByID finds no %s", tr.ID)
		}
	}
	if ByID("nothing") != nil {
		t.Error("ByID found a hymn that isn't there")
	}
}

// Ensure writes the hymns' audio to disk once, and says where.
func TestEnsureWrites(t *testing.T) {
	dir := t.TempDir()
	paths, err := Ensure(filepath.Join(dir, "hymns"))
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) != len(Tracks) {
		t.Fatalf("Ensure wrote %d files", len(paths))
	}
	for _, tr := range Tracks {
		fi, err := os.Stat(paths[tr.ID])
		if err != nil {
			t.Fatalf("%s: %v", tr.ID, err)
		}
		if fi.Size() == 0 {
			t.Errorf("%s wrote no bytes", tr.ID)
		}
	}
	// A second run writes nothing anew: the files stay as they were.
	paths2, err := Ensure(filepath.Join(dir, "hymns"))
	if err != nil {
		t.Fatal(err)
	}
	for id, p := range paths {
		if paths2[id] != p {
			t.Errorf("%s moved from %s to %s", id, p, paths2[id])
		}
	}
}
