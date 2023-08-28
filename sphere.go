package main

import "math"

type Sphere struct {
	center    Point3
	radius    float64
	mat       Material
	centerVec Vec3 // Movement vector
	bbox      Aabb // Bounding box
}

// Stationary sphere
func NewSphere(center Point3, radius float64, mat Material) Sphere {
	rvec := NewVec3(radius, radius, radius)
	return Sphere{center: center, radius: radius, mat: mat, bbox: NewAabb(center.Sub(rvec), center.Add(rvec))}
}

func NewMovingSphere(startCenter, endCenter Point3, radius float64, mat Material) Sphere {
	rvec := NewVec3(radius, radius, radius)
	bbox1 := NewAabb(startCenter.Sub(rvec), startCenter.Add(rvec))
	bbox2 := NewAabb(endCenter.Sub(rvec), endCenter.Add(rvec))
	return Sphere{center: startCenter, radius: radius, mat: mat, centerVec: endCenter.Sub(startCenter), bbox: bbox1.Union(bbox2)}
}

// Returns the (u,v) coordinates of a point on a unit sphere centered at the origin, where:
// 0 <= u <= 1 is the angle around the Y axis, counted starting from X=-1
// 0 <= v <= 1 is the angle from Y=-1 to Y=+1
// Basically we use a spherical coordinate system where Y is the polar axis (aka zenith direction),
// u is the polar angle and v is the azimuthal angle. However, we normalize all angles to be
// in the [0,1] interval so that our texture coordinates system remain consistent across different surfaces.
func getSphereUV(p Point3) (float64, float64) {
	theta := math.Acos(-p.Y)
	phi := math.Atan2(-p.Z, p.X) + math.Pi

	return phi / (2 * math.Pi), theta / math.Pi
}

// Implement the Hittable interface
func (s Sphere) Hit(ray Ray, rayTmin, rayTmax float64, rec *HitRecord) bool {
	// Linearly interpolate from startCenter (time=0) to endCenter (time=1)
	center := s.center.Add(s.centerVec.Mul(ray.Time()))

	oc := ray.Origin().Sub(center)
	a := ray.Direction().Dot(ray.Direction())
	half_b := oc.Dot(ray.Direction())
	c := oc.Dot(oc) - s.radius*s.radius
	discriminant := half_b*half_b - a*c

	if discriminant < 0 {
		return false // No intersection (any point where t < 0 is behind the camera)
	}

	// Find the nearest root that lies in the allowed range
	sqrtd := math.Sqrt(discriminant)
	root := (-half_b - sqrtd) / a // First intersection (closest to the camera)
	if root <= rayTmin || root >= rayTmax {
		root = (-half_b + sqrtd) / a // Second intersection
		if root <= rayTmin || root >= rayTmax {
			return false
		}
	}

	rec.T = root
	rec.P = ray.At(rec.T)
	outwardNormal := rec.P.Sub(center).Div(s.radius) // Divide by the sphere radius as it's cheaper that calling UnitVector() and gets the same result here
	rec.SetFaceNormal(ray, outwardNormal)
	rec.Mat = s.mat
	rec.U, rec.V = getSphereUV(outwardNormal)

	return true
}

func (s Sphere) BoundingBox() Aabb {
	return s.bbox
}
