package main

import (
	"io"
)

func Image5(w io.Writer) {
	earthTexture := NewImageTexture("earthmap.jpg")
	earthSurface := NewTextureLambertianMaterial(earthTexture)
	globe := NewSphere(NewPoint3(0, 0, 0), 2, earthSurface)

	cam := NewCamera()
	cam.SetLookFrom(NewPoint3(14, 2, 3)) // It's (0,0,12) in the book, but that doesn't look like the picture
	cam.SetLookAt(NewPoint3(0, 0, 0))
	cam.SetVerticalFieldOfView(20)

	world := NewHittableList()
	world.Add(globe)

	cam.Render(w, world)
}
