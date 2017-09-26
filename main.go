package main

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
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
	gc.MoveTo(float64(start[0]), float64(start[1]))
	for i := 0; i < 2000; i++ {
		p := walk.step()
		fmt.Printf("walk is now at (%d, %d)\n", p.coords[0], p.coords[1])

		p.coords[0] *= 10
		p.coords[1] *= 10

		ip := translateCoords(p, rect)
		fmt.Printf("drawing at (%d, %d)\n", ip.coords[0], ip.coords[1])
		gc.LineTo(float64(ip.coords[0]), float64(ip.coords[1]))
		gc.FillStroke()
		gc.MoveTo(float64(ip.coords[0]), float64(ip.coords[1]))
	}
	gc.Close()

	// Save to file
	draw2dimg.SaveToPngFile("randomwalk.png", dest)
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
	for i := 0; i < w.dim; i++ {
		next.coords[i] += w.randomStep()
	}
	w.walk = append(w.walk, next)

	return copyPoint(next)
}
