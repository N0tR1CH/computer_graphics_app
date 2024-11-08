package main

import "image/color"

type Cmyk struct {
	C uint8 `json:"c"`
	M uint8 `json:"m"`
	Y uint8 `json:"y"`
	K uint8 `json:"k"`
}

type Rgb struct {
	R uint8 `json:"r"`
	G uint8 `json:"g"`
	B uint8 `json:"b"`
}

func (a *App) RgbToCmyk(r, g, b uint8) Cmyk {
	c, m, y, k := color.RGBToCMYK(r, g, b)
	return Cmyk{
		C: c,
		M: m,
		Y: y,
		K: k,
	}
}

func (a *App) CmykToRgb(c, m, y, k uint8) Rgb {
	r, g, b := color.CMYKToRGB(c, m, y, k)
	return Rgb{r, g, b}
}
