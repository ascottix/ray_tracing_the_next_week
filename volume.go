package main

import (
	"fmt"
	"math"
	"os"
)

type ConstantMedium struct {
	boundary      Hittable
	negInvDensity float64
	phaseFunction Material
}

func NewConstantMedium(b Hittable, d float64, a Texture) ConstantMedium {
	return ConstantMedium{boundary: b, negInvDensity: -1 / d, phaseFunction: NewIsotropicMaterial(a)}
}

func (cm ConstantMedium) Hit(ray Ray, rayTmin, rayTmax float64, rec *HitRecord) bool {
	// Print occasional samples when debugging. To enable, set enableDebug true.
	enableDebug := false
	debugging := enableDebug && RandomDouble() < 0.00001

	var rec1, rec2 HitRecord

	if !cm.boundary.Hit(ray, math.Inf(-1), math.Inf(+1), &rec1) {
		return false
	}

	if !cm.boundary.Hit(ray, rec1.T+0.0001, math.Inf(+1), &rec2) {
		return false
	}

	if debugging {
		fmt.Fprintf(os.Stderr, "\nrayTmin=%f, rayTMax=%f\n", rec1.T, rec2.T)
	}

	if rec1.T < 0 {
		rec1.T = 0
	}

	if rec1.T < rayTmin {
		rec1.T = rayTmin
	}

	if rec2.T > rayTmax {
		rec2.T = rayTmax
	}

	if rec1.T >= rec2.T {
		return false
	}

	rayLength := ray.Direction().Length()
	distanceInsideBoundary := (rec2.T - rec1.T) * rayLength
	hitDistance := cm.negInvDensity * math.Log(RandomDouble())

	if hitDistance > distanceInsideBoundary {
		return false
	}

	rec.T = rec1.T + hitDistance/rayLength
	rec.P = ray.At(rec.T)

	if debugging {
		fmt.Fprintln(os.Stderr, "hitDistance =", hitDistance)
		fmt.Fprintln(os.Stderr, "rec.T =", rec.T)
		fmt.Fprintln(os.Stderr, "rec.P =", rec.P)
	}

	rec.Normal = NewVec3(1, 0, 0) // Arbitrary
	rec.FrontFace = true          // Arbitrary
	rec.Mat = cm.phaseFunction

	return true
}

func (cm ConstantMedium) BoundingBox() Aabb {
	return cm.boundary.BoundingBox()
}
