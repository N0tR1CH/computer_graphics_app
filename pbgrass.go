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

func pixelGreen(r, g, b uint32, threshold uint8) bool {
	r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)
	if int(g8) < int(r8)+int(threshold) || int(g8) < int(b8)+int(threshold) {
		return false
	}
	return true
}

func (a *App) HandleGrassTask(base64img string, threshold uint8) string {
	m, err := decodeBasePngToImg(base64img, a.ctx)
	if err != nil {
		return ""
	}

	m, binary, percent := binarizeOtsuForBfsWithGreenPercentCalculation(m, threshold)
	runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:          runtime.InfoDialog,
		Title:         "Green percentage on the image",
		Message:       fmt.Sprintf("Green percentage on the image is: %f%%", percent),
		DefaultButton: "Ok",
	})
	largestGroup := findLargestGroup(binary)

	coloredM := image.NewRGBA(m.Bounds())
	for y := m.Bounds().Min.Y; y < m.Bounds().Max.Y; y++ {
		for x := m.Bounds().Min.X; x < m.Bounds().Max.X; x++ {
			coloredM.Set(x, y, m.At(x, y))
		}
	}
	for _, p := range largestGroup {
		coloredM.Set(p.X, p.Y, color.RGBA{255, 0, 0, 255})
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, coloredM); err != nil {
		return ""
	}
	base64img = base64.StdEncoding.EncodeToString(buf.Bytes())
	return fmt.Sprintf("data:image/png;base64,%s", base64img)
}

func binarizeOtsuForBfsWithGreenPercentCalculation(m image.Image, greenThreshold uint8) (image.Image, [][]bool, float64) {
	b := m.Bounds()
	histogram := make([]int, 256)
	totalPxCount := b.Dy() * b.Dx()
	greenPxCount := 0
	binary := make([][]bool, b.Dy())
	for y := 0; y < b.Dy(); y++ {
		binary[y] = make([]bool, b.Dx())
		for x := 0; x < b.Dx(); x++ {
			r, g, b, _ := m.At(x, y).RGBA()
			if pixelGreen(r, g, b, greenThreshold) {
				greenPxCount++
			}
			lum := color.GrayModel.Convert(m.At(x, y)).(color.Gray).Y
			histogram[lum]++
			binary[y][x] = false
		}
	}
	total := b.Dx() * b.Dy()
	sum := 0
	for i := 0; i < 256; i++ {
		sum += i * histogram[i]
	}
	sumB, wB, max := 0, 0, 0.0
	threshold := uint8(0)
	for t := 0; t < 256; t++ {
		wB += histogram[t]
		if wB == 0 {
			continue
		}
		wF := total - wB
		if wF == 0 {
			break
		}
		sumB += t * histogram[t]
		mB := float64(sumB) / float64(wB)
		mF := float64(sum-sumB) / float64(wF)
		variance := float64(wB) * float64(wF) * (mB - mF) * (mB - mF)
		if variance > max {
			max = variance
			threshold = uint8(t)
		}
	}
	binaryImg := image.NewGray(b)
	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			lum := color.GrayModel.Convert(m.At(x, y)).(color.Gray).Y
			if lum <= threshold {
				binaryImg.Set(x, y, color.White)
				binary[y][x] = true
			} else {
				binaryImg.Set(x, y, color.Black)
			}
		}
	}
	return binaryImg, binary, (float64(greenPxCount) / float64(totalPxCount)) * 100
}

type Point struct {
	X, Y int
}

func findLargestGroup(binary [][]bool) []Point {
	visited := make([][]bool, len(binary))
	for i := range visited {
		visited[i] = make([]bool, len(binary[0]))
	}
	var largest []Point
	for y := 0; y < len(binary); y++ {
		for x := 0; x < len(binary[0]); x++ {
			if binary[y][x] && !visited[y][x] {
				group := bfs(binary, visited, x, y)
				if len(group) > len(largest) {
					largest = group
				}
			}
		}
	}
	return largest
}

func bfs(binary [][]bool, visited [][]bool, startX, startY int) []Point {
	queue := []Point{{startX, startY}}
	visited[startY][startX] = true
	var group []Point
	directions := []Point{
		{1, 0}, {-1, 0}, {0, 1}, {0, -1},
	}
	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		group = append(group, p)
		for _, d := range directions {
			newX, newY := p.X+d.X, p.Y+d.Y
			if newY >= 0 && newY < len(binary) && newX >= 0 && newX < len(binary[0]) {
				if binary[newY][newX] && !visited[newY][newX] {
					visited[newY][newX] = true
					queue = append(queue, Point{newX, newY})
				}
			}
		}
	}
	return group
}
