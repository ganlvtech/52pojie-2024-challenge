package main

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
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

func MergePngAlpha(imageInputPath string, alphaInputPath string, outputPath string) error {
	imageInputNGRBA, err := OpenPngNRGBA(imageInputPath)
	if err != nil {
		return err
	}
	alphaInputNGRBA, err := OpenPngNRGBA(alphaInputPath)
	if err != nil {
		return err
	}
	if imageInputNGRBA.Bounds() != alphaInputNGRBA.Bounds() {
		return errors.New("input image bounds not match")
	}

	outputNRGBA := image.NewNRGBA(imageInputNGRBA.Bounds())
	for y := imageInputNGRBA.Bounds().Min.Y; y < imageInputNGRBA.Bounds().Max.Y; y++ {
		for x := imageInputNGRBA.Bounds().Min.X; x < imageInputNGRBA.Bounds().Max.X; x++ {
			imageNRGBA := imageInputNGRBA.NRGBAAt(x, y)
			alphaNRGBA := alphaInputNGRBA.NRGBAAt(x, y)
			outputNRGBA.SetNRGBA(x, y, color.NRGBA{
				R: imageNRGBA.R,
				G: imageNRGBA.G,
				B: imageNRGBA.B,
				A: alphaNRGBA.A,
			})
		}
	}

	err = SavePngNRGBA(outputPath, outputNRGBA)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	err := MergePngAlpha("flag10.png", "flag4.png", "output/flag4_flag10.png")
	if err != nil {
		panic(err)
	}
}
