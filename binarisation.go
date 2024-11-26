package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
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

func (a *App) HandleBinarizeMeanIterative(base64str string, maxIterations int) string {
	if maxIterations < 0 || maxIterations > 100 {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Message: "Max iterations count must be between 0 and 100",
		})
		return ""
	}

	m, err := decodeBasePngToImg(base64str, a.ctx)
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Message: "Failed to decode image",
		})
		return ""
	}

	newM := binalizeMeanIterative(m, maxIterations)
	var buf bytes.Buffer
	if err := png.Encode(&buf, newM); err != nil {
		return ""
	}

	base64str = base64.StdEncoding.EncodeToString(buf.Bytes())
	return fmt.Sprintf("data:image/png;base64,%s", base64str)
}

func binalizeMeanIterative(m image.Image, maxIterations int) image.Image {
	b := m.Bounds()
	pxCount := b.Dx() * b.Dy()
	lums := make([]uint8, 0, pxCount)
	var sum uint64
	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			oldPx := m.At(x, y)
			lum := color.GrayModel.Convert(oldPx).(color.Gray).Y
			lums = append(lums, lum)
			sum += uint64(lum)
		}
	}
	threshold := uint8(sum / uint64(len(lums)))
	for i := 0; i < maxIterations; i++ {
		var sumLow, countLow, sumHigh, countHigh uint64
		oldThreshold := threshold
		for _, lum := range lums {
			if lum <= threshold {
				sumLow += uint64(lum)
				countLow++
			} else {
				sumHigh += uint64(lum)
				countHigh++
			}
		}
		if countLow > 0 && countHigh > 0 {
			meanLow := uint8(sumLow / countLow)
			meanHigh := uint8(sumHigh / countHigh)
			threshold = uint8((uint64(meanLow) + uint64(meanHigh)) / 2)
		}
		if threshold == oldThreshold {
			break
		}
	}

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

func (a *App) HandleBinarizeOtsu(base64str string) string {
	m, err := decodeBasePngToImg(base64str, a.ctx)
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Message: "Failed to decode image",
		})
		return ""
	}

	newM := binarizeOtsu(m)
	var buf bytes.Buffer
	if err := png.Encode(&buf, newM); err != nil {
		return ""
	}

	base64str = base64.StdEncoding.EncodeToString(buf.Bytes())
	return fmt.Sprintf("data:image/png;base64,%s", base64str)
}

func binarizeOtsu(m image.Image) image.Image {
	b := m.Bounds()
	histogram := make([]int, 256)
	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			oldPx := m.At(x, y)
			lum := color.GrayModel.Convert(oldPx).(color.Gray).Y
			histogram[lum]++
		}
	}
	totalPixels := b.Dx() * b.Dy()
	sumTotal := 0
	for i := 0; i < 256; i++ {
		sumTotal += i * histogram[i]
	}
	sumB := 0
	wB := 0
	maxVariance := 0.0
	threshold := uint8(0)

	for t := 0; t < 256; t++ {
		wB += histogram[t]
		if wB == 0 {
			continue
		}
		wF := totalPixels - wB
		if wF == 0 {
			break
		}
		sumB += t * histogram[t]
		mB := float64(sumB) / float64(wB)
		mF := float64(sumTotal-sumB) / float64(wF)

		variance := float64(wB) * float64(wF) * (mB - mF) * (mB - mF)
		if variance > maxVariance {
			maxVariance = variance
			threshold = uint8(t)
		}
	}
	binaryM := image.NewGray(b)
	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			oldPx := m.At(x, y)
			lum := color.GrayModel.Convert(oldPx).(color.Gray).Y
			if lum > threshold {
				binaryM.Set(x, y, color.Black)
			} else {
				binaryM.Set(x, y, color.White)
			}
		}
	}
	return binaryM
}

func (a *App) HandleBinarizeNiblack(base64str string, windowSize int, k float64) string {
	if windowSize%2 == 0 || windowSize < 3 {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Message: "Window size must be odd and >= 3",
		})
		return ""
	}

	m, err := decodeBasePngToImg(base64str, a.ctx)
	if err != nil {
		return ""
	}

	newM := binarizeNiblack(m, windowSize, k)
	var buf bytes.Buffer
	if err := png.Encode(&buf, newM); err != nil {
		return ""
	}

	base64str = base64.StdEncoding.EncodeToString(buf.Bytes())
	return fmt.Sprintf("data:image/png;base64,%s", base64str)
}

func binarizeNiblack(m image.Image, windowSize int, k float64) image.Image {
	b := m.Bounds()
	padding := windowSize / 2
	binaryM := image.NewGray(b)

	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			var sum, sumSq float64
			var count int
			for wy := -padding; wy <= padding; wy++ {
				for wx := -padding; wx <= padding; wx++ {
					nx, ny := x+wx, y+wy
					if nx >= 0 && nx < b.Dx() && ny >= 0 && ny < b.Dy() {
						px := m.At(nx, ny)
						lum := color.GrayModel.Convert(px).(color.Gray).Y
						sum += float64(lum)
						sumSq += float64(lum) * float64(lum)
						count++
					}
				}
			}
			mean := sum / float64(count)
			variance := (sumSq / float64(count)) - (mean * mean)
			stdDev := math.Sqrt(variance)
			threshold := mean + k*stdDev
			px := m.At(x, y)
			lum := color.GrayModel.Convert(px).(color.Gray).Y
			if float64(lum) > threshold {
				binaryM.Set(x, y, color.Black)
			} else {
				binaryM.Set(x, y, color.White)
			}
		}
	}

	return binaryM
}

func (a *App) HandleBinarizeBernsen(
	base64str string,
	windowSize int,
	contrastThreshold uint8,
) string {
	m, err := decodeBasePngToImg(base64str, a.ctx)
	if err != nil {
		return ""
	}

	newM := binarizeBernsen(m, windowSize, contrastThreshold)
	var buf bytes.Buffer
	if err := png.Encode(&buf, newM); err != nil {
		return ""
	}

	base64str = base64.StdEncoding.EncodeToString(buf.Bytes())
	return fmt.Sprintf("data:image/png;base64,%s", base64str)
}

func binarizeBernsen(
	m image.Image,
	windowVal int,
	contrastThreshold uint8,
) image.Image {
	padding := windowVal / 2
	b := m.Bounds()
	binM := image.NewGray(b)
	minInt := func(a, b int) int {
		if a > b {
			return b
		}
		return a
	}
	maxInt := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			var (
				minLum uint8 = 255
				maxLum uint8 = 0
			)
			yMin := maxInt(y-padding, 0)
			yMax := minInt(y+padding, b.Dy()-1)
			xMin := maxInt(x-padding, 0)
			xMax := minInt(x+padding, b.Dx()-1)
			for winY := yMin; winY <= yMax; winY++ {
				for winX := xMin; winX <= xMax; winX++ {
					px := m.At(winX, winY)
					lum := color.GrayModel.Convert(px).(color.Gray).Y
					if lum < minLum {
						minLum = lum
					}
					if lum > maxLum {
						maxLum = lum
					}
				}
			}
			contrast := maxLum - minLum
			if contrast >= contrastThreshold {
				threshold := minLum/2 + maxLum/2
				px := m.At(x, y)
				lum := color.GrayModel.Convert(px).(color.Gray).Y
				if lum < threshold {
					binM.Set(x, y, color.Black)
				} else {
					binM.Set(x, y, color.White)
				}
			} else {
				binM.Set(x, y, color.White)
			}
		}
	}

	return binM
}
