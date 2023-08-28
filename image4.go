package main

import (
	"fmt"
	"io"
)

// The "earthmap.jpg" image can be downloaded directly from the online book (see README)
func Image4(w io.Writer) {
	it := NewImageTexture("earthmap.jpg")

	fmt.Fprintf(w, "P3\n") // Magic
	fmt.Fprintf(w, "%d %d\n", it.width, it.height)
	fmt.Fprintf(w, "255\n") // Maximum value of a color component

	for y := 0; y < it.height; y++ {
		for x := 0; x < it.width; x++ {
			u, v := float64(x)/float64(it.width), float64(y)/float64(it.height)
			c := it.Value(u, 1-v, NewPoint3(0, 0, 0))

			ir := int(255.999 * c.X)
			ig := int(255.999 * c.Y)
			ib := int(255.999 * c.Z)

			fmt.Fprintf(w, "%d %d %d\n", ir, ig, ib)
		}
		fmt.Fprintln(w)
	}

	fmt.Fprintln(w)
}
