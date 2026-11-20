// Codie Bible: read the Bible, keep its seasons, and let a paired codie
// device read it aloud.
package main

import (
	"embed"
	"io/fs"
	"log"

	"codie-bibel/internal/app"
	"codie-bibel/internal/bible"
)

// bibles are the Bibles the app ships with, one gzipped JSON a translation,
// as cmd/genbible generated them into internal/bible/data.
//
//go:embed internal/bible/data
var bibles embed.FS

func main() {
	data, err := fs.Sub(bibles, "internal/bible/data")
	if err != nil {
		log.Fatal(err)
	}
	biblesReg, err := bible.NewRegistry(data)
	if err != nil {
		log.Fatal(err)
	}
	if err := app.Run(biblesReg); err != nil {
		log.Fatal(err)
	}
}
