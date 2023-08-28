package main

import (
	"io"
)

func Image8(w io.Writer) {
	world := NewHittableList()

	checker := NewRandomBlockTexture(3.2)
	material := NewTextureLambertianMaterial(checker)
	world.Add(NewSphere(NewPoint3(0, -1000, 0), 1000, material))
	world.Add(NewSphere(NewPoint3(0, 2, 0), 2, material))

	cam := NewCamera()
	cam.SetLookFrom(NewPoint3(13, 2, 3))
	cam.SetLookAt(NewPoint3(0, 0, 0))
	cam.SetVerticalFieldOfView(20)
	cam.SetDefocusAngle(0)

	cam.Render(w, world)
}
