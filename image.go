package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"io"
	"os"
	"strings"

	"github.com/spakin/netpbm"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type (
	ImageFormat    string
	ImageFormatErr string
)

const (
	jpg                   ImageFormat    = "jpeg"
	pbmP1                 ImageFormat    = "pbmP1"
	pbmP4                 ImageFormat    = "pbmP4"
	pgmP2                 ImageFormat    = "pgmP2"
	pgmP5                 ImageFormat    = "pgmP5"
	ppmP3                 ImageFormat    = "ppmP3"
	ppmP6                 ImageFormat    = "ppmP6"
	errImageFormatUnknown ImageFormatErr = "incorrect format"
)

var AllImageFormats = []struct {
	Value  ImageFormat
	TSName string
}{
	{jpg, "jpg"},
	{pbmP1, "pbmP1"},
	{pbmP4, "pbmP4"},
	{pgmP2, "pgmP2"},
	{pgmP5, "pgmP5"},
	{ppmP3, "ppmP3"},
	{ppmP6, "ppmP6"},
}

func (format ImageFormat) validate() error {
	switch format {
	case jpg, pbmP1, pbmP4, pgmP2, pgmP5, ppmP3, ppmP6:
		return nil
	default:
		return errImageFormatUnknown
	}
}

func (e ImageFormatErr) Error() string {
	return string(e)
}

func (a *App) SaveCanvasImg(base64Image string, format ImageFormat) {
	if err := format.validate(); err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:          runtime.InfoDialog,
			Title:         "Could not proceed with the operation",
			Message:       fmt.Sprintf("'%s' is invalid file format, possible ones are jpeg, pbm, pgm, ppm", format),
			DefaultButton: "Ok",
		})
		return
	}

	runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:          runtime.InfoDialog,
		Title:         "Could not proceed with the operation",
		Message:       fmt.Sprintf("'%s' is invalid file format, possible ones are jpeg, pbm, pgm, ppm", format),
		DefaultButton: "Ok",
	})

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

	reader := bytes.NewReader(imgBytes)
	img, err := jpeg.Decode(reader)

	var buf bytes.Buffer
	if err := netpbm.Encode(&buf, img, &netpbm.EncodeOptions{
		Format:   netpbm.PBM,
		Plain:    true,
		Comments: []string{"dupa", "sraka"},
	}); err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:          runtime.InfoDialog,
			Title:         "Encoding problem",
			Message:       "Image could not be encoded",
			DefaultButton: "Ok",
		})
	}

	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:          runtime.InfoDialog,
			Title:         "Decoding problem",
			Message:       "Image could not be decoded",
			DefaultButton: "Ok",
		})
		return
	}

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

	bytes, err := io.ReadAll(&buf)
	if err := os.WriteFile("./file.pbm", bytes, 0644); err != nil {

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
