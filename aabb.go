package main

import "math"

type Aabb struct {
	Min, Max Point3 // The reference implementation uses 3 intervals instead
}

func NewAabb(a, b Point3) Aabb {
	return Aabb{
		Min: NewPoint3(math.Min(a.X, b.X), math.Min(a.Y, b.Y), math.Min(a.Z, b.Z)),
		Max: NewPoint3(math.Max(a.X, b.X), math.Max(a.Y, b.Y), math.Max(a.Z, b.Z))}
}

func (aabb Aabb) Union(aabb2 Aabb) Aabb {
	return Aabb{
		NewPoint3(math.Min(aabb2.Min.X, aabb.Min.X), math.Min(aabb2.Min.Y, aabb.Min.Y), math.Min(aabb2.Min.Z, aabb.Min.Z)),
		NewPoint3(math.Max(aabb2.Max.X, aabb.Max.X), math.Max(aabb2.Max.Y, aabb.Max.Y), math.Max(aabb2.Max.Z, aabb.Max.Z)),
	}
}

func (aabb Aabb) Pad() Aabb {
	delta := 0.0001

	pad := func(cmin, cmax float64) (float64, float64) {
		if cmax-cmin < delta {
			cmin -= delta / 2
			cmax += delta / 2
		}
		return cmin, cmax
	}

	xmin, xmax := pad(aabb.Min.X, aabb.Max.X)
	ymin, ymax := pad(aabb.Min.Y, aabb.Max.Y)
	zmin, zmax := pad(aabb.Min.Z, aabb.Max.Z)

	return NewAabb(NewPoint3(xmin, ymin, zmin), NewPoint3(xmax, ymax, zmax))
}

func (aabb Aabb) Hit(ray Ray, tMin, tMax float64) bool {
	checkAxis := func(axisMin, axisMax, axisRayOrig, axisRayDir float64) bool {
		axisRayInvDir := 1 / axisRayDir

		t0 := (axisMin - axisRayOrig) * axisRayInvDir
		t1 := (axisMax - axisRayOrig) * axisRayInvDir

		if axisRayInvDir < 0 {
			t0, t1 = t1, t0
		}

		if t0 > tMin {
			tMin = t0
		}

		if t1 < tMax {
			tMax = t1
		}

		return tMax > tMin
	}

	return checkAxis(aabb.Min.X, aabb.Max.X, ray.Origin().X, ray.Direction().X) &&
		checkAxis(aabb.Min.Y, aabb.Max.Y, ray.Origin().Y, ray.Direction().Y) &&
		checkAxis(aabb.Min.Z, aabb.Max.Z, ray.Origin().Z, ray.Direction().Z)
}
