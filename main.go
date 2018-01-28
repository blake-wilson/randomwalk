package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/llgcode/draw2d/draw2dimg"
)

func main() {
	// Initialize the graphic context on an RGBA image
	rect := image.Rect(0, 0, 850, 850)
	dest := image.NewRGBA(rect)
	gc := draw2dimg.NewGraphicContext(dest)

	// Set some properties
	gc.SetFillColor(color.RGBA{0x44, 0xff, 0x44, 0xff})
	gc.SetStrokeColor(color.RGBA{0x44, 0x44, 0x44, 0xff})
	gc.SetLineWidth(5)

	walk := NewRandomWalk(2)
	start := translateCoords(walk.position(), rect).coords
	// gc.MoveTo(float64(start[0]), float64(start[1]))
	// for i := 0; i < 2000; i++ {
	// 	p := walk.step()
	// 	fmt.Printf("walk is now at (%d, %d)\n", p.coords[0], p.coords[1])

	// 	p.coords[0] *= 10
	// 	p.coords[1] *= 10

	// 	ip := translateCoords(p, rect)
	// 	fmt.Printf("drawing at (%d, %d)\n", ip.coords[0], ip.coords[1])
	// 	gc.LineTo(float64(ip.coords[0]), float64(ip.coords[1]))
	// 	gc.FillStroke()
	// 	gc.MoveTo(float64(ip.coords[0]), float64(ip.coords[1]))
	// }
	// gc.Close()

	// Save to file
	// draw2dimg.SaveToPngFile("randomwalk.png", dest)

	http.HandleFunc("/step", func(w http.ResponseWriter, r *http.Request) {
		gc := draw2dimg.NewGraphicContext(dest)

		// Set some properties
		gc.SetFillColor(color.RGBA{0x44, 0xff, 0x44, 0xff})
		gc.SetStrokeColor(color.RGBA{0x44, 0x44, 0x44, 0xff})
		gc.SetLineWidth(5)

		gc.MoveTo(float64(start[0]), float64(start[1]))

		p := walk.step()

		for i := 0; i < len(walk.walk); i++ {

			fmt.Printf("walk is now at (%d, %d)\n", p.coords[0], p.coords[1])

			pos := copyPoint(walk.walk[i])
			pos.coords[0] *= 50
			pos.coords[1] *= 50

			ip := translateCoords(pos, rect)
			fmt.Printf("drawing at (%d, %d)\n", ip.coords[0], ip.coords[1])
			gc.LineTo(float64(ip.coords[0]), float64(ip.coords[1]))
			gc.FillStroke()
			gc.MoveTo(float64(ip.coords[0]), float64(ip.coords[1]))
		}

		gc.Close()
		img := dest.SubImage(rect)
		writeImage(w, &img)
	})

	http.ListenAndServe(`:8080`, nil)
}

// translate the given coordinate for the walk
// and translate it to a point on the image.
// The origin point should be in the middle of the image.
func translateCoords(p point, rect image.Rectangle) point {
	imgPoint := newPoint(p.dim)
	imgPoint.coords[0] = p.coords[0] + (rect.Dx())/2
	imgPoint.coords[1] = p.coords[1] + (rect.Dy())/2
	return imgPoint
}

type point struct {
	dim    int
	coords []int
}

func newPoint(dims int) point {
	return point{
		dim:    dims,
		coords: make([]int, dims),
	}
}

func copyPoint(p point) point {
	cpy := newPoint(p.dim)
	copy(cpy.coords, p.coords)
	return cpy
}

type randomWalk struct {
	dim  int
	walk []point
	rng  *rand.Rand
}

func NewRandomWalk(dims int) *randomWalk {
	initialValue := newPoint(dims)
	return &randomWalk{
		dim:  dims,
		walk: []point{initialValue},
		rng:  rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Get the last position of the walk
func (w *randomWalk) position() point {
	return w.walk[len(w.walk)-1]
}

func (w *randomWalk) randomStep() int {
	return ((w.rng.Int()%2)+1)*2 - 3
}

func (w *randomWalk) step() point {
	next := newPoint(w.dim)
	copy(next.coords, w.walk[len(w.walk)-1].coords)

	idx := w.rng.Int() % w.dim
	next.coords[idx] += w.randomStep()

	w.walk = append(w.walk, next)

	return copyPoint(next)
}

// writeImage encodes an image 'img' in png format and writes it into ResponseWriter.
func writeImage(w http.ResponseWriter, img *image.Image) {

	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, *img); err != nil {
		log.Println("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
}
