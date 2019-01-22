package main

import (
	"fmt"
	"github.com/asticode/go-astilectron"
	"github.com/shawnridgeway/wavefunctioncollapse"
	"image"
	"time"
)

type WfcGen struct {
	model   *wavefunctioncollapse.OverlappingModel
	version int
	frame   int
}

const (
	destWidth, destHeight = 50, 50
)

func (wg *WfcGen) Generate(inputImg image.Image, thisVersion int) (image.Image, bool) {
	if thisVersion < wg.version {
		return nil, false
	}
	wg.version = thisVersion

	// var result image.Image
	wg.GetNewModel(inputImg)
	for {
		result, success := wg.model.Generate()
		if thisVersion != wg.version {
			fmt.Printf("Version: %v / %v | Canceled.\n", thisVersion, wg.version)
			return result, false
		}
		if success {
			fmt.Printf("Version: %v / %v | Succeeded.\n", thisVersion, wg.version)
			return result, true
		} else {
			fmt.Printf("Version: %v / %v | Failed.\n", thisVersion, wg.version)
			// Retry
		}
	}
	// fmt.Printf("Version: %v / %v | Unknown Result.\n", thisVersion, wg.version)
	// return result, false
}

func (wg *WfcGen) Iterate(inputImg image.Image, thisVersion int, win *astilectron.Window) (image.Image, bool) {
	if thisVersion < wg.version {
		return nil, false
	}
	wg.version = thisVersion

	wg.GetNewModel(inputImg)
	result, finished, success := wg.model.Iterate(1)
	if thisVersion != wg.version {
		fmt.Printf("Version: %v / %v | Canceled.\n", thisVersion, wg.version)
		return result, false
	}
	if finished {
		if success {
			fmt.Printf("Version: %v / %v | Succeeded.\n", thisVersion, wg.version)
			return result, true
		} else {
			fmt.Printf("Version: %v / %v | Failed.\n", thisVersion, wg.version)
			return result, true
		}
	}
	go func() {
		c := make(chan image.Image)
		quit := make(chan bool)
		restart := make(chan bool)
		go IterateRoutine(c, quit, restart, wg, inputImg, thisVersion)
		var imgResponse image.Image
		for {
			select {
			case imgResponse = <-c:
				// Send new frame
				message := Img{imgToArray(imgResponse), imgResponse.Bounds().Max.X, imgResponse.Bounds().Max.Y, thisVersion, true}
				win.SendMessage(message)
			case <-restart:
				wg.model.Clear()
				go IterateRoutine(c, quit, restart, wg, inputImg, thisVersion)
			case <-quit:
				return
			}
		}
	}()
	// Send first frame, go routine will continue the rest
	return result, true
}

func IterateRoutine(c chan image.Image, quit chan bool, restart chan bool, wg *WfcGen, inputImg image.Image, thisVersion int) {
	frame := 1
	iterationsPerFrame := 25
	for {
		frameStart := time.Now()

		result, finished, success := wg.model.Iterate(iterationsPerFrame)

		timeDiff := (10 * time.Millisecond) - time.Since(frameStart)
		if timeDiff > 0 {
			time.Sleep(timeDiff)
		}

		if thisVersion != wg.version {
			fmt.Printf("Version: %v / %v | Frames: %v | Canceled.\n", thisVersion, wg.version, frame)
			quit <- true
			return
		}
		if finished {
			if success {
				fmt.Printf("Version: %v / %v | Frames: %v | Succeeded.\n", thisVersion, wg.version, frame)
				c <- result
				quit <- true
				return
			} else {
				fmt.Printf("Version: %v / %v | Frames: %v | Failed.\n", thisVersion, wg.version, frame)
				c <- result
				restart <- true
				return
			}
		} else {
			// Send this frame, continue on to the next
			c <- result
		}
		frame += iterationsPerFrame
	}
}

func (wg *WfcGen) GetNewModel(inputImg image.Image) {
	wg.model = wavefunctioncollapse.NewOverlappingModel(inputImg, 3, destWidth, destHeight, true, true, 2, false)
}

func (wg *WfcGen) Cancel(thisVersion int) {
	if thisVersion > wg.version {
		wg.version = thisVersion
	}
}
