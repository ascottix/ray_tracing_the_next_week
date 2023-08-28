package main

import (
	"io"
)

func Image13(w io.Writer) {
	world := NewHittableList()

	noise := NewNoiseTextureWithGenerator(4, NewVectorPerlin())
	material := NewTextureLambertianMaterial(noise)
	world.Add(NewSphere(NewPoint3(0, -1000, 0), 1000, material))
	world.Add(NewSphere(NewPoint3(0, 2, 0), 2, material))

	cam := NewCamera()
	cam.SetLookFrom(NewPoint3(13, 2, 3))
	cam.SetLookAt(NewPoint3(0, 0, 0))
	cam.SetVerticalFieldOfView(20)
	cam.SetDefocusAngle(0)

	cam.Render(w, world)
}
