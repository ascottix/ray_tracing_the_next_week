package main

import (
	"io"
)

func Image23(w io.Writer) {
	world := NewHittableList()

	boxes1 := NewHittableList()
	ground := NewLambertianMaterial(NewColor(0.48, 0.83, 0.53))

	const boxesPerSide = 20
	for i := 0; i < boxesPerSide; i++ {
		for j := 0; j < boxesPerSide; j++ {
			w := 100.0
			x0, y0, z0 := -1000+float64(i)*w, 0.0, -1000+float64(j)*w
			x1, y1, z1 := x0+w, RandomDoubleInInterval(1, 101), z0+w

			boxes1.Add(createBox(NewPoint3(x0, y0, z0), NewPoint3(x1, y1, z1), ground))
		}
	}

	world.Add(NewBhvTree(boxes1))

	light := NewDiffuseLight(NewSolidColorTexture(NewColor(7, 7, 7)))
	world.Add(NewQuad(NewPoint3(123, 554, 147), NewVec3(300, 0, 0), NewVec3(0, 0, 265), light))

	center1 := NewPoint3(400, 400, 200)
	center2 := center1.Add(NewVec3(30, 0, 0))
	sphereMaterial := NewLambertianMaterial(NewColor(0.7, 0.3, 0.1))
	world.Add(NewMovingSphere(center1, center2, 50, sphereMaterial))

	world.Add(NewSphere(NewPoint3(260, 150, 45), 50, NewDielectricMaterial(1.5)))
	world.Add(NewSphere(NewPoint3(0, 150, 145), 50, NewMetalMaterial(NewColor(0.8, 0.8, 0.9), 1)))

	boundary := NewSphere(NewPoint3(360, 150, 145), 70, NewDielectricMaterial(1.5))
	world.Add(boundary)
	world.Add(NewConstantMedium(boundary, 0.2, NewSolidColorTexture(NewColor(0.2, 0.4, 0.9))))
	boundary = NewSphere(NewPoint3(0, 0, 0), 5000, NewDielectricMaterial(1.5))
	world.Add(NewConstantMedium(boundary, 0.0001, NewSolidColorTexture(NewColor(1, 1, 1))))

	emat := NewTextureLambertianMaterial(NewImageTexture("earthmap.jpg"))
	world.Add(NewSphere(NewPoint3(400, 200, 400), 100, emat))
	pertext := NewNoiseTexture(0.1)
	world.Add(NewSphere(NewPoint3(220, 280, 300), 80, NewTextureLambertianMaterial(pertext)))

	boxes2 := NewHittableList()
	white := NewLambertianMaterial(NewColor(0.73, 0.73, 0.73))
	const ns = 1000
	for j := 0; j < ns; j++ {
		x := RandomDoubleInInterval(0, 165)
		y := RandomDoubleInInterval(0, 165)
		z := RandomDoubleInInterval(0, 165)

		boxes2.Add(NewSphere(NewPoint3(x, y, z), 10, white))
	}

	world.Add(NewTranslate(NewRotateY(NewBhvTree(boxes2), 15), NewVec3(-100, 270, 395)))

	cam := NewCamera()
	cam.SetAspectRatio(1)
	cam.SetLookFrom(NewPoint3(478, 278, -600))
	cam.SetLookAt(NewPoint3(278, 278, 0))
	cam.SetVerticalFieldOfView(40)
	cam.SetBackground(NewColor(0, 0, 0))
	cam.SetRenderingParams(10000, 40)

	cam.Render(w, world)
}
