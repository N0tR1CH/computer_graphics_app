package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
)

func grayImageFromImage(m image.Image) *image.Gray {
	b := m.Bounds()
	grayM := image.NewGray(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			oldPx := m.At(x, y)
			px := color.GrayModel.Convert(oldPx)
			grayM.Set(x, y, px)
		}
	}
	return grayM
}

func (a *App) HandleDilation(base64img string) string {
	m, err := decodeBasePngToImg(base64img, a.ctx)
	if err != nil {
		return ""
	}
	newM := dilation(m)
	var buf bytes.Buffer
	if err := png.Encode(&buf, newM); err != nil {
		return ""
	}
	base64img = base64.StdEncoding.EncodeToString(buf.Bytes())
	return fmt.Sprintf("data:image/png;base64,%s", base64img)
}

func dilation(m image.Image) image.Image {
	se := [][]bool{
		{false, true, false},
		{true, true, true},
		{false, true, false},
	}
	b := m.Bounds()
	grayM := grayImageFromImage(m)
	dilatedM := image.NewGray(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			var maxVal uint8 = 0
			for seY := 0; seY < 3; seY++ {
				for seX := 0; seX < 3; seX++ {
					if se[seY][seX] {
						nx := x + seX - 1
						ny := y + seY - 1
						if nx >= b.Min.X && nx < b.Max.X &&
							ny >= b.Min.Y && ny < b.Max.Y {
							val := grayM.GrayAt(nx, ny).Y
							if val > maxVal {
								maxVal = val
							}
						}

					}
				}
			}
			dilatedM.SetGray(x, y, color.Gray{Y: maxVal})
		}
	}
	return dilatedM
}

func (a *App) HandleErosion(base64img string) string {
	m, err := decodeBasePngToImg(base64img, a.ctx)
	if err != nil {
		return ""
	}
	newM := erosion(m)
	var buf bytes.Buffer
	if err := png.Encode(&buf, newM); err != nil {
		return ""
	}
	base64img = base64.StdEncoding.EncodeToString(buf.Bytes())
	return fmt.Sprintf("data:image/png;base64,%s", base64img)
}

func erosion(m image.Image) image.Image {
	se := [][]bool{
		{false, true, false},
		{true, true, true},
		{false, true, false},
	}
	grayImg := grayImageFromImage(m)
	b := grayImg.Bounds()
	erodedM := image.NewGray(b)
	seHeight := len(se)
	seWidth := len(se[0])
	halfHeight := seHeight / 2
	halfWidth := seWidth / 2
	type offset struct {
		dx, dy int
	}
	var offsets []offset
	for seY, row := range se {
		for seX, val := range row {
			if val {
				offsets = append(
					offsets,
					offset{dx: seX - halfWidth, dy: seY - halfHeight},
				)
			}
		}
	}
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			minVal := uint8(255)
			for _, off := range offsets {
				nx := x + off.dx
				ny := y + off.dy
				if nx >= b.Min.X && nx < b.Max.X &&
					ny >= b.Min.Y && ny < b.Max.Y {
					val := grayImg.GrayAt(nx, ny).Y
					if val < minVal {
						minVal = val
					}
				}
			}
			erodedM.SetGray(x, y, color.Gray{Y: minVal})
		}
	}
	return erodedM
}
