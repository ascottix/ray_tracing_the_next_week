package main

import (
	"math/rand"
	"slices"
)

type BvhNode struct {
	left  Hittable
	right Hittable
	bbox  Aabb
}

func NewBhvTree(list HittableList) BvhNode {
	return NewBvhNode(list.objects) // Can also use NewBvhNodeBook(list.objects, 0, len(list.objects))
}

type Comparator func(a, b Hittable) int // Unlike C++'s std::sort(), comparators need to return int instead of bool

var boolToInt = map[bool]int{false: 1, true: -1}

func boxCompareX(a, b Hittable) int {
	return boolToInt[a.BoundingBox().Min.X < b.BoundingBox().Min.X]
}

func boxCompareY(a, b Hittable) int {
	return boolToInt[a.BoundingBox().Min.Y < b.BoundingBox().Min.Y]
}

func boxCompareZ(a, b Hittable) int {
	return boolToInt[a.BoundingBox().Min.Z < b.BoundingBox().Min.Z]
}

var boxComparators = []Comparator{boxCompareX, boxCompareY, boxCompareZ}

func getRandomBoxComparator(axis int) Comparator {
	return boxComparators[axis]
}

func NewBvhNode(objects []Hittable) BvhNode {
	var left, right Hittable

	if len(objects) == 1 {
		left, right = objects[0], objects[0]
	} else if len(objects) == 2 {
		left, right = objects[0], objects[1]
	} else {
		// Split the list in half along a random axis
		slices.SortFunc(objects, getRandomBoxComparator(rand.Intn(3)))

		mid := len(objects) / 2

		left, right = NewBvhNode(objects[:mid]), NewBvhNode(objects[mid:])
	}

	return BvhNode{left: left, right: right, bbox: left.BoundingBox().Union(right.BoundingBox())}
}

// This version resembles the book's C++ code and works fine, but doesn't take advantage of Go slices
func NewBvhNodeBook(objects []Hittable, start, end int) BvhNode {
	var left, right Hittable

	comparator := getRandomBoxComparator(rand.Intn(3))

	objectSpan := end - start

	if objectSpan == 1 {
		left = objects[start]
		right = left
	} else if objectSpan == 2 {
		left = objects[start]
		right = objects[start+1]

		if comparator(right, left) < 0 {
			left, right = right, left
		}
	} else {
		slices.SortFunc(objects[start:end], comparator)

		mid := start + objectSpan/2

		left = NewBvhNodeBook(objects, start, mid)
		right = NewBvhNodeBook(objects, mid, end)
	}

	return BvhNode{left: left, right: right, bbox: left.BoundingBox().Union(right.BoundingBox())}
}

func (node BvhNode) Hit(ray Ray, rayTmin, rayTmax float64, rec *HitRecord) bool {
	if !node.bbox.Hit(ray, rayTmin, rayTmax) {
		return false
	}

	hitLeft := node.left.Hit(ray, rayTmin, rayTmax, rec)

	if hitLeft { // Update the ray max extent as we're not interested in hits that are farther away than this
		rayTmax = rec.T
	}

	hitRight := node.right.Hit(ray, rayTmin, rayTmax, rec)

	return hitLeft || hitRight
}

func (node BvhNode) BoundingBox() Aabb {
	return node.bbox
}
