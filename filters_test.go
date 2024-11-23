package main

import (
	"image"
	"image/color"
	"testing"
)

func createTestImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 3, 3))
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			if x == 1 && y == 1 {
				img.Set(x, y, color.RGBA{255, 255, 255, 255})
			} else {
				img.Set(x, y, color.RGBA{0, 0, 0, 255})
			}
		}
	}
	return img
}

func TestApplyAveragingFilter(t *testing.T) {
	testImg := createTestImage()
	filteredImg := ApplyAveragingFilter(testImg)

	tests := []struct {
		x, y    int
		wantR   uint8
		message string
	}{
		{0, 0, 63, "Corner pixel"},
		{0, 1, 42, "Edge pixel"},
		{1, 1, 28, "Center pixel"},
	}

	for _, tt := range tests {
		r, _, _, _ := filteredImg.At(tt.x, tt.y).RGBA()
		r8 := uint8(r >> 8)
		if r8 != tt.wantR {
			t.Errorf("%s at (%d,%d): got %d, want %d",
				tt.message, tt.x, tt.y, r8, tt.wantR)
		}
	}
}

func TestApplyMedianFilter(t *testing.T) {
	testImg := createTestImage()
	filteredImg := ApplyMedianFilter(testImg)
	tests := []struct {
		x, y    int
		wantR   uint8
		message string
	}{
		{0, 0, 0, "Corner pixel"},
		{0, 1, 0, "Edge pixel"},
		{1, 1, 0, "Center pixel"},
	}
	for _, tt := range tests {
		r, _, _, _ := filteredImg.At(tt.x, tt.y).RGBA()
		r8 := uint8(r >> 8)
		if r8 != tt.wantR {
			t.Errorf("%s at (%d,%d): got %d, want %d",
				tt.message, tt.x, tt.y, r8, tt.wantR)
		}
	}
}
