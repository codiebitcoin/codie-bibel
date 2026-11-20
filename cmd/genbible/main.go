// Command genbible generates the Bibles the app ships with: it fetches each
// translation from its upstream, parses it into the app's shape, checks it
// against the canon, and writes internal/bible/data — one gzipped JSON per
// translation, an index of them all, and the licenses. The generated files
// are committed: a build never needs the network.
//
// Usage: go run ./cmd/genbible [-cache dir] [-out internal/bible/data] [-only id]
package main

import (
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"codie-bibel/internal/bible"
)

// source is one translation to fetch and parse, and what the index and the
// licenses say of it.
type source struct {
	ID       string // the file's stem, as kjv
	Name     string // as people read it
	Short    string
	Language string // BCP-47, as en, zh-Hans
	License  string
	Source   string // where it came from, for the licenses
	URL      string
	Kind     string // usfm, usfx, zefania, osis
	Order    int    // where the translation list shows it
}

// sources are the translations the app ships with, all of them public
// domain: NIV, ESV and the like never enter this table. English comes
// first, the rest follow by language.
var sources = []source{
	{ID: "kjv", Name: "King James Version", Short: "KJV", Language: "en", License: "Public domain", Source: "ebible.org (eng-kjv)", URL: "https://ebible.org/Scriptures/eng-kjv_usfm.zip", Kind: "usfm", Order: 1},
	{ID: "asv", Name: "American Standard Version", Short: "ASV", Language: "en", License: "Public domain", Source: "ebible.org (eng-asv)", URL: "https://ebible.org/Scriptures/eng-asv_usfm.zip", Kind: "usfm", Order: 2},
	{ID: "web", Name: "World English Bible", Short: "WEB", Language: "en", License: "Public domain", Source: "ebible.org (eng-web)", URL: "https://ebible.org/Scriptures/eng-web_usfm.zip", Kind: "usfm", Order: 3},
	{ID: "cuv-simp", Name: "和合本 (简体)", Short: "和合本", Language: "zh-Hans", License: "Public domain", Source: "seven1m/open-bibles (chi-cuv-simp)", URL: "https://raw.githubusercontent.com/seven1m/open-bibles/master/chi-cuv-simp.usfx.xml", Kind: "usfx", Order: 4},
	{ID: "cuv", Name: "和合本 (繁體)", Short: "和合本", Language: "zh-Hant", License: "Public domain", Source: "seven1m/open-bibles (chi-cuv)", URL: "https://raw.githubusercontent.com/seven1m/open-bibles/master/chi-cuv.usfx.xml", Kind: "usfx", Order: 5},
	{ID: "luther1912", Name: "Lutherbibel 1912", Short: "Luth1912", Language: "de", License: "Public domain", Source: "seven1m/open-bibles (deu-luther1912)", URL: "https://raw.githubusercontent.com/seven1m/open-bibles/master/deu-luther1912.osis.xml", Kind: "osis", Order: 6},
	{ID: "rv1909", Name: "Reina-Valera 1909", Short: "RV1909", Language: "es", License: "Public domain", Source: "ebible.org (spaRV1909)", URL: "https://ebible.org/Scriptures/spaRV1909_usfm.zip", Kind: "usfm", Order: 7},
	{ID: "ls1910", Name: "Louis Segond 1910", Short: "LS1910", Language: "fr", License: "Public domain", Source: "ebible.org (fraLSG)", URL: "https://ebible.org/Scriptures/fraLSG_usfm.zip", Kind: "usfm", Order: 8},
	{ID: "almeida", Name: "Almeida (edição de 1900/1911)", Short: "Almeida", Language: "pt", License: "Public domain", Source: "getbible.net (almeida)", URL: "https://api.getbible.net/v2/almeida.json", Kind: "getbible", Order: 9},
	{ID: "synodal", Name: "Синодальный перевод", Short: "Синод.", Language: "ru", License: "Public domain", Source: "seven1m/open-bibles (rus-synodal)", URL: "https://raw.githubusercontent.com/seven1m/open-bibles/master/rus-synodal.zefania.xml", Kind: "zefania", Order: 10},
	{ID: "kougo", Name: "口語訳聖書 (1955)", Short: "口語訳", Language: "ja", License: "Public domain", Source: "seven1m/open-bibles (jpn-kougo)", URL: "https://raw.githubusercontent.com/seven1m/open-bibles/master/jpn-kougo.osis.xml", Kind: "osis", Order: 11},
	{ID: "krv", Name: "개역한글 (1961)", Short: "개역한글", Language: "ko", License: "Public domain", Source: "getbible.net (korean)", URL: "https://api.getbible.net/v2/korean.json", Kind: "getbible", Order: 12},
	{ID: "svd", Name: "Smith & Van Dyke (1865)", Short: "SVD", Language: "ar", License: "Public domain", Source: "m0ty/bible-io-json (arb-svd-1865)", URL: "https://raw.githubusercontent.com/m0ty/bible-io-json/main/Arabic/arb-svd-1865.json", Kind: "svdjson", Order: 13},
}

// chapterExtras are the books a translation divides by its own count of
// chapters, where the canon's is the modern one. Luther 1912 keeps the
// Hebrew's four chapters of Joel, and reads Malachi in three, where the
// modern editions have three of Joel and four of Malachi.
var chapterExtras = map[string]map[string]int{
	"luther1912": {"JOL": 4, "MAL": 3},
}

// osisIDs maps the OSIS book names to the app's USFM codes: the standard
// names the OSIS spec and its converters write.
var osisIDs = map[string]string{
	"Gen": "GEN", "Exod": "EXO", "Lev": "LEV", "Num": "NUM", "Deut": "DEU",
	"Josh": "JOS", "Judg": "JDG", "Ruth": "RUT", "1Sam": "1SA", "2Sam": "2SA",
	"1Kgs": "1KI", "2Kgs": "2KI", "1Chr": "1CH", "2Chr": "2CH", "Ezra": "EZR",
	"Neh": "NEH", "Esth": "EST", "Job": "JOB", "Ps": "PSA", "Prov": "PRO",
	"Eccl": "ECC", "Song": "SNG", "Isa": "ISA", "Jer": "JER", "Lam": "LAM",
	"Ezek": "EZK", "Dan": "DAN", "Hos": "HOS", "Joel": "JOL", "Amos": "AMO",
	"Obad": "OBA", "Jonah": "JON", "Mic": "MIC", "Nah": "NAM", "Hab": "HAB",
	"Zeph": "ZEP", "Hag": "HAG", "Zech": "ZEC", "Mal": "MAL",
	"Matt": "MAT", "Mark": "MRK", "Luke": "LUK", "John": "JHN", "Acts": "ACT",
	"Rom": "ROM", "1Cor": "1CO", "2Cor": "2CO", "Gal": "GAL", "Eph": "EPH",
	"Phil": "PHP", "Col": "COL", "1Thess": "1TH", "2Thess": "2TH",
	"1Tim": "1TI", "2Tim": "2TI", "Titus": "TIT", "Phlm": "PHM",
	"Heb": "HEB", "Jas": "JAS", "1Pet": "1PE", "2Pet": "2PE",
	"1John": "1JN", "2John": "2JN", "3John": "3JN", "Jude": "JUD", "Rev": "REV",
}

func main() {
	cache := flag.String("cache", ".genbible-cache", "where the upstream files are kept")
	out := flag.String("out", "../../internal/bible/data", "where the generated data goes")
	only := flag.String("only", "", "one translation's ID, to redo only it")
	flag.Parse()

	if err := os.MkdirAll(*out, 0o755); err != nil {
		fatal(err)
	}
	var index []bible.Info
	var licenses strings.Builder
	licenses.WriteString("# Bibles the app ships with\n\nAll are public domain. Each is what its upstream ships; the app edits\nnothing of the text.\n\n| ID | Name | Language | License | Source |\n|---|---|---|---|---|\n")
	for _, s := range sources {
		if s.URL == "" {
			fmt.Printf("– %-9s no upstream yet, skipped\n", s.ID)
			continue
		}
		if *only != "" && *only != s.ID {
			continue
		}
		local := filepath.Join(*cache, s.ID+filepath.Ext(s.URL))
		if err := fetch(s.URL, local); err != nil {
			fatal(err)
		}
		chapters, err := parse(s.Kind, local)
		if err != nil {
			fatal(fmt.Errorf("%s: %w", s.ID, err))
		}
		t, err := translation(s, chapters)
		if err != nil {
			fatal(fmt.Errorf("%s: %w", s.ID, err))
		}
		file := filepath.Join(*out, s.ID+".json.gz")
		if err := writeGZ(file, t); err != nil {
			fatal(err)
		}
		verses := 0
		for _, b := range t.Books {
			for _, ch := range b.Chapters {
				verses += len(ch)
			}
		}
		fmt.Printf("✓ %-9s %3d books %5d verses %8d bytes gz\n", s.ID, len(t.Books), verses, mustSize(file))
		if *only == "" {
			index = append(index, t.Info)
			licenses.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n", s.ID, s.Name, s.Language, s.License, s.Source))
		}
	}
	if *only != "" {
		return // a single redo touches its file alone
	}
	writeIndex(*out, index)
	os.WriteFile(filepath.Join(*out, "LICENSES.md"), []byte(licenses.String()), 0o644)
	fmt.Printf("wrote %s/index.json and LICENSES.md\n", *out)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "genbible:", err)
	os.Exit(1)
}

func mustSize(path string) int64 {
	fi, err := os.Stat(path)
	if err != nil {
		return -1
	}
	return fi.Size()
}

// fetch brings the upstream file to the cache, once: a file that is there
// stays, so the command can run again cheaply.
func fetch(url, local string) error {
	if fi, err := os.Stat(local); err == nil && fi.Size() > 0 {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(local), 0o755); err != nil {
		return err
	}
	tmp := local + ".part"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	defer f.Close()
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("fetching %s: %s", url, res.Status)
	}
	if _, err := io.Copy(f, res.Body); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, local)
}

// translation shapes what a parser found into a Translation, in the canon's
// order, and refuses one that misses a book or a chapter: better the build
// fails than a Bible ships incomplete.
func translation(s source, chapters map[string][][]string) (*bible.Translation, error) {
	t := &bible.Translation{Info: bible.Info{
		ID: s.ID, Name: s.Name, Short: s.Short, Language: s.Language, License: s.License, Source: s.Source,
	}}
	for _, code := range bible.CanonCodes {
		chs, ok := chapters[code]
		if !ok {
			return nil, fmt.Errorf("no %s", code)
		}
		want := bible.Canon[code].Chapters
		if extra := chapterExtras[s.ID][code]; extra != 0 {
			want = extra
		}
		if len(chs) > want {
			// A translation that keeps the deuterocanonical additions in the
			// book's own chapters, as the Synodal Daniel's 13 and 14 are:
			// the canon's chapters are what the app reads.
			fmt.Printf("  %-9s %s: %d chapters, the canon's %d kept\n", s.ID, code, len(chs), want)
			chs = chs[:want]
		}
		if len(chs) != want {
			return nil, fmt.Errorf("%s has %d chapters, want %d", code, len(chs), want)
		}
		name, _ := bible.BookName(code)
		t.Books = append(t.Books, &bible.Book{Code: code, Name: name.En, Chapters: chs})
	}
	return t, nil
}

func writeGZ(path string, v any) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	gz := gzip.NewWriter(f)
	if err := json.NewEncoder(gz).Encode(v); err != nil {
		return err
	}
	if err := gz.Close(); err != nil {
		return err
	}
	return f.Close()
}

func writeIndex(out string, index []bible.Info) {
	b, err := json.MarshalIndent(index, "", " ")
	if err != nil {
		fatal(err)
	}
	if err := os.WriteFile(filepath.Join(out, "index.json"), b, 0o644); err != nil {
		fatal(err)
	}
}

// ---- the parsers, one per format the upstreams come in ----

func parse(kind, path string) (map[string][][]string, error) {
	switch kind {
	case "usfm":
		return parseUSFMZip(path)
	case "getbible", "svdjson":
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		if kind == "getbible" {
			return parseGetBible(f)
		}
		return parseSVDJSON(f)
	case "usfx", "zefania", "osis":
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		switch kind {
		case "usfx":
			return parseUSFX(f)
		case "zefania":
			return parseZefania(f)
		default:
			return parseOSIS(f)
		}
	default:
		return nil, fmt.Errorf("unknown kind %q", kind)
	}
}

// ---- USFM: a zip of book files of markers, as ebible.org ships ----

var (
	usfmFile   = regexp.MustCompile(`^\d{2}-([A-Z0-9]{3})`)
	usfmEnd    = regexp.MustCompile(`\\[a-z0-9+]+\*`)
	usfmMarker = regexp.MustCompile(`\\[a-z0-9+]+(-[a-z0-9]+)?`)
	usfmAttr   = regexp.MustCompile(`\|[^\\\n]*`)
	// usfmNotes take out the study matter a verse carries: a footnote (f or
	// fe) and a cross-reference (x), each to its own end marker, since the
	// content itself carries nested word markers. The wording markers (+w,
	// add and their kin) stay.
	usfmNotes = []*regexp.Regexp{
		regexp.MustCompile(`(?s)\\f .{0,999}?\\f\*`),
		regexp.MustCompile(`(?s)\\fe .{0,999}?\\fe\*`),
		regexp.MustCompile(`(?s)\\x .{0,999}?\\x\*`),
	}
	usfmSpaces = regexp.MustCompile(`\s{2,}`)
)

// parseUSFMZip reads a zip of USFM book files, as ebible.org ships them:
// 70-MATeng-kjv.usfm holds Matthew, its verses under \v.
func parseUSFMZip(path string) (map[string][][]string, error) {
	zr, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer zr.Close()
	chapters := map[string][][]string{}
	for _, zf := range zr.File {
		m := usfmFile.FindStringSubmatch(zf.Name)
		if m == nil {
			continue // front matter, directory entries
		}
		code := m[1]
		if _, canon := bible.Canon[code]; !canon {
			continue
		}
		rc, err := zf.Open()
		if err != nil {
			return nil, err
		}
		b, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return nil, err
		}
		if chs := parseUSFMBook(string(b)); len(chs) > 0 {
			chapters[code] = chs
		}
	}
	return chapters, nil
}

// parseUSFMBook reads one book's markers: \c opens a chapter, \v opens a
// verse, the rest of a marker alphabet is either wording (\w, \add) the text
// keeps, or study matter (\f, \x, \e) it drops.
func parseUSFMBook(text string) [][]string {
	for _, note := range usfmNotes {
		text = note.ReplaceAllString(text, "")
	}
	var chs [][]string
	var cur []string
	flush := func() {
		if cur != nil {
			chs = append(chs, cur)
			cur = nil
		}
	}
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, `\c `) {
			if _, ok := atoi(line[3:]); ok {
				flush()
			}
			continue
		}
		if strings.HasPrefix(line, `\v `) {
			nstr, text, _ := strings.Cut(line[3:], " ")
			if _, ok := atoi(nstr); ok {
				cur = append(cur, cleanVerse(text))
			}
			continue
		}
	}
	flush()
	return chs
}

// cleanVerse tidies a verse: the markers' attributes, ends and marks go,
// the spaces they left collapse, and the ends lose their padding.
func cleanVerse(s string) string {
	s = usfmAttr.ReplaceAllString(s, "")
	s = usfmEnd.ReplaceAllString(s, "")
	s = usfmMarker.ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, "¶", "")
	s = usfmSpaces.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

func atoi(s string) (int, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, false
		}
		n = n*10 + int(c-'0')
	}
	return n, n > 0
}

// ---- USFX: XML with milestone verses, as <v id="1"/>text<ve/> ----

// verseDoc is the state the XML walkers share: the book and its chapters
// under way, and the verse collecting.
type verseDoc struct {
	chapters  map[string][][]string
	code      string
	chs       [][]string
	cur       *strings.Builder
	depthNote int
}

func (d *verseDoc) openBook(code string) {
	d.code = code
	d.chs = nil
}

func (d *verseDoc) openVerse() {
	if len(d.chs) == 0 {
		d.chs = append(d.chs, nil)
	}
	d.cur = &strings.Builder{}
}

func (d *verseDoc) closeVerse() {
	if d.cur != nil && len(d.chs) > 0 {
		last := len(d.chs) - 1
		d.chs[last] = append(d.chs[last], strings.TrimSpace(d.cur.String()))
		d.cur = nil
	}
}

func (d *verseDoc) closeBook() {
	if d.code != "" {
		d.chapters[d.code] = d.chs
		d.code, d.chs = "", nil
	}
}

// textStarters are the elements whose content is a heading or a note, which
// a verse's text never has.
var textStarters = map[string]bool{"note": true, "h": true, "title": true, "mt": true, "ms": true, "mr": true}

// parseUSFX walks the XML: a book by id, a chapter by <c>, a verse from <v>
// to <ve>, dropping headings and notes.
func parseUSFX(r io.Reader) (map[string][][]string, error) {
	dec := xml.NewDecoder(r)
	d := &verseDoc{chapters: map[string][][]string{}}
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "book":
				d.openBook(attr(t, "id"))
			case "c":
				d.chs = append(d.chs, nil)
			case "v":
				d.openVerse()
			case "ve":
				d.closeVerse()
			default:
				if textStarters[t.Name.Local] {
					d.depthNote++
				}
			}
		case xml.EndElement:
			switch t.Name.Local {
			case "book":
				d.closeBook()
			default:
				if textStarters[t.Name.Local] && d.depthNote > 0 {
					d.depthNote--
				}
			}
		case xml.CharData:
			if d.cur != nil && d.depthNote == 0 {
				d.cur.Write(t)
			}
		}
	}
	return d.chapters, nil
}

// ---- Zefania: XML with <BIBLEBOOK bnumber>, <CHAPTER cnumber>, <VERS> ----

// parseZefania walks the XML: the book's bnumber is its place in the canon,
// so GEN is 1 and REV is 66.
func parseZefania(r io.Reader) (map[string][][]string, error) {
	dec := xml.NewDecoder(r)
	d := &verseDoc{chapters: map[string][][]string{}}
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "BIBLEBOOK":
				n, ok := atoi(attr(t, "bnumber"))
				if !ok || n < 1 || n > len(bible.CanonCodes) {
					d.code = ""
					continue
				}
				d.openBook(bible.CanonCodes[n-1])
			case "CHAPTER":
				d.chs = append(d.chs, nil)
			case "VERS":
				d.openVerse()
			default:
				if t.Name.Local == "CAPTION" || t.Name.Local == "REMARK" || t.Name.Local == "NOTE" {
					d.depthNote++
				}
			}
		case xml.EndElement:
			switch t.Name.Local {
			case "BIBLEBOOK":
				d.closeBook()
			case "VERS":
				d.closeVerse()
			case "CAPTION", "REMARK", "NOTE":
				if d.depthNote > 0 {
					d.depthNote--
				}
			}
		case xml.CharData:
			if d.cur != nil && d.depthNote == 0 {
				d.cur.Write(t)
			}
		}
	}
	return d.chapters, nil
}

// ---- OSIS: XML whose verses carry osisIDs, as Gen.1.1 ----

// parseOSIS walks the XML and reads each verse from its osisID: the book's
// standard OSIS name turns into a code, and the chapter follows. Whatever
// containers the file wraps its verses in, the osisIDs are where the text
// belongs.
func parseOSIS(r io.Reader) (map[string][][]string, error) {
	dec := xml.NewDecoder(r)
	chapters := map[string][][]string{}
	var cur *strings.Builder
	var curBook string
	var curCh int
	depthNote := 0
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "verse":
				if book, ch, ok := osisRef(attr(t, "osisID")); ok {
					curBook, curCh = book, ch
					cur = &strings.Builder{}
				} else {
					curBook, cur = "", nil
				}
			case "note":
				depthNote++
			}
		case xml.EndElement:
			switch t.Name.Local {
			case "verse":
				if cur != nil && curBook != "" {
					chs := chaptersFor(chapters, curBook)
					for len(chs) < curCh {
						chs = append(chs, nil)
					}
					chs[curCh-1] = append(chs[curCh-1], strings.TrimSpace(cur.String()))
					chapters[curBook] = chs
					cur = nil
				}
			case "note":
				if depthNote > 0 {
					depthNote--
				}
			}
		case xml.CharData:
			if cur != nil && depthNote == 0 {
				cur.Write(t)
			}
		}
	}
	return chapters, nil
}

// chaptersFor is the chapter slice of a book in the accumulating set,
// making one as it's needed.
func chaptersFor(chapters map[string][][]string, code string) [][]string {
	chs := chapters[code]
	if chs == nil {
		chs = [][]string{}
		chapters[code] = chs
	}
	return chs
}

// osisRef reads an osisID as the app's book code and chapter, as "Gen.1"
// is GEN 1; a verse's ID has its verse after, which here goes unused.
func osisRef(id string) (book string, ch int, ok bool) {
	book, rest, _ := strings.Cut(id, ".")
	book = osisIDs[book]
	if book == "" {
		return "", 0, false
	}
	rest, _, _ = strings.Cut(rest, ".")
	ch, ok = atoi(rest)
	if !ok {
		return "", 0, false
	}
	return book, ch, true
}

func attr(t xml.StartElement, name string) string {
	for _, a := range t.Attr {
		if a.Name.Local == name {
			return a.Value
		}
	}
	return ""
}

// ---- the two JSON shapes: getbible.net's and the SVD file's ----

// getbibleBook is one book of api.getbible.net's shape: nr is the book's
// place in the canon, and the chapters each name their own verses.
type getbibleBook struct {
	Nr       int `json:"nr"`
	Chapters []struct {
		Chapter int `json:"chapter"`
		Verses  []struct {
			Verse int    `json:"verse"`
			Text  string `json:"text"`
		} `json:"verses"`
	} `json:"chapters"`
}

type getbibleDoc struct {
	Books []getbibleBook `json:"books"`
}

// parseGetBible reads api.getbible.net's whole-Bible JSON, as the Almeida
// and the Korean Revised Version come.
func parseGetBible(r io.Reader) (map[string][][]string, error) {
	var doc getbibleDoc
	if err := json.NewDecoder(r).Decode(&doc); err != nil {
		return nil, err
	}
	chapters := map[string][][]string{}
	for _, b := range doc.Books {
		if b.Nr < 1 || b.Nr > len(bible.CanonCodes) {
			continue
		}
		code := bible.CanonCodes[b.Nr-1]
		chs := make([][]string, len(b.Chapters))
		for i, ch := range b.Chapters {
			verses := make([]string, len(ch.Verses))
			for j, v := range ch.Verses {
				verses[j] = v.Text
			}
			chs[i] = verses
		}
		chapters[code] = chs
	}
	return chapters, nil
}

// svdAbbrevs maps the SVD file's book keys to the app's codes. Its own
// names, as the JSON has them: jd is Jude, jud is Judges, jo is John, and
// jn is Jonah.
var svdAbbrevs = map[string]string{
	"gn": "GEN", "ex": "EXO", "lv": "LEV", "nm": "NUM", "dt": "DEU",
	"js": "JOS", "jud": "JDG", "rt": "RUT", "1sm": "1SA", "2sm": "2SA",
	"1kgs": "1KI", "2kgs": "2KI", "1ch": "1CH", "2ch": "2CH", "ezr": "EZR",
	"ne": "NEH", "et": "EST", "job": "JOB", "ps": "PSA", "prv": "PRO",
	"ec": "ECC", "so": "SNG", "is": "ISA", "jr": "JER", "lm": "LAM",
	"ez": "EZK", "dn": "DAN", "ho": "HOS", "jl": "JOL", "am": "AMO",
	"ob": "OBA", "jn": "JON", "mi": "MIC", "na": "NAM", "hk": "HAB",
	"zp": "ZEP", "hg": "HAG", "zc": "ZEC", "ml": "MAL",
	"mt": "MAT", "mk": "MRK", "lk": "LUK", "jo": "JHN", "act": "ACT",
	"rm": "ROM", "1co": "1CO", "2co": "2CO", "gl": "GAL", "eph": "EPH",
	"ph": "PHP", "cl": "COL", "1ts": "1TH", "2ts": "2TH", "1tm": "1TI",
	"2tm": "2TI", "tt": "TIT", "phm": "PHM", "hb": "HEB", "jm": "JAS",
	"1pe": "1PE", "2pe": "2PE", "1jo": "1JN", "2jo": "2JN", "3jo": "3JN",
	"jd": "JUD", "re": "REV",
}

type svdDoc struct {
	Books map[string]struct {
		Name     string `json:"name"`
		Chapters map[string]map[string]string `json:"chapters"`
	} `json:"books"`
}

// parseSVDJSON reads the Arabic Smith & Van Dyke file: its books are keyed
// by their own short names, and its chapters and verses by their numbers,
// as JSON maps must be. The direction marks the file scatters through its
// text go.
func parseSVDJSON(r io.Reader) (map[string][][]string, error) {
	var doc svdDoc
	if err := json.NewDecoder(r).Decode(&doc); err != nil {
		return nil, err
	}
	marks := regexp.MustCompile(`[\x{200e}\x{200f}\x{202a}-\x{202e}\x{2066}-\x{2069}]`)
	chapters := map[string][][]string{}
	for abbr, b := range doc.Books {
		code, ok := svdAbbrevs[abbr]
		if !ok {
			continue
		}
		chs := make([][]string, len(b.Chapters))
		for chStr, verses := range b.Chapters {
			ch, ok := atoi(chStr)
			if !ok || ch < 1 || ch > len(chs) {
				continue
			}
			out := make([]string, len(verses))
			for vStr, text := range verses {
				v, ok := atoi(vStr)
				if !ok || v < 1 || v > len(out) {
					continue
				}
				out[v-1] = marks.ReplaceAllString(text, "")
			}
			chs[ch-1] = out
		}
		chapters[code] = chs
	}
	return chapters, nil
}
