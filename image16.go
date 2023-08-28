package main

import (
	"io"
)

func Image16(w io.Writer) {
	world := NewHittableList()

	// Materials
	leftRed := NewLambertianMaterial(NewColor(1.0, 0.2, 0.2))
	backGreen := NewLambertianMaterial(NewColor(0.2, 1.0, 0.2))
	rightBlue := NewLambertianMaterial(NewColor(0.2, 0.2, 1.0))
	upperOrange := NewLambertianMaterial(NewColor(1.0, 0.5, 0.0))
	lowerTeal := NewLambertianMaterial(NewColor(0.2, 0.8, 0.8))

	// Quads
	world.Add(NewQuad(NewPoint3(-3, -2, 5), NewVec3(0, 0, -4), NewVec3(0, 4, 0), leftRed))
	world.Add(NewQuad(NewPoint3(-2, -2, 0), NewVec3(4, 0, 0), NewVec3(0, 4, 0), backGreen))
	world.Add(NewQuad(NewPoint3(3, -2, 1), NewVec3(0, 0, 4), NewVec3(0, 4, 0), rightBlue))
	world.Add(NewQuad(NewPoint3(-2, 3, 1), NewVec3(4, 0, 0), NewVec3(0, 0, 4), upperOrange))
	world.Add(NewQuad(NewPoint3(-2, -3, 5), NewVec3(4, 0, 0), NewVec3(0, 0, -4), lowerTeal))

	cam := NewCamera()
	cam.SetLookFrom(NewPoint3(0, 0, 9))
	cam.SetLookAt(NewPoint3(0, 0, 0))
	cam.SetVerticalFieldOfView(80)
	cam.SetDefocusAngle(0)
	cam.SetAspectRatio(1)

	cam.Render(w, world)
}
