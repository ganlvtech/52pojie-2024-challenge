package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
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

func ApplyShader(img *image.NRGBA, shaderFunc func(nrgba color.NRGBA, x int, y int, bounds image.Rectangle) color.NRGBA) *image.NRGBA {
	bounds := img.Bounds()
	output := image.NewNRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			output.SetNRGBA(x, y, shaderFunc(img.NRGBAAt(x, y), x, y, bounds))
		}
	}
	return output
}

func main() {
	inputNRGBA, err := OpenPngNRGBA("flag1_image.png")
	if err != nil {
		panic(err)
	}
	centerX := float64(inputNRGBA.Bounds().Min.X+inputNRGBA.Bounds().Max.X) / 2.0
	centerY := float64(inputNRGBA.Bounds().Min.Y+inputNRGBA.Bounds().Max.Y) / 2.0
	waveHalfWidth := centerY / 5.0
	var taskList []func()
	for i := 0; i < 60; i++ {
		(func(i int) {
			taskList = append(taskList, func() {
				waveCenterRadius := float64(i) * (math.Hypot(centerX, centerY) + waveHalfWidth) / 60.0
				outputNRGBA := ApplyShader(inputNRGBA, func(nrgba color.NRGBA, x int, y int, bounds image.Rectangle) color.NRGBA {
					r := math.Hypot(float64(x)-centerX, float64(y)-centerY)
					if r < waveCenterRadius-waveHalfWidth { // 圆内部全透明
						return color.NRGBA{
							R: 0,
							G: 0,
							B: 0,
							A: 0,
						}
					} else if r < waveCenterRadius { // 圆内部边缘逐渐半透明
						return color.NRGBA{
							R: nrgba.R,
							G: nrgba.G,
							B: nrgba.B,
							A: uint8((1.0 - (waveCenterRadius-r)/waveHalfWidth) * 255.0),
						}
					} else if r < waveCenterRadius+waveHalfWidth { // 圆外部边缘黑色渐隐
						return color.NRGBA{
							R: uint8((1.0 - (r-waveCenterRadius)/waveHalfWidth) * float64(nrgba.R)),
							G: uint8((1.0 - (r-waveCenterRadius)/waveHalfWidth) * float64(nrgba.G)),
							B: uint8((1.0 - (r-waveCenterRadius)/waveHalfWidth) * float64(nrgba.B)),
							A: 255,
						}
					} else {
						return color.NRGBA{ // 圆外部纯黑色
							R: 0,
							G: 0,
							B: 0,
							A: 255,
						}
					}
				})
				err := SavePngNRGBA(fmt.Sprintf("output/flag1_%02d.png", i), outputNRGBA)
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
