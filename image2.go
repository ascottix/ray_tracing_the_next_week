package main

import (
	"io"
)

func Image2(w io.Writer) {
	world := NewHittableList()

	checker := NewBicolorCheckerTexture(3.2, NewColor(.2, .3, .1), NewColor(.9, .9, .9))
	materialGround := NewTextureLambertianMaterial(checker)
	world.Add(NewSphere(NewPoint3(0.0, -1000, 0), 1000, materialGround))

	addRandomSpheresToWorld(&world)

	material1 := NewDielectricMaterial(1.5)
	world.Add(NewSphere(NewPoint3(0, 1, 0), 1, material1))
	material2 := NewLambertianMaterial(NewColor(0.4, 0.2, 0.1))
	world.Add(NewSphere(NewPoint3(-4, 1, 0), 1, material2))
	material3 := NewMetalMaterial(NewColor(0.7, 0.6, 0.5), 0)
	world.Add(NewSphere(NewPoint3(4, 1, 0), 1, material3))

	cam := NewCamera()
	cam.SetLookFrom(NewPoint3(13, 2, 3))
	cam.SetLookAt(NewPoint3(0, 0, 0))
	cam.SetVerticalFieldOfView(20)
	cam.SetFocusDistance(10)
	cam.SetDefocusAngle(0.02)

	world_bvh := NewBhvTree(world)

	cam.Render(w, world_bvh)
}
