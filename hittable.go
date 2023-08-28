package main

type HitRecord struct {
	P         Point3   // Hit point on surface
	Normal    Vec3     // Normal to surface at point P
	T         float64  // Ray extension at hit point
	FrontFace bool     // Whether the ray hit the surface from outside (true) or inside (false)
	Mat       Material // Surface material
	U, V      float64  // Coordinates of hit point relative to surface
}

type Hittable interface {
	Hit(ray Ray, rayTmin, rayTmax float64, rec *HitRecord) bool

	BoundingBox() Aabb
}

// A normal to an object surface may point outwards or inwards... how do we choose?
// There are two main conventions:
// 1. the normal always points outwards
// 2. the normal always points against the ray
// In the first case we can use the dot product between ray and normal to determine whether
// the ray is outside the object (dot product is negative) or inside the object (dot product is positive).
// In the second case the dot product will always be negative so we need to determine that information
// first, and store it for later.
// This book takes the second approach.
// Note: outwardNormal is assumed to have unit length
func (h *HitRecord) SetFaceNormal(ray Ray, outwardNormal Vec3) {
	h.FrontFace = ray.Direction().Dot(outwardNormal) < 0
	if h.FrontFace {
		// Ray is outside the sphere
		h.Normal = outwardNormal
	} else {
		// Ray is inside the sphere, flip the normal
		h.Normal = outwardNormal.Negate()
	}
}
