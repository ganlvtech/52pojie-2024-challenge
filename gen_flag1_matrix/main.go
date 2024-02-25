package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
)

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

func GenMatrixImage(bounds image.Rectangle, distance float64, intensityFunc func(r float64) float64) *image.NRGBA {
	outputNRGBA := image.NewNRGBA(bounds)
	centerX := float64(outputNRGBA.Bounds().Min.X+outputNRGBA.Bounds().Max.X) / 2.0
	centerY := float64(outputNRGBA.Bounds().Min.Y+outputNRGBA.Bounds().Max.Y) / 2.0
	for y := outputNRGBA.Bounds().Min.Y; y < outputNRGBA.Bounds().Max.Y; y++ {
		for x := outputNRGBA.Bounds().Min.X; x < outputNRGBA.Bounds().Max.X; x++ {
			dx1 := math.Mod(math.Abs(float64(x)-centerX), distance)
			dy1 := math.Mod(math.Abs(float64(y)-centerY), distance)
			dx2 := distance - dx1
			dy2 := distance - dy1

			grayF := intensityFunc(math.Hypot(dx1, dy1)) + intensityFunc(math.Hypot(dx1, dy2)) + intensityFunc(math.Hypot(dx2, dy1)) + intensityFunc(math.Hypot(dx2, dy2))

			var gray uint8
			if grayF >= 1.0 {
				gray = 0xff
			} else {
				gray = uint8(255.0 * grayF)
			}

			outputNRGBA.SetNRGBA(x, y, color.NRGBA{
				R: gray,
				G: gray,
				B: gray,
				A: 0xff,
			})
		}
	}
	return outputNRGBA
}

func main() {
	err := SavePngNRGBA("output/flag1_matrix.png", GenMatrixImage(image.Rect(0, 0, 3840, 2160), 40, func(r float64) float64 {
		if r <= 2.0 {
			return 1.0
		}
		if r <= 5.0 {
			return 1.0 - (r-2.0)/3.0
		}
		return 0.0
	}))
	if err != nil {
		panic(err)
	}
}
