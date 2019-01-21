package main

import (
	"fmt"
	"github.com/shawnridgeway/wavefunctioncollapse"
	"image"
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

	var result image.Image
	success := false
	wg.model = wavefunctioncollapse.NewOverlappingModel(inputImg, 3, destWidth, destHeight, true, true, 2, false)
	for !success {
		result, success = wg.model.Generate()
		if thisVersion != wg.version {
			fmt.Printf("Version: %v / %v | Canceled.\n", thisVersion, wg.version)
			return result, false
		}
		if success {
			fmt.Printf("Version: %v / %v | Succeeded.\n", thisVersion, wg.version)
			break
		} else {
			fmt.Printf("Version: %v / %v | Failed, trying again.\n", thisVersion, wg.version)
		}
	}
	return result, true
}

func (wg *WfcGen) Iterate() (image.Image, bool) {
	if thisVersion < wg.version {
		return nil, false
	}
	wg.version = thisVersion

	var result image.Image
	success := false
	wg.model = wavefunctioncollapse.NewOverlappingModel(inputImg, 3, destWidth, destHeight, true, true, 2, false)
}

func IterateRoutine(c chan int) {

}

func (wg *WfcGen) Cancel(thisVersion int) {
	if thisVersion > wg.version {
		wg.version = thisVersion
	}
}
