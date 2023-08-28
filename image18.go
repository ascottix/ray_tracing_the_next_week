package main

import (
	"io"
)

func Image18(w io.Writer) {
	world := NewHittableList()

	noise := NewMarbleTexture(4)
	material := NewTextureLambertianMaterial(noise)
	world.Add(NewSphere(NewPoint3(0, -1000, 0), 1000, material))
	world.Add(NewSphere(NewPoint3(0, 2, 0), 2, material))

	diffLight := NewDiffuseLight(NewSolidColorTexture(NewColor(4, 4, 4)))
	world.Add(NewSphere(NewPoint3(0, 7, 0), 2, diffLight))
	world.Add(NewQuad(NewPoint3(3, 1, -2), NewVec3(2, 0, 0), NewVec3(0, 2, 0), diffLight))

	cam := NewCamera()
	cam.SetLookFrom(NewPoint3(26, 3, 6))
	cam.SetLookAt(NewPoint3(0, 2, 0))
	cam.SetVerticalFieldOfView(20)
	cam.SetBackground(NewColor(0, 0, 0))

	cam.Render(w, world)
}
