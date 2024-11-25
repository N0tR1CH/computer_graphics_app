package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"slices"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) HandleBinarizeManual(base64str string, threshold uint8) string {
	if threshold < 0 || threshold > 255 {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Message: "Threshold must be between 0 and 255",
		})
		return ""
	}
	m, err := decodeBasePngToImg(base64str, a.ctx)
	if err != nil {
		return ""
	}
	newM := binarizeManual(m, threshold)
	var buf bytes.Buffer
	if err := png.Encode(&buf, newM); err != nil {
		return ""
	}
	base64str = base64.StdEncoding.EncodeToString(buf.Bytes())
	dataUrl := fmt.Sprintf("data:image/png;base64,%s", base64str)
	return dataUrl
}

func binarizeManual(m image.Image, threshold uint8) image.Image {
	bounds := m.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	newM := image.NewGray(bounds)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := m.At(x, y).RGBA()
			gray := uint8(0.299*float64(r/257) + 0.587*float64(g/257) + 0.114*float64(b/257))
			if gray > threshold {
				newM.Set(x, y, color.White)
			} else {
				newM.Set(x, y, color.Black)
			}
		}
	}
	return newM
}

func (a *App) HandleBinarizePercentBlack(base64str string, percent float64) string {
	if percent < 0 || percent > 100 {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Message: "Percent must be between 0 and 100",
		})
		return ""
	}
	m, err := decodeBasePngToImg(base64str, a.ctx)
	if err != nil {
		return ""
	}
	newM := binalizePercentBlack(m, percent)
	var buf bytes.Buffer
	if err := png.Encode(&buf, newM); err != nil {
		return ""
	}
	base64str = base64.StdEncoding.EncodeToString(buf.Bytes())
	dataUrl := fmt.Sprintf("data:image/png;base64,%s", base64str)
	return dataUrl
}

func binalizePercentBlack(m image.Image, percent float64) image.Image {
	b := m.Bounds()
	pxCount := b.Dx() * b.Dy()
	lums := make([]uint8, 0, pxCount)
	for y := 0; y < b.Max.Y; y++ {
		for x := 0; x < b.Max.X; x++ {
			oldPx := m.At(x, y)
			lum := color.GrayModel.Convert(oldPx).(color.Gray).Y
			lums = append(lums, lum)
		}
	}

	lumsSorted := make([]uint8, len(lums))
	copy(lumsSorted, lums)
	slices.Sort(lumsSorted)

	thresholdIdx := int((percent / 100) * float64(pxCount))
	threshold := lumsSorted[thresholdIdx]

	binaryM := image.NewGray(b)
	var loopIdx int
	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			if lums[loopIdx] > threshold {
				binaryM.Set(x, y, color.Black)
			} else {
				binaryM.Set(x, y, color.White)
			}
			loopIdx++
		}
	}
	return binaryM
}
