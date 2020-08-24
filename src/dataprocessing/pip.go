package dataprocessing

import (
	"math"
)

// After https://github.com/paulmach/orb
func rayCast(point, s, e []float64) (bool, bool) {
	if s[0] > e[0] {
		s, e = e, s
	}

	if point[0] == s[0] {
		if point[1] == s[1] {
			// point == start
			return false, true
		} else if s[0] == e[0] {
			// vertical segment (s -> e)
			// return true if within the line, check to see if start or end is greater.
			if s[1] > e[1] && s[1] >= point[1] && point[1] >= e[1] {
				return false, true
			}

			if e[1] > s[1] && e[1] >= point[1] && point[1] >= s[1] {
				return false, true
			}
		}

		// Move the y coordinate to deal with degenerate case
		point[0] = math.Nextafter(point[0], math.Inf(1))
	} else if point[0] == e[0] {
		if point[1] == e[1] {
			// matching the end point
			return false, true
		}

		point[0] = math.Nextafter(point[0], math.Inf(1))
	}

	if point[0] < s[0] || point[0] > e[0] {
		return false, false
	}

	if s[1] > e[1] {
		if point[1] > s[1] {
			return false, false
		} else if point[1] < e[1] {
			return true, false
		}
	} else {
		if point[1] > e[1] {
			return false, false
		} else if point[1] < s[1] {
			return true, false
		}
	}

	rs := (point[1] - s[1]) / (point[0] - s[0])
	ds := (e[1] - s[1]) / (e[0] - s[0])

	if rs == ds {
		return false, true
	}

	return rs <= ds, false
}

// After https://github.com/paulmach/orb
func polygonContains(polygon *[][]float64, point []float64) bool {
	b, on := rayCast(point, (*polygon)[0], (*polygon)[len(*polygon)-1])
	if on {
		return true
	}

	for i := 0; i < len(*polygon)-1; i++ {
		inter, on := rayCast(point, (*polygon)[i], (*polygon)[i+1])
		if on {
			return true
		}
		if inter {
			b = !b
		}
	}
	return b
}

func isLand(tree *boundingTree, point []float64, polygons *Polygons) bool {
	land := false
	if boundingContains(&tree.boundingBox, point) {
		for _, child := range (*tree).children {
			land = isLand(&child, point, polygons)
			if land {
				return land
			}
		}
		if (*tree).id >= 0 {
			land = polygonContains(&(*polygons)[(*tree).id].Points, point)
		}
	}
	return land
}

func isLandNBT(allBoundingBoxes *[]map[string]float64, point []float64, polygons *Polygons) bool {
	for i, j := range *allBoundingBoxes {
		if boundingContains(&j, []float64{point[0], point[1]}) {
			if polygonContains(&(*polygons)[i].Points, []float64{point[0], point[1]}) {
				return true
			}
		}
	}
	return false
}
