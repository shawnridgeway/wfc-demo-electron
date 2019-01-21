package main

import (
	"encoding/json"
	// "fmt"
	// "reflect"
	// "io/ioutil"
	// "os"
	// "os/user"
	// "path/filepath"
	// "sort"
	// "strconv"
	"image"
	"image/color"

	// "github.com/asticode/go-astichartjs"
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
)

type Img struct {
	Data     []uint `json:"data"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Version  int    `json:"version"`
	Animated bool   `json:"animated"`
}

func (img Img) ColorModel() color.Model {
	return color.RGBAModel
}

func (img Img) Bounds() image.Rectangle {
	return image.Rect(0, 0, img.Width, img.Height)
}

func (img Img) At(x, y int) color.Color {
	ind := 4 * (x + y*img.Width)
	return color.RGBA{uint8(img.Data[ind]), uint8(img.Data[ind+1]), uint8(img.Data[ind+2]), uint8(img.Data[ind+3])}
}

type CancelMessagePayload struct {
	Version int `json:"version"`
}

var model *WfcGen

// handleMessages handles messages
func HandleMessages(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	switch m.Name {
	case "new":
		imgRequest := Img{}
		err := json.Unmarshal([]byte(m.Payload), &imgRequest)
		if err != nil {
			panic(err)
		}
		if model == nil {
			model = &WfcGen{}
		}
		if imgRequest.Animated {
			imgResponse, ok := model.Iterate(imgRequest, imgRequest.Version)
		} else {
			imgResponse, ok := model.Generate(imgRequest, imgRequest.Version)
		}
		if ok {
			payload = Img{imgToArray(imgResponse), imgResponse.Bounds().Max.X, imgResponse.Bounds().Max.Y, imgRequest.Version, imgRequest.Animated}
		}
	case "cancel-new":
		cancelRequest := CancelMessagePayload{}
		err := json.Unmarshal([]byte(m.Payload), &cancelRequest)
		if err != nil {
			panic(err)
		}
		model.Cancel(cancelRequest.Version)
		payload = CancelMessagePayload{cancelRequest.Version}
	}
	return
}

func imgToArray(img image.Image) []uint {
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y
	result := make([]uint, 4*width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			ind := 4 * (x + y*width)
			r, g, b, a := img.At(x, y).RGBA()
			result[ind] = uint(uint8(r))
			result[ind+1] = uint(uint8(g))
			result[ind+2] = uint(uint8(b))
			result[ind+3] = uint(uint8(a))
		}
	}
	return result
}
