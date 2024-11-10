package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
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

func (format ImageFormat) validate(ctx context.Context) error {
	switch format {
	case jpg, pbmP1, pbmP4, pgmP2, pgmP5, ppmP3, ppmP6:
		return nil
	default:
		runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
			Type:          runtime.InfoDialog,
			Title:         "Could not proceed with the operation",
			Message:       fmt.Sprintf("'%s' is invalid file format, possible ones are jpeg, pbm, pgm, ppm", format),
			DefaultButton: "Ok",
		})
		return errImageFormatUnknown
	}
}

func (format ImageFormat) netpbm() bool {
	switch format {
	case pbmP1, pbmP4, pgmP2, pgmP5, ppmP3, ppmP6:
		return true
	default:
		return false
	}
}

func (format ImageFormat) filters() (displayName string, pattern string) {
	displayName, pattern = "JPEG Image", "*.jpg"
	if format.netpbm() {
		displayName = fmt.Sprintf("%s", strings.ToUpper(string(format[:3])))
		pattern = fmt.Sprintf("*.%s", string(format[:3]))
	}
	return displayName, pattern
}

func (format ImageFormat) plain() bool {
	switch format {
	case pbmP1, pgmP2, ppmP3:
		return true
	default:
		return false
	}
}

func (format ImageFormat) format() netpbm.Format {
	switch format {
	case pbmP1, pbmP4:
		return netpbm.PBM
	case pgmP2, pgmP5:
		return netpbm.PGM
	case ppmP3, ppmP6:
		return netpbm.PPM
	}
	return -1
}

func (e ImageFormatErr) Error() string {
	return string(e)
}

func dataFromBase64(ctx context.Context, base64Image string) (string, error) {
	parts := strings.Split(base64Image, ",")
	if len(parts) != 2 {
		msg := "Invalid data URI format"
		runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
			Type:          runtime.InfoDialog,
			Title:         "Could not proceed with the operation",
			Message:       msg,
			DefaultButton: "Ok",
		})
		return "", errors.New(msg)
	}

	return parts[1], nil
}

func rightImgBytes(
	imgBytes []byte,
	format ImageFormat,
	comments []string,
	ctx context.Context,
) ([]byte, error) {
	if !format.netpbm() {
		return imgBytes, nil
	}

	img, err := jpeg.Decode(bytes.NewReader(imgBytes))
	var buf bytes.Buffer
	if err := netpbm.Encode(&buf, img, &netpbm.EncodeOptions{
		Format:   format.format(),
		Plain:    format.plain(),
		Comments: comments,
	}); err != nil {
		runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
			Type:          runtime.InfoDialog,
			Title:         "Encoding problem",
			Message:       "Image could not be encoded",
			DefaultButton: "Ok",
		})
	}

	imgBytes, err = io.ReadAll(&buf)
	if err != nil {
		runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
			Type:          runtime.InfoDialog,
			Title:         "Decoding problem",
			Message:       "Image could not be decoded",
			DefaultButton: "Ok",
		})
		return nil, errors.New("Image could not be decoded")
	}
	return imgBytes, nil
}

func (a *App) SaveCanvasImg(
	base64Image string,
	format ImageFormat,
	comments []string,
) {
	fmt.Println(comments)
	if err := format.validate(a.ctx); err != nil {
		return
	}

	data, err := dataFromBase64(a.ctx, base64Image)
	if err != nil {
		return
	}

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

	displayName, pattern := format.filters()
	filepath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		DefaultFilename: "canvas",
		Filters: []runtime.FileFilter{
			{
				DisplayName: displayName,
				Pattern:     pattern,
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

	imgBytes, err = rightImgBytes(imgBytes, format, comments, a.ctx)
	if err != nil {
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
