package main

import (
	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"time"
)

const (
	rows = 30
	cols = 14

	blockSize       = 20
	blockLineThin   = 1
	blockMarginSize = 4

	screenWidth  = cols * blockSize
	screenHeight = rows * blockSize
)

var (
	blockColor = sdl.Color{84, 85, 80, 255}
	bgColor    = sdl.Color{0, 0, 0, 255}
)

type Point struct {
	x int
	y int
}

func filledRectangle(r *sdl.Renderer, x, y, w, h int, c sdl.Color) bool {
	lx := int16(x)
	rx := int16(x + w)
	ty := int16(y)
	by := int16(y + h)
	return gfx.FilledPolygonColor(r, []int16{lx, rx, rx, lx, lx}, []int16{ty, ty, by, by, ty}, c)
}

func drawBlock(r *sdl.Renderer, col, row int, invert bool) {
	bg := bgColor
	fg := blockColor

	if invert {
		bg, fg = fg, bg
	}

	filledRectangle(r, col*blockSize, row*blockSize, blockSize, blockSize, fg)
	filledRectangle(r, col*blockSize+blockLineThin, row*blockSize+blockLineThin, blockSize-blockLineThin, blockSize-blockLineThin, bg)
	filledRectangle(r, col*blockSize+blockMarginSize, row*blockSize+blockMarginSize, blockSize-blockMarginSize, blockSize-blockMarginSize, fg)
}

func main() {
	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_TIMER); err != nil {
		log.Fatal(err)
	}

	window, err := sdl.CreateWindow("gameboy",
		sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		screenWidth, screenHeight, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer window.Destroy()

	render, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatal(err)
	}
	defer render.Destroy()

	render.SetDrawColor(0, 0, 0, sdl.ALPHA_OPAQUE)
	render.Clear()

	for x := 0; x < cols; x++ {
		for y := 0; y < rows; y++ {
			drawBlock(render, x, y, false)
		}
	}
	render.Present()

	brick := []Point{Point{10, 10}, Point{11, 10}, Point{11, 11}, Point{11, 12}}
	offset := 0

	tick := time.Tick(500 * time.Millisecond)
	go func(c <-chan time.Time) {
		for {
			<-c
			for _, b := range brick {
				drawBlock(render, b.x, b.y+offset%rows, false)
			}

			for _, b := range brick {
				drawBlock(render, b.x, (b.y+offset+1)%rows, true)
			}

			offset = offset + 1
			render.Present()
		}
	}(tick)

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
	}

	sdl.Quit()
}
