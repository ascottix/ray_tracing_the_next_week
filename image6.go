package main

import (
	"fmt"
	"io"
)

func Image6(w io.Writer) {
	width, height := 400, 225

	fmt.Fprintf(w, "P3\n") // Magic
	fmt.Fprintf(w, "%d %d\n", width, height)
	fmt.Fprintf(w, "255\n") // Maximum value of a color component

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			b := int(RandomDouble() * 255.999)

			fmt.Fprintf(w, "%d %d %d\n", b, b, b)
		}
		fmt.Fprintln(w)
	}

	fmt.Fprintln(w)
}
