package main

import (
	"io"
)

func Image3(w io.Writer) {
	world := NewHittableList()

	checker := NewBicolorCheckerTexture(3.2, NewColor(.2, .3, .1), NewColor(.9, .9, .9))
	materialGround := NewTextureLambertianMaterial(checker)
	world.Add(NewSphere(NewPoint3(0, -10, 0), 10, materialGround))
	world.Add(NewSphere(NewPoint3(0, 10, 0), 10, materialGround))

	cam := NewCamera()
	cam.SetLookFrom(NewPoint3(13, 2, 3))
	cam.SetLookAt(NewPoint3(0, 0, 0))
	cam.SetVerticalFieldOfView(20)
	cam.SetFocusDistance(10)
	cam.SetDefocusAngle(0)

	world_bvh := NewBhvTree(world)

	cam.Render(w, world_bvh)
}
