package bible

import (
	"os"
	"strings"
	"testing"
)

// The Bibles the app ships with all open, and the classic verses read as
// they should: this is the data genbible generated, read back whole.
func TestShippedDataOpens(t *testing.T) {
	if _, err := os.Stat("data"); err != nil {
		t.Skip("no data directory: genbible hasn't run here")
	}
	r, err := NewRegistry(os.DirFS("data"))
	if err != nil {
		t.Fatal(err)
	}
	infos := r.List()
	if len(infos) < 13 {
		t.Fatalf("the app ships %d translations", len(infos))
	}
	for _, info := range infos {
		if _, err := r.Get(info.ID); err != nil {
			t.Errorf("%s: %v", info.ID, err)
		}
	}
	kjv, err := r.Get("kjv")
	if err != nil {
		t.Fatal(err)
	}
	got, ok := kjv.Verse("JHN", 3, 16)
	if !ok || !strings.HasPrefix(got, "For God so loved the world") {
		t.Errorf("KJV John 3:16: %q", got)
	}
	cuv, err := r.Get("cuv-simp")
	if err != nil {
		t.Fatal(err)
	}
	got, ok = cuv.Verse("JHN", 3, 16)
	if !ok || !strings.HasPrefix(got, "神爱世人") {
		t.Errorf("CUV John 3:16: %q", got)
	}
}
