package dataprocessing

import (
	"math"
)

func transformLon(newNorth, point []float64) float64 {
	var transformedLon float64

	// Degrees to Radians
	var dtr = math.Pi / 180.0

	// New north is already the north pole
	if newNorth[1] == 90.0 {
		transformedLon = point[0]
	} else {
		var t = math.Sin((point[0]-newNorth[0])*dtr) * math.Cos(point[1]*dtr)
		var b = math.Sin(dtr*point[1])*math.Cos(newNorth[1]*dtr) - math.Cos(point[1]*dtr)*math.Sin(newNorth[1]*dtr)*math.Cos((point[0]-newNorth[0])*dtr)
		// Radians to Degrees
		transformedLon = math.Atan2(t, b) / dtr
	}
	return transformedLon
}

// Direction of the shortest path from a to b
// 1 = east, -1 = west, 0 = neither
func eastOrWest(aLon, bLon float64) int {
	var out int
	var del = bLon - aLon
	if del > 180.0 {
		del = del - 360.0
	}
	if del < -180.0 {
		del = del + 360.0
	}
	if del > 0.0 && del != 180.0 {
		out = -1
	} else if del < 0.0 && del != -180.0 {
		out = 1
	} else {
		out = 0
	}
	return out
}

// Test if point is inside polygon, use north-pole for the known-to-be-outside point
func pointInPolygonSphere(poly *Polygon, point []float64) bool {
	var inside = false
	var strike = false

	// Point is the north-pole
	if point[1] == 90.0 {
		return false
	}

	for i := 0; i < len(poly.Points); i++ {
		var a = (poly.Points)[i]
		var b = (poly.Points)[(i+1)%len(poly.Points)]
		strike = false

		var nortPole = []float64{0.0, 90.0}
		var aT = a[0]
		var pT = point[0]
		var bT = b[0]
		if a[0] == b[0] {
			// a[0] -= 0.000000001
			nortPole = []float64{0.1, 89.9}
			aT = transformLon(nortPole, a)
			pT = transformLon(nortPole, point)
			bT = transformLon(nortPole, b)
		}

		if pT == aT {
			strike = true
		} else {

			var aToP = eastOrWest(aT, pT)
			var pToB = eastOrWest(pT, bT)

			if aToP == (poly.EoWNext)[i] && pToB == (poly.EoWNext)[i] {
				strike = true
			}
		}

		if strike {
			if point[1] == a[1] && point[0] == a[0] {
				return true
			}

			var pLonTransformed = transformLon(a, point)

			if (poly.LngTNext)[i] == pLonTransformed {
				return true
			}

			var bToP = eastOrWest((poly.LngTNext)[i], pLonTransformed)
			if (poly.BtoX)[i] == -bToP {
				inside = !inside
			}
		}
	}
	return inside
}

func isLandSphere(tree *boundingTree, point []float64, polygons *Polygons) bool {
	land := false
	if point[1] <= -80.0 {
		return true
	}
	if boundingContains(&tree.boundingBox, point) {
		if (*tree).id >= 0 {
			land = pointInPolygonSphere(&(*polygons)[(*tree).id], point)
			if land {
				return land
			}
		}
		for _, child := range (*tree).children {
			land = isLandSphere(&child, point, polygons)
			if land {
				return land
			}
		}
	}
	return land
}

func isLandSphereNBT(allBoundingBoxes *[]map[string]float64, point []float64, polygons *Polygons) bool {
	for i, j := range *allBoundingBoxes {
		if boundingContains(&j, point) {
			if pointInPolygonSphere(&(*polygons)[i], point) {
				return true
			}
		}
	}
	return false
}
