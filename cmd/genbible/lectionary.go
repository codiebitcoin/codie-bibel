package main

// The Revised Common Lectionary: genbible turns the reference table the
// Date-Lectionary project keeps (BSD 2-Clause, marmanold) into the app's
// own lectionary.json — the readings of each Sunday and holy day, keyed by
// a code the app computes from the calendar, in the three years and the
// fixed holy days.

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const rclURL = "https://raw.githubusercontent.com/marmanold/Date-Lectionary/master/share/rcl_lect.xml"

// lectXML is the shape of the source table.
type lectXML struct {
	Years []struct {
		Name string `xml:"name,attr"`
		Days []struct {
			Name    string `xml:"name,attr"`
			Lessons []struct {
				Order int    `xml:"order,attr"`
				Text  string `xml:",chardata"`
			} `xml:"lesson"`
		} `xml:"day"`
	} `xml:"year"`
}

// lectOut is what the app reads: the readings by code, by year, and the
// date ranges the after-Pentecost Sundays go by.
type lectOut struct {
	Years map[string]map[string][4]string `json:"years"`
	// Ranges hold the after-Pentecost Sundays' date ranges, as
	// "proper-7-10" reads July 10 to July 16.
	Ranges map[string][2]string `json:"ranges"`
}

// lectCodes maps the source table's day names to the app's codes.
var lectCodes = map[string]string{
	"First Sunday of Advent":            "advent1",
	"Second Sunday of Advent":           "advent2",
	"Third Sunday of Advent":            "advent3",
	"Fourth Sunday of Advent":           "advent4",
	"Christmas, Proper I":               "christmas1",
	"Christmas, Proper II":              "christmas2",
	"Christmas, Proper III":             "christmas3",
	"First Sunday after Christmas Day":  "after-christmas1",
	"Second Sunday after Christmas Day": "after-christmas2",
	"Epiphany of the Lord":              "epiphany",
	"Baptism of the Lord":               "baptism",
	"First Sunday after the Epiphany":   "after-epiphany1",
	"Second Sunday after the Epiphany":  "after-epiphany2",
	"Third Sunday after the Epiphany":   "after-epiphany3",
	"Fourth Sunday after the Epiphany":  "after-epiphany4",
	"Fifth Sunday after the Epiphany":   "after-epiphany5",
	"Sixth Sunday after the Epiphany":   "after-epiphany6",
	"Seventh Sunday after the Epiphany": "after-epiphany7",
	"Eighth Sunday after the Epiphany":  "after-epiphany8",
	"Last Sunday after the Epiphany":    "last-after-epiphany",
	"Ash Wednesday":                     "ash-wednesday",
	"First Sunday in Lent":              "lent1",
	"Second Sunday in Lent":             "lent2",
	"Third Sunday in Lent":              "lent3",
	"Fourth Sunday in Lent":             "lent4",
	"Fifth Sunday in Lent":              "lent5",
	"Liturgy of the Palms":              "palm",
	"Liturgy of the Passion":            "palm",
	"Holy Thursday":                     "maundy",
	"Good Friday":                       "good-friday",
	"Holy Saturday":                     "holy-saturday",
	"Easter Vigil":                      "easter-vigil",
	"Easter Day":                        "easter-day",
	"Easter Evening":                    "easter-evening",
	"Second Sunday of Easter":           "easter2",
	"Third Sunday of Easter":            "easter3",
	"Fourth Sunday of Easter":           "easter4",
	"Fifth Sunday of Easter":            "easter5",
	"Sixth Sunday of Easter":            "easter6",
	"Seventh Sunday of Easter":          "easter7",
	"Ascension of the Lord":             "ascension",
	"Day of Pentecost":                  "pentecost",
	"Trinity Sunday":                    "trinity",
	"Christ the King":                   "christ-the-king",
	"All Saints":                        "all-saints",
	"Holy Name of Jesus":                "holy-name",
	"Presentaiton of the Lord":          "presentation", // the table's own spelling
	"Presentation of the Lord":          "presentation",
	"Annunciation of the Lord":          "annunciation",
	"Visitation of Mary to Elizabeth":   "visitation",
	"Holy Cross Day":                    "holy-cross",
}

// rangeDay parses "Sunday between July 10 and July 16 inclusive" into its
// code and its two ends.
func rangeDay(name string) (code string, from, to string, ok bool) {
	if !strings.HasPrefix(name, "Sunday between ") {
		return "", "", "", false
	}
	months := map[string]string{
		"May": "5", "June": "6", "July": "7", "August": "8",
		"September": "9", "October": "10", "November": "11",
	}
	rest := strings.TrimSuffix(strings.TrimPrefix(name, "Sunday between "), " inclusive")
	if !strings.Contains(rest, " and ") {
		return "", "", "", false
	}
	mid := strings.Index(rest, " and ")
	first, second := rest[:mid], rest[mid+len(" and "):]
	parse := func(s string) (string, bool) {
		fields := strings.Fields(s)
		if len(fields) != 2 {
			return "", false
		}
		m, ok := months[fields[0]]
		if !ok {
			return "", false
		}
		d, err := strconv.Atoi(fields[1])
		if err != nil || d < 1 || d > 31 {
			return "", false
		}
		return fmt.Sprintf("%s-%02d", m, d), true
	}
	from, ok1 := parse(first)
	to, ok2 := parse(second)
	if !ok1 || !ok2 {
		return "", "", "", false
	}
	code = "proper-" + from
	return code, from, to, true
}

func genLectionary(cacheDir, outDir string) error {
	local := filepath.Join(cacheDir, "rcl_lect.xml")
	if err := fetch(rclURL, local); err != nil {
		return err
	}
	f, err := os.Open(local)
	if err != nil {
		return err
	}
	defer f.Close()
	var src lectXML
	if err := xml.NewDecoder(f).Decode(&src); err != nil {
		return err
	}
	out := lectOut{Years: map[string]map[string][4]string{}, Ranges: map[string][2]string{}}
	for _, y := range src.Years {
		yearName := y.Name
		if yearName == "holidays" {
			yearName = "holydays"
		}
		for _, d := range y.Days {
			code := lectCodes[d.Name]
			from := ""
			to := ""
			if code == "" {
				var ok bool
				code, from, to, ok = rangeDay(d.Name)
				if !ok {
					continue // a day the app doesn't keep, as Holy Week's weekdays
				}
			}
			year, ok := out.Years[yearName]
			if !ok {
				year = map[string][4]string{}
				out.Years[yearName] = year
			}
			var lessons [4]string
			for _, l := range d.Lessons {
				if l.Order >= 1 && l.Order <= 4 {
					lessons[l.Order-1] = strings.Join(strings.Fields(l.Text), " ")
				}
			}
			year[code] = lessons
			if from != "" {
				out.Ranges[code] = [2]string{from, to}
			}
		}
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(out, "", " ")
	if err != nil {
		return err
	}
	path := filepath.Join(outDir, "lectionary.json")
	if err := os.WriteFile(path, b, 0o644); err != nil {
		return err
	}
	n := 0
	for _, days := range out.Years {
		n += len(days)
	}
	fmt.Printf("✓ lectionary %d days over %d years, %d ranges\n", n, len(out.Years), len(out.Ranges))
	return nil
}
