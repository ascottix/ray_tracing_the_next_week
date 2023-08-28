package main

import (
	"math"
)

type Quad struct {
	Q      Point3 // Reference corner
	u      Vec3   // Vector to 2nd corner (the 2nd corner is placed at Q+u)
	v      Vec3   // Vector to 3rd corner (the 3rd corner is placed at Q+v, the 4th corner is at Q+u+v)
	mat    Material
	bbox   Aabb
	normal Vec3
	D      float64
	w      Vec3
}

func NewQuad(q Point3, u, v Vec3, mat Material) Quad {
	bbox := NewAabb(q, q.Add(u).Add(v)).Pad()
	n := u.Cross(v)
	normal := n.UnitVector()
	d := normal.Dot(q)
	w := n.Div(n.Dot(n))

	return Quad{q, u, v, mat, bbox, normal, d, w}
}

// Implement the Hittable interface
func (quad Quad) Hit(ray Ray, rayTmin, rayTmax float64, rec *HitRecord) bool {
	denom := quad.normal.Dot(ray.Direction())

	// If the ray is parallel to the plane there is no intersection
	if math.Abs(denom) < 1e-8 {
		return false
	}

	// Compute t and check if it's within range
	t := (quad.D - quad.normal.Dot(ray.Origin())) / denom
	if t < rayTmin || t > rayTmax {
		return false
	}

	// Got a candidate intersection point
	intersection := ray.At(t)

	planarHitpointVector := intersection.Sub(quad.Q)
	alpha := quad.w.Dot(planarHitpointVector.Cross(quad.v))
	beta := quad.w.Dot(quad.u.Cross(planarHitpointVector))

	if alpha >= 0 && alpha <= 1 && beta >= 0 && beta <= 1 {
		// Got an actual intersection
		rec.T = t
		rec.P = intersection
		rec.Mat = quad.mat
		rec.SetFaceNormal(ray, quad.normal)
		rec.U = alpha
		rec.V = beta

		return true
	}

	return false
}

func (quad Quad) BoundingBox() Aabb {
	return quad.bbox
}
