package main

import (
	"io"
)

// Cornell box with two boxes made of fog
func Image22(w io.Writer) {
	world := NewHittableList()

	white := createEmptyCornellBox(&world)

	light := NewDiffuseLight(NewSolidColorTexture(NewColor(7, 7, 7)))
	world.Add(NewQuad(NewPoint3(113, 553.9, 127), NewVec3(330, 0, 0), NewVec3(0, 0, 305), light))

	box1 := createBox(NewPoint3(0, 0, 0), NewPoint3(165, 330, 165), white)
	box1m := NewTranslate(NewRotateY(box1, 15), NewVec3(265, 0, 295))
	fog1 := NewConstantMedium(box1m, 0.01, NewSolidColorTexture(NewColor(0, 0, 0)))
	world.Add(fog1)
	box2 := createBox(NewPoint3(0, 0, 0), NewPoint3(165, 165, 165), white)
	box2m := NewTranslate(NewRotateY(box2, -18), NewVec3(130, 0, 65))
	fog2 := NewConstantMedium(box2m, 0.01, NewSolidColorTexture(NewColor(1, 1, 1)))
	world.Add(fog2)

	cam := NewCamera()
	cam.SetAspectRatio(1)
	cam.SetLookFrom(NewPoint3(278, 278, -800))
	cam.SetLookAt(NewPoint3(278, 278, 0))
	cam.SetVerticalFieldOfView(40)
	cam.SetBackground(NewColor(0, 0, 0))
	cam.SetRenderingParams(200, 40)

	cam.Render(w, world)
}
