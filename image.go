package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type (
	imageFormat    string
	imageFormatErr string
)

const (
	jpg   = imageFormat("jpeg")
	pbmP1 = imageFormat("pbmP1")
	pbmP4 = imageFormat("pbmP4")
	pgmP2 = imageFormat("pgmP2")
	pgmP5 = imageFormat("pgmP5")
	ppmP3 = imageFormat("ppmP3")
	ppmP6 = imageFormat("ppmP6")

	errImageFormatUnknown = imageFormatErr("incorrect format")
)

func (format imageFormat) validate() error {
	switch format {
	case jpg, pbmP1, pbmP4, pgmP2, pgmP5, ppmP3, ppmP6:
		return nil
	default:
		return errImageFormatUnknown
	}
}

func (e imageFormatErr) Error() string {
	return string(e)
}

func (a *App) SaveCanvasImg(base64Image string, format imageFormat) {
	if err := format.validate(); err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:          runtime.InfoDialog,
			Title:         "Could not proceed with the operation",
			Message:       fmt.Sprintf("'%s' is invalid file format, possible ones are jpeg, pbm, pgm, ppm", format),
			DefaultButton: "Ok",
		})
		return
	}

	parts := strings.Split(base64Image, ",")
	if len(parts) != 2 {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:          runtime.InfoDialog,
			Title:         "Could not proceed with the operation",
			Message:       "Invalid data URI format",
			DefaultButton: "Ok",
		})
		return
	}

	data := parts[1]
	imgBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:          runtime.InfoDialog,
			Title:         "Could not proceed with the operation",
			Message:       "Couldn't decode the image sry",
			DefaultButton: "Ok",
		})
		return
	}

	// reader := bytes.NewReader(imgBytes)
	// img, err := jpeg.Decode(reader)
	// if err != nil {
	// 	runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
	// 		Type:          runtime.InfoDialog,
	// 		Title:         "Decoding problem",
	// 		Message:       "Image could not be decoded",
	// 		DefaultButton: "Ok",
	// 	})
	// 	return
	// }

	filepath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		DefaultFilename: "canvas",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "JPEG Image *.jpg",
				Pattern:     "*.jpg",
			},
		},
	})

	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:          runtime.InfoDialog,
			Title:         "Could not proceed with the operation",
			Message:       "Filepath problem ;(",
			DefaultButton: "Ok",
		})
		return
	}

	if err := os.WriteFile(filepath, imgBytes, 0644); err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:          runtime.InfoDialog,
			Title:         "File was not saved",
			Message:       "Operation was cancelled or there was a problem with file permissions",
			DefaultButton: "Ok",
		})
	}
}
