// Codie Bible: read the Bible, keep its seasons, and let a paired codie
// device read it aloud.
package main

import (
	"log"

	"codie-bibel/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
