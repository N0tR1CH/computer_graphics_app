package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) HandleHistogram(base64str string) string {
	var newM image.Image
	m, err := decodeBasePngToImg(base64str, a.ctx)
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:          runtime.WarningDialog,
			Title:         "Problem with parsing image data",
			Message:       fmt.Sprintf("Problem: %s", err.Error()),
			DefaultButton: "Ok",
		})
	}
	selection, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Title:   "Filters",
		Message: "Choose a filter you want to apply",
		Buttons: []string{"stretch", "equalize"},
	})
	switch selection {
	case "stretch":
		newM = stretchHistogram(m)
	case "equalize":
		newM = equalizeHistogram(m)
	default:
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:          runtime.WarningDialog,
			Title:         "Problem with parsing image data",
			Message:       fmt.Sprintf("Problem: %s", err.Error()),
			DefaultButton: "Ok",
		})
		return ""
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, newM); err != nil {
		return ""
	}
	base64str = base64.StdEncoding.EncodeToString(buf.Bytes())
	dataUrl := fmt.Sprintf("data:image/png;base64,%s", base64str)
	return dataUrl
}

func equalizeHistogram(img image.Image) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	histogram := make([]int, 256)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			gray := uint8(0.299*float64(r/257) + 0.587*float64(g/257) + 0.114*float64(b/257))
			histogram[gray]++
		}
	}
	cdf := make([]int, 256)
	cdf[0] = histogram[0]
	for i := 1; i < 256; i++ {
		cdf[i] = cdf[i-1] + histogram[i]
	}
	cdfMin := 0
	for i := 0; i < 256; i++ {
		if cdf[i] > 0 {
			cdfMin = cdf[i]
			break
		}
	}
	lookupTable := make([]uint8, 256)
	for i := 0; i < 256; i++ {
		lookupTable[i] = uint8(float64(cdf[i]-cdfMin) / float64(width*height-cdfMin) * 255.0)
	}
	equalized := image.NewGray(bounds)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			gray := uint8(0.299*float64(r/257) + 0.587*float64(g/257) + 0.114*float64(b/257))
			equalized.SetGray(x, y, color.Gray{Y: lookupTable[gray]})
		}
	}
	return equalized
}

func stretchHistogram(img image.Image) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	var min, max uint8 = 255, 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			gray := uint8(0.299*float64(r/257) + 0.587*float64(g/257) + 0.114*float64(b/257))
			if gray < min {
				min = gray
			}
			if gray > max {
				max = gray
			}
		}
	}
	if max == min {
		return img
	}
	scale := 255.0 / float64(max-min)
	stretched := image.NewGray(bounds)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			gray := uint8(0.299*float64(r/257) + 0.587*float64(g/257) + 0.114*float64(b/257))
			newGray := uint8(float64(gray-min) * scale)
			stretched.SetGray(x, y, color.Gray{Y: newGray})
		}
	}
	return stretched
}
