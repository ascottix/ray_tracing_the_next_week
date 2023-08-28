package main

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
)

type Texture interface {
	Value(u, v float64, p Point3) Color
}

// Single color
type SolidColorTexture struct {
	value Color
}

// Checkered pattern using two other textures as the "white" and "black" squares
type CheckerTexture struct {
	scale float64
	even  Texture
	odd   Texture
}

// Image-based texture
type ImageTexture struct {
	data   []Color
	width  int
	height int
}

// Random-block texture used for image 8
type RandomBlockTexture struct {
	scale float64
	size  int
	data  []float64
}

// Noise texture
type NoiseTexture struct {
	scale float64
	noise NoiseGenerator
}

func NewSolidColorTexture(c Color) SolidColorTexture {
	return SolidColorTexture{value: c}
}

func (sct SolidColorTexture) Value(u, v float64, p Point3) Color {
	return sct.value
}

func NewCheckerTexture(scale float64, even, odd Texture) CheckerTexture {
	return CheckerTexture{scale, even, odd}
}

// Utility to build a checkered texture from two colors
func NewBicolorCheckerTexture(scale float64, even, odd Color) CheckerTexture {
	return CheckerTexture{scale: scale, even: NewSolidColorTexture(even), odd: NewSolidColorTexture(odd)}
}

func (ct CheckerTexture) Value(u, v float64, p Point3) Color {
	x := math.Floor(ct.scale * p.X)
	y := math.Floor(ct.scale * p.Y)
	z := math.Floor(ct.scale * p.Z)

	even := int(x+y+z)%2 == 0

	if even {
		return ct.even.Value(u, v, p)
	} else {
		return ct.odd.Value(u, v, p)
	}
}

// Creates an image texture from a PNG or JPEG file
func NewImageTexture(filename string) ImageTexture {
	f, err := os.Open(filename)

	if err == nil {
		defer f.Close()

		image, _, err := image.Decode(f)

		if err == nil {
			bounds := image.Bounds()

			w := bounds.Max.X - bounds.Min.X
			h := bounds.Max.Y - bounds.Min.Y

			t := ImageTexture{data: make([]Color, w*h), width: w, height: h}

			// Convert the image pixels to our internal color format
			for y := 0; y < h; y++ {
				for x := 0; x < w; x++ {
					r, g, b, _ := image.At(bounds.Min.X+x, bounds.Min.Y+y).RGBA()

					o := x + y*w // Offset into our data array

					t.data[o].X = float64(r) / 0xffff
					t.data[o].Y = float64(g) / 0xffff
					t.data[o].Z = float64(b) / 0xffff
				}
			}

			return t
		} else {
			panic("Cannot decode image: " + filename)
		}
	} else {
		panic("Cannot load file: " + filename)
	}
}

func (it ImageTexture) Value(u, v float64, p Point3) Color {
	clamp := func(value float64) float64 {
		return math.Max(0, math.Min(1, value))
	}

	u, v = clamp(u), clamp(v)
	i, j := int(u*float64(it.width-1)), int((1-v)*float64(it.height-1))
	o := i + j*it.width

	return it.data[o]
}

func NewRandomBlockTexture(scale float64) RandomBlockTexture {
	const size = 4 // Must be a power of 2
	data := make([]float64, size*size*size)
	for i := range data {
		data[i] = RandomDouble()
	}
	return RandomBlockTexture{scale, size, data}
}

func (rbt RandomBlockTexture) Value(u, v float64, p Point3) Color {
	f := func(coord float64) int {
		return int(rbt.scale*coord) & (rbt.size - 1)
	}

	x, y, z := f(p.X), f(p.Y), f(p.Z)
	i := z*rbt.size*rbt.size + y*rbt.size + x
	b := rbt.data[i]

	return Color{b, b, b}
}

func NewNoiseTextureWith(mode int) NoiseTexture {
	return NoiseTexture{scale: 1, noise: NewPerlin(mode)}
}

func NewNoiseTexture(scale float64) NoiseTexture {
	return NoiseTexture{scale: scale, noise: NewPerlin(NoiseTrilinearInterpolationWithHermitianSmoothing)}
}

func NewNoiseTextureWithGenerator(scale float64, noise NoiseGenerator) NoiseTexture {
	return NoiseTexture{scale, noise}
}

func (nt NoiseTexture) Value(u, v float64, p Point3) Color {
	return Color{1, 1, 1}.Mul(nt.noise.Noise(p.Mul(nt.scale)))
}

func NewMarbleTexture(scale float64) NoiseTexture {
	// To get the marble effect right, turbulence should use the unscaled point, that's why the 1/scale factor
	return NewNoiseTextureWithGenerator(scale, NewTurbulenceNoiseWithPhase(10, 1/scale, 7))
}
