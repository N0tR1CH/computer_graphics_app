package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"sort"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) HandleFilterApplying(base64str string) string {
	m, err := decodeBasePngToImg(base64str, a.ctx)
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:          runtime.InfoDialog,
			Title:         "Decoding problem",
			Message:       fmt.Sprintf("Image could not be decoded: %s", err.Error()),
			DefaultButton: "Ok",
		})
	}
	selection, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Title:         "Filters",
		Message:       "Choose a filter you want to apply",
		Buttons:       []string{"average", "median", "sobel", "gaussian"},
		DefaultButton: "average",
	})

	var newM image.Image
	switch selection {
	case "average":
		newM = ApplyAveragingFilter(m)
	case "median":
		newM = ApplyMedianFilter(m)
	case "sobel":
		newM = ApplySobelFilter(m)
	case "gaussian":
		newM = ApplyGaussianBlur(m)
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, newM); err != nil {
		return ""
	}
	base64str = base64.StdEncoding.EncodeToString(buf.Bytes())
	dataUrl := fmt.Sprintf("data:image/png;base64,%s", base64str)
	return dataUrl
}

func ApplyAveragingFilter(img image.Image) image.Image {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var sumR, sumG, sumB int
			var count int

			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					nx, ny := x+dx, y+dy
					if nx < bounds.Min.X || nx >= bounds.Max.X || ny < bounds.Min.Y || ny >= bounds.Max.Y {
						continue
					}
					r16, g16, b16, _ := img.At(nx, ny).RGBA()
					sumR += int(r16 >> 8)
					sumG += int(g16 >> 8)
					sumB += int(b16 >> 8)
					count++
				}
			}
			if count > 0 {
				avgR := uint8(sumR / count)
				avgG := uint8(sumG / count)
				avgB := uint8(sumB / count)
				_, _, _, a16 := img.At(x, y).RGBA()
				a8 := uint8(a16 >> 8)

				newImg.Set(x, y, color.RGBA{
					R: avgR,
					G: avgG,
					B: avgB,
					A: a8,
				})
			}
		}
	}
	return newImg
}

func ApplyMedianFilter(img image.Image) image.Image {
	bounds := img.Bounds()
	result := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var rVals, gVals, bVals []int
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					nx, ny := x+dx, y+dy
					if nx >= bounds.Min.X && nx < bounds.Max.X && ny >= bounds.Min.Y && ny < bounds.Max.Y {
						r, g, b, _ := img.At(nx, ny).RGBA()
						rVals = append(rVals, int(r>>8))
						gVals = append(gVals, int(g>>8))
						bVals = append(bVals, int(b>>8))
					}
				}
			}
			sort.Ints(rVals)
			sort.Ints(gVals)
			sort.Ints(bVals)
			medianIdx := len(rVals) / 2
			result.Set(x, y, color.RGBA{
				uint8(rVals[medianIdx]),
				uint8(gVals[medianIdx]),
				uint8(bVals[medianIdx]),
				255,
			})
		}
	}

	return result
}

func ApplySobelFilter(img image.Image) image.Image {
	bounds := img.Bounds()
	result := image.NewRGBA(bounds)
	sobelX := [][]int{
		{-1, 0, 1},
		{-2, 0, 2},
		{-1, 0, 1},
	}
	sobelY := [][]int{
		{-1, -2, -1},
		{0, 0, 0},
		{1, 2, 1},
	}
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var sumX, sumY float64
			for i := -1; i <= 1; i++ {
				for j := -1; j <= 1; j++ {
					nx, ny := x+j, y+i
					if nx < bounds.Min.X {
						nx = bounds.Min.X
					}
					if nx >= bounds.Max.X {
						nx = bounds.Max.X - 1
					}
					if ny < bounds.Min.Y {
						ny = bounds.Min.Y
					}
					if ny >= bounds.Max.Y {
						ny = bounds.Max.Y - 1
					}

					pixel := img.At(nx, ny)
					r, _, _, _ := pixel.RGBA()
					intensity := float64(r >> 8)

					sumX += intensity * float64(sobelX[i+1][j+1])
					sumY += intensity * float64(sobelY[i+1][j+1])
				}
			}
			magnitude := math.Sqrt(sumX*sumX + sumY*sumY)
			normalizedVal := uint8(math.Min(255, magnitude))
			result.Set(x, y, color.RGBA{normalizedVal, normalizedVal, normalizedVal, 255})
		}
	}
	return result
}

func ApplyGaussianBlur(img image.Image) image.Image {
	clamp := func(val, min, max int) int {
		if val < min {
			return min
		}
		if val > max {
			return max
		}
		return val
	}
	bounds := img.Bounds()
	result := image.NewRGBA(bounds)
	kernel := [][]int{
		{1, 2, 1},
		{2, 4, 2},
		{1, 2, 1},
	}
	kernelSum := 16
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var sumR, sumG, sumB int
			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					nx := clamp(x+kx, bounds.Min.X, bounds.Max.X-1)
					ny := clamp(y+ky, bounds.Min.Y, bounds.Max.Y-1)

					r, g, b, _ := img.At(nx, ny).RGBA()
					weight := kernel[ky+1][kx+1]

					sumR += int(r>>8) * weight
					sumG += int(g>>8) * weight
					sumB += int(b>>8) * weight
				}
			}
			blurredR := uint8(sumR / kernelSum)
			blurredG := uint8(sumG / kernelSum)
			blurredB := uint8(sumB / kernelSum)
			result.Set(x, y, color.RGBA{blurredR, blurredG, blurredB, 255})
		}
	}
	return result
}
