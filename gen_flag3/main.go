package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

func OpenPngNRGBA(inputPath string) (*image.NRGBA, error) {
	inputBytes, err := os.ReadFile(inputPath)
	if err != nil {
		return nil, err
	}
	inputImage, err := png.Decode(bytes.NewReader(inputBytes))
	if err != nil {
		return nil, err
	}
	inputNGRBA, ok := inputImage.(*image.NRGBA)
	if !ok {
		return nil, errors.New("input image is not in rgba color format")
	}
	return inputNGRBA, nil
}

func SavePngNRGBA(outputPath string, img *image.NRGBA) error {
	err := os.MkdirAll(filepath.Dir(outputPath), 0755)
	if err != nil {
		return err
	}
	outputFile, err := os.OpenFile(outputPath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	err = png.Encode(outputFile, img)
	if err != nil {
		return err
	}
	return nil
}

func RunTasks(taskList []func()) {
	maxParallel := runtime.NumCPU()
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxParallel)
	for _, task := range taskList {
		wg.Add(1)
		semaphore <- struct{}{}
		go func(t func()) {
			defer func() {
				<-semaphore
				wg.Done()
			}()
			t()
		}(task)
	}
	wg.Wait()
}

func GenNoiseImage(seed int64, bonunds image.Rectangle) *image.NRGBA {
	r := rand.NewSource(seed)
	outputNRGBA := image.NewNRGBA(bonunds)

	for y := outputNRGBA.Bounds().Min.Y; y < outputNRGBA.Bounds().Max.Y; y += 2 {
		for x := outputNRGBA.Bounds().Min.X; x < outputNRGBA.Bounds().Max.X; x += 2 {
			gray := uint8(r.Int63() & 0xff)
			c := color.NRGBA{
				R: gray,
				G: gray,
				B: gray,
				A: 0xff,
			}
			outputNRGBA.SetNRGBA(x, y, c)
			outputNRGBA.SetNRGBA(x+1, y, c)
			outputNRGBA.SetNRGBA(x, y+1, c)
			outputNRGBA.SetNRGBA(x+1, y+1, c)
		}
	}
	return outputNRGBA
}

func GenFlag3(flagInputNRGBA *image.NRGBA, noiseInputNRGBA *image.NRGBA, outputPath string, delta int) error {
	if flagInputNRGBA.Bounds() != noiseInputNRGBA.Bounds() {
		return errors.New("input image bounds not match")
	}

	width := noiseInputNRGBA.Bounds().Max.X - noiseInputNRGBA.Bounds().Min.X
	outputNRGBA := image.NewNRGBA(noiseInputNRGBA.Bounds())
	for y := outputNRGBA.Bounds().Min.Y; y < outputNRGBA.Bounds().Max.Y; y++ {
		for x := outputNRGBA.Bounds().Min.X; x < outputNRGBA.Bounds().Max.X; x++ {
			flagNRGBA := flagInputNRGBA.NRGBAAt(x, y)
			var rgba color.NRGBA
			if flagNRGBA.R >= 128 {
				rgba = noiseInputNRGBA.NRGBAAt((x-delta+width)%width, y)
			} else {
				rgba = noiseInputNRGBA.NRGBAAt((x+delta)%width, y)
			}
			outputNRGBA.SetNRGBA(x, y, color.NRGBA{
				R: rgba.R,
				G: rgba.G,
				B: rgba.B,
				A: rgba.A,
			})
		}
	}

	err := SavePngNRGBA(outputPath, outputNRGBA)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	flagInputNRGBA, err := OpenPngNRGBA("flag3_image.png")
	if err != nil {
		panic(err)
	}
	noiseInputNRGBA := GenNoiseImage(0, flagInputNRGBA.Bounds())
	var taskList []func()
	for i := 0; i < 42; i++ {
		(func(i int) {
			taskList = append(taskList, func() {
				err := GenFlag3(flagInputNRGBA, noiseInputNRGBA, fmt.Sprintf("output/flag3_%02d.png", i), i*5)
				if err != nil {
					fmt.Println(i, err)
					panic(err)
				} else {
					fmt.Println(i, "OK")
				}
			})
		})(i)
	}
	RunTasks(taskList)
}
