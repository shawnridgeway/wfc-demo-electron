package main

import (
	"fmt"
	"os"
	"time"
	"math/rand"
	// "io"
	"image"
	"strconv"
	// "image/color"
	"image/draw"
	"image/png"
	"github.com/shawnridgeway/wfc/grid"
)

type Tile image.Image

func loadImg(file string) image.Image {
	raw, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer raw.Close()
	img, _, err := image.Decode(raw)
	if err != nil {
		panic(err)
	}
	return img
}

func saveImg(file, path string, img image.Image) {
	rErr := os.RemoveAll(path)
	if rErr != nil {
		fmt.Println(rErr)
	}
	mErr := os.MkdirAll(path, 0777)
	if mErr != nil {
		fmt.Println(mErr)
	}
	f, cErr := os.Create(path + file)
	if cErr != nil {
		panic(cErr)
	}
	defer f.Close()
	png.Encode(f, img)
}

func wfc(g *grid.Grid) {
	rand.Seed(time.Now().UnixNano())
	bounds := len(*g)
	frontier := make([]grid.Coor, 0)
	frontier = append(frontier, grid.Coor{X: int(bounds / 2), Y: int(bounds / 2)})
	for len(frontier) > 0 {
		// Get current tile
		current := frontier[0]
		frontier = frontier[1:]

		if (*g)[current.Y][current.X] != 0 {
			continue
		}

		// Get neighbors
		neighbors := make([]grid.Coor, 0)
		if current.X + 1 < bounds {
			neighbors = append(neighbors, grid.Coor{X: current.X + 1, Y: current.Y})
		}
		if current.X - 1 >= 0 {
			neighbors = append(neighbors, grid.Coor{X: current.X - 1, Y: current.Y})
		}
		if current.Y + 1 < bounds {
			neighbors = append(neighbors, grid.Coor{X: current.X, Y: current.Y + 1})
		}
		if current.Y - 1 >= 0 {
			neighbors = append(neighbors, grid.Coor{X: current.X, Y: current.Y - 1})
		}

		lands := 0
		waters := 0
		for _, n := range neighbors {
			if (*g)[n.Y][n.X] == 1 {
				lands++
			}
			if (*g)[n.Y][n.X] == 2 {
				waters++
			}
		}

		tile := 2
		if (lands * 2 > waters && rand.Intn(bounds) > 1) || (lands * 2 <= waters && rand.Intn(bounds) < 1) {
			tile = 1
		}

		// Set current tile
		(*g)[current.Y][current.X] = tile

		// Enquque non-explored neighbors
		for _, n := range neighbors {
			if (*g)[n.Y][n.X] == 0 {
				frontier = append(frontier, n)
			}
		}
	}
}

func Boring() <-chan int {
	c := make(chan int)
	go func() {
		for {
			c <- rand.Intn(10)
		}
	}()
	return c
}

func merge(c1, c2 <-chan int) <-chan int {
	out := make(chan int)
	go func() { for { out <- <-c1 } }()
	go func() { for { out <- <-c2 } }()
	return out
}

func TestBoring() {
	c1 := Boring()
	c2 := Boring()
	for v := range merge(c1, c2) {
		fmt.Println(v)
	}
}

func GetImg(staticPath string) string {
	tileSize := 16
	gridSize := 32
	staticPath += "/resources/app/"
	outputPath := "output/"
	fileName := "img_" + strconv.Itoa(int(time.Now().UnixNano())) + ".png" 

	// Load tiles
	images := make(map[int]image.Image)
	images[1] = loadImg("resources/img/grass.png")
	images[2] = loadImg("resources/img/water.png")

	// Make the grid
	g := grid.New(gridSize, gridSize)
	wfc(g)

	sig := 0

	// Make image
	img := image.NewRGBA(image.Rect(0, 0, tileSize * gridSize, tileSize * gridSize))
	for iy, row := range *g {
		for ix, t := range row {
			draw.Draw(img, image.Rect(tileSize * ix, tileSize * iy, tileSize * (ix + 1), tileSize * (iy + 1)), images[t], image.ZP, draw.Src)
			if t == 1 {
				sig++
			}
		}
	}

	// Save to file
	saveImg(fileName, staticPath + outputPath, img)

	return outputPath + fileName
}