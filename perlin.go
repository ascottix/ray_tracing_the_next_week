package main

import (
	"math"
	"math/rand"
)

const (
	NoiseNoInterpolation = iota
	NoiseTrilinearInterpolation
	NoiseTrilinearInterpolationWithHermitianSmoothing
)

type NoiseGenerator interface {
	Noise(p Point3) float64
}

type Perlin struct {
	ranfloat []float64
	perm_x   []int
	perm_y   []int
	perm_z   []int
	mode     int
}

type VectorPerlin struct {
	ranvec []Vec3
	perm_x []int
	perm_y []int
	perm_z []int
}

type TurbulenceNoise struct {
	perlin VectorPerlin
	depth  int
}

type TurbulenceNoiseWithPhase struct {
	tn    TurbulenceNoise
	amp   float64
	scale float64
}

func _generatePerm() []int {
	p := make([]int, 256)

	for i := range p {
		p[i] = i
	}

	rand.Shuffle(len(p), func(i, j int) { // Permute the array
		p[i], p[j] = p[j], p[i]
	})

	return p
}

func NewPerlin(mode int) NoiseGenerator {
	ranfloat := make([]float64, 256)
	for i := range ranfloat {
		ranfloat[i] = RandomDouble()
	}

	perm_x := _generatePerm()
	perm_y := _generatePerm()
	perm_z := _generatePerm()

	return Perlin{ranfloat, perm_x, perm_y, perm_z, mode}
}

func (perlin Perlin) Noise(p Point3) float64 {
	if perlin.mode == NoiseNoInterpolation {
		f := func(coord float64) int {
			return int(coord*4) & 255
		}

		i, j, k := f(p.X), f(p.Y), f(p.Z)

		return perlin.ranfloat[perlin.perm_x[i]^perlin.perm_y[j]^perlin.perm_z[k]]
	} else {
		i, u := int(math.Floor(p.X)), p.X-math.Floor(p.X)
		j, v := int(math.Floor(p.Y)), p.Y-math.Floor(p.Y)
		k, w := int(math.Floor(p.Z)), p.Z-math.Floor(p.Z)

		if perlin.mode == NoiseTrilinearInterpolationWithHermitianSmoothing {
			u = u * u * (3 - 2*u)
			v = v * v * (3 - 2*v)
			w = w * w * (3 - 2*w)
		}

		var c [2][2][2]float64

		for di := 0; di < 2; di++ {
			for dj := 0; dj < 2; dj++ {
				for dk := 0; dk < 2; dk++ {
					idx := perlin.perm_x[(i+di)&255] ^ perlin.perm_y[(j+dj)&255] ^ perlin.perm_z[(k+dk)&255]
					c[di][dj][dk] = perlin.ranfloat[idx]
				}
			}
		}

		// Trilinear interpolation
		v1 := c[1][1][0]*u*v + c[0][1][0]*(1-u)*v + c[1][0][0]*u*(1-v) + c[0][0][0]*(1-u)*(1-v) // Bilinear on k=0
		v2 := c[1][1][1]*u*v + c[0][1][1]*(1-u)*v + c[1][0][1]*u*(1-v) + c[0][0][1]*(1-u)*(1-v) // Bilinear on k=1

		return v1*(1-w) + v2*w
	}
}

func NewVectorPerlin() VectorPerlin {
	ranvec := make([]Vec3, 256)
	for i := range ranvec {
		ranvec[i] = NewRandomInIntervalVec3(-1, 1).UnitVector()
	}

	perm_x := _generatePerm()
	perm_y := _generatePerm()
	perm_z := _generatePerm()

	return VectorPerlin{ranvec, perm_x, perm_y, perm_z}
}

func (perlin VectorPerlin) Noise(p Point3) float64 {
	i, u := int(math.Floor(p.X)), p.X-math.Floor(p.X)
	j, v := int(math.Floor(p.Y)), p.Y-math.Floor(p.Y)
	k, w := int(math.Floor(p.Z)), p.Z-math.Floor(p.Z)

	var c [2][2][2]Vec3

	for di := 0; di < 2; di++ {
		for dj := 0; dj < 2; dj++ {
			for dk := 0; dk < 2; dk++ {
				idx := perlin.perm_x[(i+di)&255] ^ perlin.perm_y[(j+dj)&255] ^ perlin.perm_z[(k+dk)&255]
				c[di][dj][dk] = perlin.ranvec[idx]
			}
		}
	}

	// Trilinear interpolation
	accum := 0.0

	uu := u * u * (3 - 2*u) // Adjust the interpolation factors with an Hermitian cubic function
	vv := v * v * (3 - 2*v)
	ww := w * w * (3 - 2*w)

	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			for k := 0; k < 2; k++ {
				fi, fj, fk := float64(i), float64(j), float64(k)
				weight := NewVec3(u-fi, v-fj, w-fk)
				accum += (fi*uu + (1-fi)*(1-uu)) *
					(fj*vv + (1-fj)*(1-vv)) *
					(fk*ww + (1-fk)*(1-ww)) *
					weight.Dot(c[i][j][k])
			}
		}
	}

	// This kind of interpolation may produce negative values, make sure we return a positive value between 0 and 1
	return (accum + 1) / 2
}

// Turbulence is the sum of noice at different frequencies
func (perlin VectorPerlin) Turbulence(p Point3, depth int) float64 {
	accum := 0.0
	weight := 1.0

	for i := 0; i < depth; i++ {
		noise := 2*perlin.Noise(p) - 1 // Here we can use the negative values, so get back the "original" value
		accum += weight * noise
		weight *= 0.5
		p = p.Mul(2)
	}

	return math.Abs(accum)
}

func NewTurbulenceNoise(depth int) TurbulenceNoise {
	return TurbulenceNoise{perlin: NewVectorPerlin(), depth: depth}
}

func (tn TurbulenceNoise) Noise(p Point3) float64 {
	return tn.perlin.Turbulence(p, tn.depth)
}

func NewTurbulenceNoiseWithPhase(amp, scale float64, depth int) TurbulenceNoiseWithPhase {
	return TurbulenceNoiseWithPhase{tn: NewTurbulenceNoise(depth), amp: amp, scale: scale}
}

func (tnwp TurbulenceNoiseWithPhase) Noise(p Point3) float64 {
	phase := tnwp.amp * tnwp.tn.perlin.Turbulence(p.Mul(tnwp.scale), tnwp.tn.depth)
	return (1 + math.Sin(p.Z+phase)) / 2
}
