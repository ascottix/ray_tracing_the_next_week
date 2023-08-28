package main

import "math"

// An instance of a Hittable object that is translated by some offset
type Translate struct {
	object Hittable // The original object is stored as-is, all the magic happens in the Hit() function
	offset Vec3
	bbox   Aabb
}

// An instance of a Hittable object that is rotated along the Y-axis
type RotateY struct {
	object   Hittable // The original object is stored as-is, all the magic happens in the Hit() function
	sinTheta float64
	cosTheta float64
	bbox     Aabb
}

func NewTranslate(object Hittable, offset Vec3) Translate {
	return Translate{object, offset, NewAabb(object.BoundingBox().Min.Add(offset), object.BoundingBox().Max.Add(offset))}
}

func (t Translate) Hit(ray Ray, rayTmin, rayTmax float64, rec *HitRecord) bool {
	// Move the ray backwards by the offset
	offsetRay := NewRay(ray.Origin().Sub(t.offset), ray.Direction(), ray.Time())

	// Determine where (if any) an intersection occurs along the offset ray
	if !t.object.Hit(offsetRay, rayTmin, rayTmax, rec) {
		return false
	}

	// Move the intersection point forwards by the offset
	rec.P = rec.P.Add(t.offset)

	return true
}

func (t Translate) BoundingBox() Aabb {
	return t.bbox
}

func NewRotateY(object Hittable, angleInDegrees float64) RotateY {
	theta := DegreesToRadians(angleInDegrees)
	sinTheta := math.Sin(theta)
	cosTheta := math.Cos(theta)

	// We need to build a new bounding taking into account all rotations
	rotate := func(p Point3) Point3 {
		return NewPoint3(cosTheta*p.X+sinTheta*p.Z, p.Y, -sinTheta*p.X+cosTheta*p.Z)
	}

	p1 := rotate(object.BoundingBox().Min)
	p2 := rotate(object.BoundingBox().Max)

	min := NewPoint3(math.Min(p1.X, p2.X), p1.Y, math.Min(p1.Z, p2.Z))
	max := NewPoint3(math.Max(p1.X, p2.X), p2.Y, math.Max(p1.Z, p2.Z))

	return RotateY{object, sinTheta, cosTheta, NewAabb(min, max)}
}

func (roty RotateY) Hit(ray Ray, rayTmin, rayTmax float64, rec *HitRecord) bool {
	// Change the ray from world space to object space
	rotatedOrigin := NewPoint3(roty.cosTheta*ray.Origin().X-roty.sinTheta*ray.Origin().Z, ray.Origin().Y, roty.sinTheta*ray.Origin().X+roty.cosTheta*ray.Origin().Z)
	rotatedDirection := NewPoint3(roty.cosTheta*ray.Direction().X-roty.sinTheta*ray.Direction().Z, ray.Direction().Y, roty.sinTheta*ray.Direction().X+roty.cosTheta*ray.Direction().Z)
	rotatedRay := NewRay(rotatedOrigin, rotatedDirection, ray.Time())

	// Determine where (if any) an intersection occurs in object space
	if !roty.object.Hit(rotatedRay, rayTmin, rayTmax, rec) {
		return false
	}

	// Change the intersection point from object space to world space
	rec.P = NewPoint3(roty.cosTheta*rec.P.X+roty.sinTheta*rec.P.Z, rec.P.Y, -roty.sinTheta*rec.P.X+roty.cosTheta*rec.P.Z)

	// Change the normal from object space to world space
	rec.Normal = NewVec3(roty.cosTheta*rec.Normal.X+roty.sinTheta*rec.Normal.Z, rec.Normal.Y, -roty.sinTheta*rec.Normal.X+roty.cosTheta*rec.Normal.Z)

	return true
}

func (roty RotateY) BoundingBox() Aabb {
	return roty.bbox
}
