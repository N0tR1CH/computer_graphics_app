package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"strconv"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type operation uint8

const (
	addition operation = iota
	substraction
	multiplication
	division
)

type pointWiseRgbValues struct {
	r  uint8
	g  uint8
	b  uint8
	op operation
}

func (a *App) HandleAlphaPointWiseTransformations(alphaVal uint8, base64str string) string {
	var newBase64str string
	data, err := dataFromBase64(a.ctx, base64str)
	if err != nil {
		return newBase64str
	}
	imgBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:          runtime.InfoDialog,
			Title:         "Could not proceed with the operation",
			Message:       fmt.Sprintf("Couldn't decode the image sry: %s", err.Error()),
			DefaultButton: "Ok",
		})
		return newBase64str
	}
	img, err := png.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:          runtime.InfoDialog,
			Title:         "Could not proceed with the operation",
			Message:       fmt.Sprintf("Couldn't decode the image sry: %s", err.Error()),
			DefaultButton: "Ok",
		})
	}
	newBase64str = generateNewAlphaBase64Str(alphaVal, img)
	if newBase64str == "" {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:          runtime.InfoDialog,
			Title:         "Could not proceed with the operation",
			Message:       "Who else know what",
			DefaultButton: "Ok",
		})
	}
	return newBase64str
}

func generateNewAlphaBase64Str(a8 uint8, m image.Image) string {
	if m == nil {
		return ""
	}
	// From stdlib example, trying this out because positions might not start from 0
	// although in example they do so idk/idc
	bounds := m.Bounds()
	newM := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r16, g16, b16, _ := m.At(x, y).RGBA()
			r8 := uint8(r16 >> 8)
			g8 := uint8(g16 >> 8)
			b8 := uint8(b16 >> 8)
			newM.Set(x, y, color.RGBA{r8, g8, b8, a8})
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, newM); err != nil {
		return ""
	}
	base64str := base64.StdEncoding.EncodeToString(buf.Bytes())
	dataUrl := fmt.Sprintf("data:image/png;base64,%s", base64str)
	return dataUrl
}

func (a *App) HandleRgbPointWiseTransformations(values []string, base64str string) string {
	var newBase64str string
	pwrv, err := parseRgb(values)
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:          runtime.InfoDialog,
			Title:         "Could not proceed with the operation",
			Message:       fmt.Sprintf("problem with values from the form: %s", err.Error()),
			DefaultButton: "Ok",
		})
		return newBase64str
	}

	data, err := dataFromBase64(a.ctx, base64str)
	if err != nil {
		return newBase64str
	}
	imgBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:          runtime.InfoDialog,
			Title:         "Could not proceed with the operation",
			Message:       fmt.Sprintf("Couldn't decode the image sry: %s", err.Error()),
			DefaultButton: "Ok",
		})
		return newBase64str
	}
	img, err := png.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:          runtime.InfoDialog,
			Title:         "Could not proceed with the operation",
			Message:       fmt.Sprintf("Couldn't decode the image sry: %s", err.Error()),
			DefaultButton: "Ok",
		})
	}
	newBase64str = generateNewRgbBase64Str(*pwrv, img)
	if newBase64str == "" {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:          runtime.InfoDialog,
			Title:         "Could not proceed with the operation",
			Message:       "Who else know what",
			DefaultButton: "Ok",
		})
	}
	return newBase64str
}

func generateNewRgbBase64Str(pwrv pointWiseRgbValues, img image.Image) string {
	if img == nil {
		return ""
	}
	newImg := image.NewRGBA(img.Bounds())
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			clr := img.At(x, y)
			r16, g16, b16, a16 := clr.RGBA()
			r8 := uint8(r16 >> 8)
			g8 := uint8(g16 >> 8)
			b8 := uint8(b16 >> 8)
			a8 := uint8(a16 >> 8)
			switch pwrv.op {
			case addition:
				r8 = uint8(math.Min(float64(r8+pwrv.r), 255))
				g8 = uint8(math.Min(float64(g8+pwrv.r), 255))
				b8 = uint8(math.Min(float64(b8+pwrv.r), 255))
			case substraction:
				r8 = uint8(math.Max(float64(r8-pwrv.r), 0))
				g8 = uint8(math.Max(float64(g8-pwrv.r), 0))
				b8 = uint8(math.Max(float64(b8-pwrv.r), 0))
			case division:
				r8 = uint8(math.Max(float64(r8/pwrv.r), 0))
				g8 = uint8(math.Max(float64(g8/pwrv.r), 0))
				b8 = uint8(math.Max(float64(b8/pwrv.r), 0))
			case multiplication:
				r8 = uint8(math.Min(float64(r8*pwrv.r), 255))
				g8 = uint8(math.Min(float64(g8*pwrv.r), 255))
				b8 = uint8(math.Min(float64(b8*pwrv.r), 255))
			}
			newImg.Set(x, y, color.RGBA{r8, g8, b8, a8})
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, newImg); err != nil {
		return ""
	}
	base64str := base64.StdEncoding.EncodeToString(buf.Bytes())
	dataUrl := fmt.Sprintf("data:image/png;base64,%s", base64str)
	return dataUrl
}

func parseRgb(values []string) (*pointWiseRgbValues, error) {
	if len(values) != 4 {
		return nil, errors.New("Operation require 4 elements")
	}
	valuesStruct := new(pointWiseRgbValues)
	switch values[0] {
	case "addition":
		valuesStruct.op = addition
	case "substraction":
		valuesStruct.op = substraction
	case "multiplication":
		valuesStruct.op = multiplication
	case "division":
		valuesStruct.op = division
	default:
		return nil, errors.New("Unknown operation for rgb")
	}

	r, err := strconv.Atoi(values[1])
	if err != nil {
		return nil, errors.New("Problem with converting integer")
	}
	if iLimitExceed(r) {
		return nil, errors.New("r must be value between 0 and 255")
	}
	valuesStruct.r = uint8(r)

	g, err := strconv.Atoi(values[2])
	if err != nil {
		return nil, errors.New("Problem with converting integer")
	}
	if iLimitExceed(g) {
		return nil, errors.New("g must be value between 0 and 255")
	}
	valuesStruct.g = uint8(g)

	b, err := strconv.Atoi(values[3])
	if err != nil {
		return nil, errors.New("Problem with converting integer")
	}
	if iLimitExceed(b) {
		return nil, errors.New("b must be value between 0 and 255")
	}
	valuesStruct.b = uint8(b)

	return valuesStruct, nil
}

func iLimitExceed(v int) bool {
	if v < 0 || v > 255 {
		return true
	}
	return false
}
