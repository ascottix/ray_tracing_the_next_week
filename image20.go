package main

import (
	"io"
)

// Returns the 3D box (six sides) that contains the two opposite vertices a and b
func createBox(a, b Point3, mat Material) HittableList {
	sides := NewHittableList()
	bbox := NewAabb(a, b)
	dx := NewVec3(bbox.Max.X-bbox.Min.X, 0, 0)
	dy := NewVec3(0, bbox.Max.Y-bbox.Min.Y, 0)
	dz := NewVec3(0, 0, bbox.Max.Z-bbox.Min.Z)

	sides.Add(NewQuad(NewPoint3(bbox.Min.X, bbox.Min.Y, bbox.Max.Z), dx, dy, mat))          // Front
	sides.Add(NewQuad(NewPoint3(bbox.Max.X, bbox.Min.Y, bbox.Max.Z), dz.Negate(), dy, mat)) // Right
	sides.Add(NewQuad(NewPoint3(bbox.Max.X, bbox.Min.Y, bbox.Min.Z), dx.Negate(), dy, mat)) // Back
	sides.Add(NewQuad(NewPoint3(bbox.Min.X, bbox.Min.Y, bbox.Min.Z), dz, dy, mat))          // Left
	sides.Add(NewQuad(NewPoint3(bbox.Min.X, bbox.Max.Y, bbox.Max.Z), dx, dz.Negate(), mat)) // Top
	sides.Add(NewQuad(NewPoint3(bbox.Min.X, bbox.Min.Y, bbox.Min.Z), dx, dz, mat))          // Bottom

	return sides
}

// Cornell box with two boxes
func Image20(w io.Writer) {
	world := NewHittableList()

	white := createEmptyCornellBox(&world)
	world.Add(createBox(NewPoint3(130, 0, 65), NewPoint3(295, 165, 230), white))
	world.Add(createBox(NewPoint3(265, 0, 295), NewPoint3(430, 330, 460), white))

	cam := NewCamera()
	cam.SetAspectRatio(1)
	cam.SetLookFrom(NewPoint3(278, 278, -800))
	cam.SetLookAt(NewPoint3(278, 278, 0))
	cam.SetVerticalFieldOfView(40)
	cam.SetBackground(NewColor(0, 0, 0))
	cam.SetRenderingParams(800, 50)

	cam.Render(w, world)
}
