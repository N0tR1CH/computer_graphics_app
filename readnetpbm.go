package main

import (
	"bufio"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

type ParsingState int

// Below code is a mix of ai and myself as binary work (bit manipulation etc.) was too much for my skillzzz... yet!!!

var currentImgId int

const (
	MagicNumReading ParsingState = iota
	ParamsReading
	PixelsReading
)

// P1 parsing
func parsePbmAscii(r io.Reader) (image.Image, []string, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	var (
		comments      []string
		magicNum      string
		width, height int
		img           *image.RGBA
		x, y          int
	)

	currentState := MagicNumReading

	for scanner.Scan() {
		line := scanner.Text()

		if idx := strings.Index(line, "#"); idx != -1 {
			comments = append(comments, line[idx:])
			line = line[:idx]
		}

		fields := strings.Fields(line)
		for _, field := range fields {
			switch currentState {
			case MagicNumReading:
				magicNum = field
				if magicNum != "P1" {
					return nil, comments, fmt.Errorf(
						"invalid magic number: expected 'P1', got '%s'",
						magicNum,
					)
				}
				currentState = ParamsReading
			case ParamsReading:
				if width == 0 {
					num, err := strconv.Atoi(field)
					if err != nil {
						return nil, comments, fmt.Errorf("invalid width '%s': %v", field, err)
					}
					if num <= 0 {
						return nil, comments, fmt.Errorf("width must be greater than 0, got %d", num)
					}
					width = num
					continue
				}
				if height == 0 {
					num, err := strconv.Atoi(field)
					if err != nil {
						return nil, comments, fmt.Errorf("invalid height '%s': %v", field, err)
					}
					if num <= 0 {
						return nil, comments, fmt.Errorf("height must be greater than 0, got %d", num)
					}
					height = num
					img = image.NewRGBA(image.Rect(0, 0, width, height))
					currentState = PixelsReading
					continue
				}
			case PixelsReading:
				for _, char := range field {
					if char == ' ' || char == '\t' || char == '\n' {
						continue
					}
					value := char - '0'
					if value == 0 {
						img.Set(x, y, color.White)
					} else if value == 1 {
						img.Set(x, y, color.Black)
					} else {
						return nil, comments, fmt.Errorf("pixel value must be 0 or 1, got '%c'", char)
					}
					x++
					if x >= width {
						x = 0
						y++
						if y >= height {
							break
						}
					}
				}
			}
		}
		if currentState == PixelsReading && y >= height {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, comments, fmt.Errorf("error reading PBM data: %v", err)
	}

	if y != height || (y == height && x != 0) {
		return nil, comments, errors.New("incomplete pixel data")
	}

	switch r := r.(type) {
	case *os.File:
		err := r.Close()
		if err != nil {
			return nil, nil, fmt.Errorf("problem with closing file: %v", err)
		}
	}

	return img, comments, nil
}

// P4 Parsing
func parsePbmBinary(r io.Reader) (image.Image, []string, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	var (
		comments            []string
		magicNum            string
		width, height, x, y int
		img                 *image.Gray
	)

	currentState := MagicNumReading

	for scanner.Scan() {
		line := scanner.Text()

		if idx := strings.Index(line, "#"); idx != -1 {
			comments = append(comments, line[idx:])
			line = line[:idx]
		}

		fields := strings.Fields(line)
		for _, field := range fields {
			switch currentState {
			case MagicNumReading:
				magicNum = field
				if magicNum != "P4" {
					return nil, comments, fmt.Errorf("invalid magic number: expected 'P4', got '%s'", magicNum)
				}
				currentState = ParamsReading
			case ParamsReading:
				if width == 0 {
					num, err := strconv.Atoi(field)
					if err != nil {
						return nil, comments, fmt.Errorf("invalid width '%s': %v", field, err)
					}
					if num <= 0 {
						return nil, comments, fmt.Errorf("width must be greater than 0, got %d", num)
					}
					width = num
					continue
				}
				if height == 0 {
					num, err := strconv.Atoi(field)
					if err != nil {
						return nil, comments, fmt.Errorf("invalid height '%s': %v", field, err)
					}
					if num <= 0 {
						return nil, comments, fmt.Errorf("height must be greater than 0, got %d", num)
					}
					height = num
					img = image.NewGray(image.Rect(0, 0, width, height))
					currentState = PixelsReading
					continue
				}
			}
		}

		if currentState == PixelsReading {
			break
		}
	}

	if currentState != PixelsReading {
		return nil, comments, errors.New("missing width/height before pixel data")
	}

	bytesNum := int(math.Ceil(float64(width)/8.0) * float64(height))
	fmt.Println(bytesNum)
	if scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		for _, field := range fields {
			num, err := strconv.Atoi(field)
			if err != nil {
				return nil, nil, fmt.Errorf("problem with string conversion to integer: %+v", err)
			}
			ui8num := uint8(num)
			for bit := 7; bit >= 0; bit-- {
				if x >= width {
					x = 0
					y++
					if y >= height {
						break
					}
					break
				}
				pixel := (ui8num >> bit) & 1 // extract the bit as a pixel (0 or 1)
				switch pixel {
				case 0:
					img.SetGray(x, y, color.Gray{Y: 255})
				case 1:
					img.SetGray(x, y, color.Gray{Y: 0})
				}
				x++
			}
			if y >= height {
				break
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, comments, fmt.Errorf("error reading PBM data: %v", err)
	}

	switch r := r.(type) {
	case *os.File:
		err := r.Close()
		if err != nil {
			return nil, nil, fmt.Errorf("problem with closing file: %v", err)
		}
	}

	return img, comments, nil
}

// P2 Parsing
func parsePgmAscii(r io.Reader) (image.Image, []string, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	var (
		comments              []string
		magicNum              string
		width, height, maxNum int
		scalingFactor         float64
		img                   *image.Gray
		x                     int
		y                     int
	)

	currentState := MagicNumReading
	scalingFactor = 1.0

	for scanner.Scan() {
		line := scanner.Text()

		if idx := strings.Index(line, "#"); idx != -1 {
			comments = append(comments, line[idx:])
			line = line[:idx]
		}

		fields := strings.Fields(line)
		for _, field := range fields {
			switch currentState {
			case MagicNumReading:
				magicNum = field
				if magicNum != "P2" {
					return nil, comments, fmt.Errorf(
						"invalid magic number: expected 'P2', got '%s'",
						magicNum,
					)
				}
				currentState = ParamsReading
			case ParamsReading:
				if width == 0 {
					num, err := strconv.Atoi(field)
					if err != nil {
						return nil, comments, fmt.Errorf("invalid width '%s': %v", field, err)
					}
					if num <= 0 {
						return nil, comments, fmt.Errorf("width must be greater than 0, got %d", num)
					}
					width = num
					continue
				}
				if height == 0 {
					num, err := strconv.Atoi(field)
					if err != nil {
						return nil, comments, fmt.Errorf("invalid height '%s': %v", field, err)
					}
					if num <= 0 {
						return nil, comments, fmt.Errorf("height must be greater than 0, got %d", num)
					}
					height = num
					continue
				}
				if maxNum == 0 {
					num, err := strconv.Atoi(field)
					if err != nil {
						return nil, comments, fmt.Errorf("invalid max color value '%s': %v", field, err)
					}
					if num <= 0 {
						return nil, comments, fmt.Errorf("max color value must be greater than 0, got %d", num)
					}
					maxNum = num
					if maxNum != 255 {
						scalingFactor = 255.0 / float64(maxNum)
						fmt.Printf("Warning: max color value is %d, scaling will be applied.\n", maxNum)
					}
					img = image.NewGray(image.Rect(0, 0, width, height))
					currentState = PixelsReading
					continue
				}
			case PixelsReading:
				num, err := strconv.Atoi(field)
				if err != nil {
					return nil, comments, fmt.Errorf("invalid pixel value '%s': %v", field, err)
				}
				if num < 0 || num > maxNum {
					return nil, comments, fmt.Errorf("pixel value %d out of range (0-%d)", num, maxNum)
				}
				grayness := uint8(float64(num)*scalingFactor + 0.5)
				img.SetGray(x, y, color.Gray{Y: grayness})
				x++
				if x == width {
					x = 0
					y++
					if y >= height {
						break
					}
				}
			}
		}

		if currentState == PixelsReading && y == height {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, comments, fmt.Errorf("error reading PPM data: %v", err)
	}

	if y != height || (y == height && x != 0) {
		return nil, comments, errors.New("incomplete pixel data")
	}

	switch r := r.(type) {
	case *os.File:
		err := r.Close()
		if err != nil {
			return nil, nil, fmt.Errorf("problem with closing file: %v", err)
		}
	}

	return img, comments, nil
}

// P5 Parsing
func parsePgmBinary(r io.Reader) (image.Image, []string, error) {
	bufReader := bufio.NewReader(r)
	var comments []string

	magic, err := bufReader.ReadString('\n')
	if err != nil {
		return nil, nil, err
	}
	magic = strings.TrimSpace(magic)
	if magic != "P5" {
		return nil, nil, fmt.Errorf("invalid magic number: expected P5, got %s", magic)
	}

	var width, height, maxVal int
	for {
		line, err := bufReader.ReadString('\n')
		if err != nil {
			return nil, nil, err
		}
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			comments = append(comments, strings.TrimPrefix(line, "# "))
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 3 {
			_, err := fmt.Sscanf(parts[0], "%d", &width)
			if err != nil {
				return nil, nil, err
			}
			_, err = fmt.Sscanf(parts[1], "%d", &height)
			if err != nil {
				return nil, nil, err
			}
			_, err = fmt.Sscanf(parts[2], "%d", &maxVal)
			if err != nil {
				return nil, nil, err
			}
			break
		}
	}

	if maxVal != 255 {
		return nil, nil, fmt.Errorf("unsupported max value: %d", maxVal)
	}

	img := image.NewGray(image.Rect(0, 0, width, height))
	pixelData := make([]byte, width*height)
	_, err = io.ReadFull(bufReader, pixelData)
	if err != nil {
		return nil, nil, err
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			idx := y*width + x
			grayVal := pixelData[idx]
			img.SetGray(x, y, color.Gray{Y: grayVal})
		}
	}

	return img, comments, nil
}

// P3 Parsing
func parsePpmAscii(r io.Reader) (image.Image, []string, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	var (
		comments              []string
		magicNum              string
		width, height, maxNum int
		scalingFactor         float64
		img                   *image.RGBA
		x                     int
		y                     int
		rgb                   [3]int
		colorIdx              int
	)

	currentState := MagicNumReading
	scalingFactor = 1.0

	for scanner.Scan() {
		line := scanner.Text()

		if idx := strings.Index(line, "#"); idx != -1 {
			comments = append(comments, line[idx:])
			line = line[:idx]
		}

		fields := strings.Fields(line)
		for _, field := range fields {
			switch currentState {
			case MagicNumReading:
				magicNum = field
				if magicNum != "P3" {
					return nil, comments, fmt.Errorf(
						"invalid magic number: expected 'P3', got '%s'",
						magicNum,
					)
				}
				currentState = ParamsReading
			case ParamsReading:
				if width == 0 {
					num, err := strconv.Atoi(field)
					if err != nil {
						return nil, comments, fmt.Errorf("invalid width '%s': %v", field, err)
					}
					if num <= 0 {
						return nil, comments, fmt.Errorf("width must be greater than 0, got %d", num)
					}
					width = num
					continue
				}
				if height == 0 {
					num, err := strconv.Atoi(field)
					if err != nil {
						return nil, comments, fmt.Errorf("invalid height '%s': %v", field, err)
					}
					if num <= 0 {
						return nil, comments, fmt.Errorf("height must be greater than 0, got %d", num)
					}
					height = num
					continue
				}
				if maxNum == 0 {
					num, err := strconv.Atoi(field)
					if err != nil {
						return nil, comments, fmt.Errorf("invalid max color value '%s': %v", field, err)
					}
					if num <= 0 {
						return nil, comments, fmt.Errorf("max color value must be greater than 0, got %d", num)
					}
					maxNum = num
					if maxNum != 255 {
						scalingFactor = 255.0 / float64(maxNum)
						fmt.Printf("Warning: max color value is %d, scaling will be applied.\n", maxNum)
					}
					img = image.NewRGBA(image.Rect(0, 0, width, height))
					currentState = PixelsReading
					continue
				}
			case PixelsReading:
				num, err := strconv.Atoi(field)
				if err != nil {
					return nil, comments, fmt.Errorf("invalid pixel value '%s': %v", field, err)
				}
				if num < 0 || num > maxNum {
					return nil, comments, fmt.Errorf("pixel value %d out of range (0-%d)", num, maxNum)
				}
				rgb[colorIdx] = num
				colorIdx++
				if colorIdx == 3 {
					r := uint8(float64(rgb[0])*scalingFactor + 0.5)
					g := uint8(float64(rgb[1])*scalingFactor + 0.5)
					b := uint8(float64(rgb[2])*scalingFactor + 0.5)
					a := uint8(255)
					img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: a})
					colorIdx = 0
					x++
					if x == width {
						x = 0
						y++
					}
				}
			}
		}

		if currentState == PixelsReading && y == height {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, comments, fmt.Errorf("error reading PPM data: %v", err)
	}

	if y != height || (y == height && x != 0) {
		return nil, comments, errors.New("incomplete pixel data")
	}

	switch r := r.(type) {
	case *os.File:
		err := r.Close()
		if err != nil {
			return nil, nil, fmt.Errorf("problem with closing file: %v", err)
		}
	}

	return img, comments, nil
}

// P6 parsing
func parsePpmBinary(r io.Reader) (image.Image, []string, error) {
	bufReader := bufio.NewReader(r)
	var comments []string

	// Read magic number
	magic, err := bufReader.ReadString('\n')
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read magic number: %v", err)
	}
	magic = strings.TrimSpace(magic)
	if magic != "P6" {
		return nil, nil, fmt.Errorf("invalid magic number: expected P6, got %s", magic)
	}

	// Initialize variables to store header information
	var width, height, maxVal int
	tokensCollected := 0

	for tokensCollected < 3 {
		line, err := bufReader.ReadString('\n')
		if err != nil {
			return nil, nil, fmt.Errorf("failed to read header line: %v", err)
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		if strings.HasPrefix(line, "#") {
			// Entire line is a comment
			comments = append(comments, strings.TrimSpace(strings.TrimPrefix(line, "#")))
			continue
		}

		// Handle inline comments by splitting at '#'
		if idx := strings.Index(line, "#"); idx != -1 {
			line = line[:idx]
		}

		parts := strings.Fields(line)
		for _, part := range parts {
			switch tokensCollected {
			case 0:
				width, err = strconv.Atoi(part)
				if err != nil {
					return nil, nil, fmt.Errorf("invalid width '%s': %v", part, err)
				}
			case 1:
				height, err = strconv.Atoi(part)
				if err != nil {
					return nil, nil, fmt.Errorf("invalid height '%s': %v", part, err)
				}
			case 2:
				maxVal, err = strconv.Atoi(part)
				if err != nil {
					return nil, nil, fmt.Errorf("invalid maxVal '%s': %v", part, err)
				}
			}
			tokensCollected++
			if tokensCollected == 3 {
				break
			}
		}
	}

	if maxVal != 255 {
		return nil, nil, fmt.Errorf("unsupported max value: %d", maxVal)
	}

	// Read binary data
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	pixelData := make([]byte, width*height*3)
	_, err = io.ReadFull(bufReader, pixelData)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read pixel data: %v", err)
	}

	// Populate the image with pixel data
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			idx := (y*width + x) * 3
			if idx+2 >= len(pixelData) {
				return nil, nil, fmt.Errorf("unexpected end of pixel data")
			}
			r := pixelData[idx]
			g := pixelData[idx+1]
			b := pixelData[idx+2]
			img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}

	return img, comments, nil
}

func parseNetPbm(file *os.File) (image.Image, []string, error) {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var magicNum string

	if scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		magicNum = fields[0]
	}
	if _, err := file.Seek(0, 0); err != nil {
		return nil, nil, fmt.Errorf("problem setting offset back: +%v", err)
	}

	switch magicNum {
	case "P1":
		img, comments, err := parsePbmAscii(file)
		if err != nil {
			return nil, nil, err
		}
		return img, comments, nil
	case "P2":
		img, comments, err := parsePgmAscii(file)
		if err != nil {
			return nil, nil, err
		}
		return img, comments, nil
	case "P3":
		img, comments, err := parsePpmAscii(file)
		if err != nil {
			return nil, nil, err
		}
		return img, comments, nil
	case "P4":
		img, comments, err := parsePbmBinary(file)
		if err != nil {
			return nil, nil, err
		}
		return img, comments, nil
	case "P5":
		img, comments, err := parsePgmBinary(file)
		if err != nil {
			return nil, nil, err
		}
		return img, comments, nil
	case "P6":
		img, comments, err := parsePpmBinary(file)
		if err != nil {
			return nil, nil, err
		}
		return img, comments, nil
	default:
		return nil, nil, errors.New("no format")
	}
}

// This was made in another program as it was to hard to deal with wails and also
// iterate over such hard thing
// keeping this in the same package as this workaround with static server hosting files
// is really painful but i know that in the long run it seems to be the better choice
// func main() {
// 	file, err := os.Open("ppm-test-06-p6.ppm")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	img, comments, err := parseNetPbm(file)
// 	fmt.Println(comments)
//
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	outFile, err := os.Create(fmt.Sprintf("file%d.png", currentImgId))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	if err := png.Encode(outFile, img); err != nil {
// 		panic(err)
// 	}
// 	if err := outFile.Close(); err != nil {
// 		panic(err)
// 	}
// }
