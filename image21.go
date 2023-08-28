package main

import (
	"io"
)

// Cornell box with two rotated boxes
func Image21(w io.Writer) {
	world := NewHittableList()

	white := createEmptyCornellBox(&world)
	box1 := createBox(NewPoint3(0, 0, 0), NewPoint3(165, 330, 165), white)
	world.Add(NewTranslate(NewRotateY(box1, 15), NewVec3(265, 0, 295)))
	box2 := createBox(NewPoint3(0, 0, 0), NewPoint3(165, 165, 165), white)
	world.Add(NewTranslate(NewRotateY(box2, -18), NewVec3(130, 0, 65)))

	cam := NewCamera()
	cam.SetAspectRatio(1)
	cam.SetLookFrom(NewPoint3(278, 278, -800))
	cam.SetLookAt(NewPoint3(278, 278, 0))
	cam.SetVerticalFieldOfView(40)
	cam.SetBackground(NewColor(0, 0, 0))
	cam.SetRenderingParams(100, 10)

	cam.Render(w, world)
}
