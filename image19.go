package main

import (
	"io"
)

func createEmptyCornellBox(world *HittableList) Material {
	red := NewLambertianMaterial(NewColor(0.65, 0.05, 0.05))
	green := NewLambertianMaterial(NewColor(0.12, 0.45, 0.15))
	white := NewLambertianMaterial(NewColor(0.73, 0.73, 0.73))
	light := NewDiffuseLight(NewSolidColorTexture(NewColor(15, 15, 15)))

	world.Add(NewQuad(NewPoint3(555, 0, 0), NewVec3(0, 555, 0), NewVec3(0, 0, 555), green))
	world.Add(NewQuad(NewPoint3(0, 0, 0), NewVec3(0, 555, 0), NewVec3(0, 0, 555), red))
	world.Add(NewQuad(NewPoint3(343, 554, 332), NewVec3(-130, 0, 0), NewVec3(0, 0, -105), light))
	world.Add(NewQuad(NewPoint3(0, 0, 0), NewVec3(555, 0, 0), NewVec3(0, 0, 555), white))
	world.Add(NewQuad(NewPoint3(555, 555, 555), NewVec3(-555, 0, 0), NewVec3(0, 0, -555), white))
	world.Add(NewQuad(NewPoint3(0, 0, 555), NewVec3(555, 0, 0), NewVec3(0, 555, 0), white))

	return white
}

// Empty Cornell box
func Image19(w io.Writer) {
	world := NewHittableList()

	createEmptyCornellBox(&world)

	cam := NewCamera()
	cam.SetAspectRatio(1)
	cam.SetLookFrom(NewPoint3(278, 278, -800))
	cam.SetLookAt(NewPoint3(278, 278, 0))
	cam.SetVerticalFieldOfView(40)
	cam.SetBackground(NewColor(0, 0, 0))
	cam.SetImageWidth(400)
	cam.SetRenderingParams(200, 50)

	cam.Render(w, world)
}
