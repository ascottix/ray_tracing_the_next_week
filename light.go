package main

type DiffuseLight struct {
	emit Texture
}

func NewDiffuseLight(texture Texture) DiffuseLight {
	return DiffuseLight{texture}
}

func (dl DiffuseLight) Scatter(ray Ray, rec *HitRecord, attenuation *Color, scattered *Ray) bool {
	return false
}

func (dl DiffuseLight) Emitted(u, v float64, p Point3) Color {
	return dl.emit.Value(u, v, p)
}
