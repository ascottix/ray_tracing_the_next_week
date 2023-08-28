package main

import (
	"fmt"
	"io"
)

func blur(noise []float64, width, height int) []float64 {
	buf := make([]float64, len(noise))

	at := func(x, y int) float64 {
		x = (x + width) % width
		y = (y + height) % height

		return noise[y*width+x]
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := at(x-1, y-1) + at(x-1, y) + at(x-1, y+1) +
				at(x, y-1) + at(x, y) + at(x, y+1) +
				at(x+1, y-1) + at(x+1, y) + at(x+1, y+1)

			buf[y*width+x] = c / 9
		}
	}

	return buf
}

func Image7(w io.Writer) {
	width, height := 400, 225

	fmt.Fprintf(w, "P3\n") // Magic
	fmt.Fprintf(w, "%d %d\n", width, height)
	fmt.Fprintf(w, "255\n") // Maximum value of a color component

	// Generate the noise in a buffer, so we can process it later
	noise := make([]float64, width*height)
	for i := range noise {
		noise[i] = RandomDouble()
	}

	// Apply a simple box filter a few times
	for w := 0; w < 4; w++ {
		noise = blur(noise, width, height)
	}

	// Output the image
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := noise[y*width+x]
			b := int(c * 255.999)

			fmt.Fprintf(w, "%d %d %d\n", b, b, b)
		}
		fmt.Fprintln(w)
	}

	fmt.Fprintln(w)
}
