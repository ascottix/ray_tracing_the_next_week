package main

import "math"

type Material interface {
	// Returns true if the surface scattered (reflected) the incoming ray, or false if it has absorbed it.
	// If the ray has been scattered, also returns the scattered ray and the attenuation color (which depends on the material).
	Scatter(ray Ray, rec *HitRecord, attenuation *Color, scattered *Ray) bool
	Emitted(u, v float64, p Point3) Color
}

// Provides the Emitted() method to materials that don't need to implement it
type BlackEmitter struct {
}

func (b BlackEmitter) Emitted(u, v float64, p Point3) Color {
	return Color{}
}

// A metal material reflects the light according to the direction of the incident ray
type MetalMaterial struct {
	BlackEmitter
	albedo Color
	fuzz   float64 // If fuzz is 0 the material is perfectly smooth, setting 0 < fuzz <= 1 adds roughness to the surface
}

// A Lambertial (or matte) material scatters the light in a random direction, with a distribution biased towards the surface normal
type TextureLambertianMaterial struct {
	BlackEmitter
	texture Texture
}

// An isotropic material scatters the light in a random direction, with a uniform distribution
type IsotropicMaterial struct {
	BlackEmitter
	albedo Texture
}

// Metal material
func NewMetalMaterial(a Color, fuzz float64) MetalMaterial {
	return MetalMaterial{albedo: a, fuzz: math.Min(fuzz, 1)}
}

func Reflect(v, n Vec3) Vec3 {
	return v.Add(n.Mul(-2 * n.Dot(v)))
}

func Refract(uv, n Vec3, etaiOverEtat float64) Vec3 {
	cosTheta := n.Dot(uv.Negate()) // It's math.Min(n.Dot(uv.Negate()), 1.0) in the original source
	rOutPerp := uv.Add(n.Mul(cosTheta)).Mul(etaiOverEtat)
	rOutParallel := n.Mul(-math.Sqrt(math.Abs(1.0 - rOutPerp.LengthSquared())))
	return rOutPerp.Add(rOutParallel)
}

func (m MetalMaterial) Scatter(ray Ray, rec *HitRecord, attenuation *Color, scattered *Ray) bool {
	reflected := Reflect(ray.Direction().UnitVector(), rec.Normal)

	*scattered = NewRay(rec.P, reflected.Add(NewRandomUnitVec3().Mul(m.fuzz)), ray.Time())
	*attenuation = m.albedo

	// We should just return true here, but because of the fuzziness it may happen that a ray is scattered below the surface.
	// If that happens, just pretend the surface has absorbed it and don't scatter.
	return rec.Normal.Dot(scattered.Direction()) > 0
}

// Dielectric material
type DielectricMaterial struct {
	BlackEmitter
	ir float64
}

func NewDielectricMaterial(indexOfRefraction float64) DielectricMaterial {
	return DielectricMaterial{ir: indexOfRefraction}
}

// Schlick's approximation for reflectance,
// it is used to handle total internal reflection
func SchlickReflectance(cosine, refIdx float64) float64 {
	r0 := (1 - refIdx) / (1 + refIdx)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow(1-cosine, 5)
}

func (m DielectricMaterial) Scatter(ray Ray, rec *HitRecord, attenuation *Color, scattered *Ray) bool {
	refractionRatio := m.ir
	if rec.FrontFace {
		refractionRatio = 1 / refractionRatio
	}

	unitDirection := ray.Direction().UnitVector()

	cosTheta := rec.Normal.Dot(unitDirection.Negate())
	sinTheta := math.Sqrt(1 - cosTheta*cosTheta)

	cannotRefract := (refractionRatio*sinTheta > 1) || SchlickReflectance(cosTheta, refractionRatio) >= RandomDouble()

	if cannotRefract {
		reflected := Reflect(unitDirection, rec.Normal)
		*scattered = NewRay(rec.P, reflected, ray.Time())
	} else {
		refracted := Refract(unitDirection, rec.Normal, refractionRatio)
		*scattered = NewRay(rec.P, refracted, ray.Time())
	}

	*attenuation = Color{1, 1, 1}

	return true
}

// Lambertian material based on texture
func NewTextureLambertianMaterial(texture Texture) TextureLambertianMaterial {
	return TextureLambertianMaterial{texture: texture}
}

func (m TextureLambertianMaterial) Scatter(ray Ray, rec *HitRecord, attenuation *Color, scattered *Ray) bool {
	scatterDirection := rec.Normal.Add(NewRandomUnitVec3())

	// Catch an edge case where the random unit vector is exactly opposite to the surface normal and nullifies the scatter direction
	if scatterDirection.NearZero() {
		scatterDirection = rec.Normal
	}

	*scattered = NewRay(rec.P, scatterDirection, ray.Time())
	*attenuation = m.texture.Value(rec.U, rec.V, rec.P)

	return true
}

// Lambertian material based on color
func NewLambertianMaterial(a Color) TextureLambertianMaterial {
	return NewTextureLambertianMaterial(NewSolidColorTexture(a))
}

// Isotropic material
func NewIsotropicMaterial(a Texture) IsotropicMaterial {
	return IsotropicMaterial{albedo: a}
}

func (m IsotropicMaterial) Scatter(ray Ray, rec *HitRecord, attenuation *Color, scattered *Ray) bool {
	*scattered = NewRay(rec.P, NewRandomUnitVec3(), ray.Time())
	*attenuation = m.albedo.Value(rec.U, rec.V, rec.P)
	return true
}
