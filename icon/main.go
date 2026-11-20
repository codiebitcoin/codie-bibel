// Command icon draws the app's icons: the app icon (resources/icon.png),
// and the tray's two — the color one for the notification area, and the
// black one macOS tints (internal/app/tray). The mark is an open Bible
// with the codie's badge at its corner.
//
// Usage: go run ./icon
package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"path/filepath"
)

func main() {
	dir, _ := os.Getwd()
	write(filepath.Join(dir, "resources", "icon.png"), drawIcon(1024, true, true))
	write(filepath.Join(dir, "internal", "app", "tray", "tray.png"), drawIcon(44, false, true))
	write(filepath.Join(dir, "internal", "app", "tray", "tray-template.png"), drawIcon(44, false, false))
	log.Print("wrote resources/icon.png and the tray's two")
}

func write(path string, img image.Image) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		log.Fatal(err)
	}
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		log.Fatal(err)
	}
}

var (
	burgundy = color.RGBA{0x7a, 0x2e, 0x2e, 0xff}
	cream    = color.RGBA{0xf7, 0xf2, 0xe8, 0xff}
	ribbon   = color.RGBA{0xa8, 0x3a, 0x3a, 0xff}
	codie    = color.RGBA{0xe8, 0x8b, 0x2e, 0xff}
)

// drawIcon draws the mark at size: the rounded square of the app icon
// (frame), the open book, and the codie's badge at the corner (badge).
// The tray's template draws the book dark on transparent, for macOS to
// tint.
func drawIcon(size int, frame, badge bool) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	if frame {
		roundRect(img, size, size/5, burgundy)
	}
	page := cream
	line := burgundy
	if !frame {
		if badge {
			page = color.RGBA{0xf0, 0xf0, 0xf0, 0xff}
		} else {
			page = color.RGBA{0x10, 0x10, 0x10, 0xff}
		}
		line = page
	}
	s := float64(size)
	spine := s * 0.5
	top := s * 0.26
	bottom := s * 0.74
	outerTop := s * 0.18
	outerBottom := s * 0.78

	// The pages, leaning from the spine to their outer edges.
	quad(img, spine, top, s*0.10, outerTop, s*0.10, outerBottom, spine, bottom, page)
	quad(img, spine, top, s*0.90, outerTop, s*0.90, outerBottom, spine, bottom, page)
	// The lines of the text, four to a page.
	for i := 0; i < 4; i++ {
		y := s * (0.38 + 0.085*float64(i))
		h := s * 0.016
		rect(img, s*0.17, y, s*0.44, y+h, line)
		rect(img, s*0.56, y, s*0.83, y+h, line)
	}
	// The ribbon of the bookmark, from the spine's foot.
	rect(img, spine-s*0.03, bottom, spine+s*0.05, s*0.88, ribbon)
	if badge {
		r := s * 0.14
		cx, cy := s*0.80, s*0.80
		circle(img, cx, cy, r, codie)
		circle(img, cx, cy, r*0.58, burgundy)
		circle(img, cx, cy, r*0.28, codie)
	}
	return img
}

// roundRect fills a rounded rectangle over the whole square.
func roundRect(img *image.RGBA, size, radius int, col color.RGBA) {
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if inRoundRect(x, y, size, radius) {
				img.SetRGBA(x, y, col)
			}
		}
	}
}

func inRoundRect(x, y, size, r int) bool {
	corner := func(cx, cy int) bool {
		dx, dy := x-cx, y-cy
		return dx*dx+dy*dy > r*r
	}
	if x < r && y < r && corner(r, r) {
		return false
	}
	if x >= size-r && y < r && corner(size-r-1, r) {
		return false
	}
	if x < r && y >= size-r && corner(r, size-r-1) {
		return false
	}
	if x >= size-r && y >= size-r && corner(size-r-1, size-r-1) {
		return false
	}
	return true
}

// quad fills a four-cornered shape by scanline, corners in order.
func quad(img *image.RGBA, x1, y1, x2, y2, x3, y3, x4, y4 float64, col color.RGBA) {
	fillTriangle(img, x1, y1, x2, y2, x3, y3, col)
	fillTriangle(img, x1, y1, x3, y3, x4, y4, col)
}

// rect fills a rectangle.
func rect(img *image.RGBA, x1, y1, x2, y2 float64, col color.RGBA) {
	for y := int(y1); y < int(y2); y++ {
		for x := int(x1); x < int(x2); x++ {
			img.SetRGBA(x, y, col)
		}
	}
}

// fillTriangle fills a triangle by barycentric test over its bounds.
func fillTriangle(img *image.RGBA, x1, y1, x2, y2, x3, y3 float64, col color.RGBA) {
	minX := int(min3(x1, x2, x3))
	minY := int(min3(y1, y2, y3))
	maxX := int(max3(x1, x2, x3)) + 1
	maxY := int(max3(y1, y2, y3)) + 1
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			px, py := float64(x)+0.5, float64(y)+0.5
			if pointIn(px, py, x1, y1, x2, y2, x3, y3) {
				img.SetRGBA(x, y, col)
			}
		}
	}
}

func pointIn(px, py, x1, y1, x2, y2, x3, y3 float64) bool {
	d := (y2-y3)*(x1-x3) + (x3-x2)*(y1-y3)
	a := ((y2-y3)*(px-x3) + (x3-x2)*(py-y3)) / d
	b := ((y3-y1)*(px-x3) + (x1-x3)*(py-y3)) / d
	c := 1 - a - b
	return a >= 0 && b >= 0 && c >= 0
}

func circle(img *image.RGBA, cx, cy, r float64, col color.RGBA) {
	for y := int(cy - r); y <= int(cy+r); y++ {
		for x := int(cx - r); x <= int(cx+r); x++ {
			dx, dy := float64(x)-cx, float64(y)-cy
			if dx*dx+dy*dy <= r*r {
				img.SetRGBA(x, y, col)
			}
		}
	}
}

func min3(a, b, c float64) float64 {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

func max3(a, b, c float64) float64 {
	if a > b {
		if a > c {
			return a
		}
		return c
	}
	if b > c {
		return b
	}
	return c
}
