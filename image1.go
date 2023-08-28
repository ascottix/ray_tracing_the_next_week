package main

import (
	"io"
)

func addRandomSpheresToWorld(world *HittableList) {
	ref := NewPoint3(4, 0.2, 0)
	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			center := NewPoint3(float64(a)+0.9*RandomDouble(), 0.2, float64(b)+0.9*RandomDouble())

			if center.Sub(ref).Length() > 0.9 {
				chooseMat := RandomDouble()
				if chooseMat < 0.8 {
					// Diffuse
					center2 := center.Add(NewVec3(0, RandomDoubleInInterval(0, 0.5), 0)) // Comment out the .Add(...) part to prevent the sphere from moving
					albedo := NewRandomVec3().MultiplyByComponent(NewRandomVec3())
					mat := NewLambertianMaterial(albedo)
					world.Add(NewMovingSphere(center, center2, 0.2, mat))
				} else if chooseMat < 0.95 {
					// Metal
					albedo := NewRandomInIntervalVec3(0.5, 1)
					fuzz := RandomDoubleInInterval(0, 0.5)
					mat := NewMetalMaterial(albedo, fuzz)
					world.Add(NewSphere(center, 0.2, mat))
				} else {
					// Glass
					mat := NewDielectricMaterial(1.5)
					world.Add(NewSphere(center, 0.2, mat))
				}
			}
		}
	}
}

func Image1(w io.Writer) {
	world := NewHittableList()

	materialGround := NewLambertianMaterial(NewColor(0.5, 0.5, 0.5))
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
