package dataprocessing

import (
	"fmt"
	"math"
)

func transformLon(newNorth, point []float64) float64 {
	var transformedLon float64

	var dtr = math.Pi / 180

	// New north is already the north pole
	if newNorth[0] == 90.0 {
		transformedLon = point[0]
	} else {
		var t = math.Sin((point[1]-newNorth[1])*dtr) * math.Cos(point[0]*dtr)
		var b = math.Sin(dtr*point[0])*math.Cos(newNorth[0]*dtr) - math.Cos(point[0]*dtr)*math.Sin(newNorth[0]*dtr)*math.Cos((point[1]-newNorth[1])*dtr)
		transformedLon = math.Atan2(t, b) / dtr
	}

	return transformedLon
}

// Direction of the shortest path from a to b
// 1 = east, -1 = west, 0 = neither
func eastOrWest(aLon, bLon float64) int {
	var out int
	var del = bLon - aLon
	if del > 180 {
		del = del - 360
	}
	if del < -180 {
		del = del + 360
	}
	if del > 0 && del != 180 {
		out = -1
	} else if del < 0 && del != -180 {
		out = 1
	} else {
		out = 0
	}
	return out
}

// Test if point is inside polygon, use north-pole for the known-to-be-outside point
func pointInPolygonSphere(polygon *[][]float64, point []float64) bool {
	var inside = false
	var strike bool

	// Point is the south-pole
	// Pontentially antipodal check
	if point[0] == -90 {
		fmt.Printf("Tried to check point antipodal to the north pole.")
		return true
	}

	// Point is the north-pole
	if point[0] == 90 {
		return false
	}

	for i := 0; i < len(*polygon); i++ {
		var a = (*polygon)[i]
		var b = (*polygon)[i+1%len(*polygon)]

		if point[1] == a[1] {
			strike = true
		} else {
			var aToB = eastOrWest(a[1], b[1])
			var aToP = eastOrWest(a[1], point[1])
			var pToB = eastOrWest(point[1], b[1])
			if aToP == aToB && pToB == aToB {
				strike = true
			}
		}

		if strike {
			if point[0] == a[0] && point[1] == a[1] {
				return true
			}

			// Possible to calculate once at polygon creation
			var northPoleLonTransformed = transformLon(a, []float64{90, 0})
			var bLonTransformed = transformLon(a, []float64{90, 0})
			// Not possible
			var pLonTransformed = transformLon(a, []float64{90, 0})

			if bLonTransformed == pLonTransformed {
				return true
			}

			var bToX = eastOrWest(b[1], northPoleLonTransformed)
			var bToP = eastOrWest(b[1], point[1])
			if bToX == -bToP {
				inside = !inside
			}
		}
	}

	return inside
}

func rayCastSphere(point, s, e []float64) (bool, bool) {
	if s[0] > e[0] {
		s, e = e, s
	}
	return true, true
}

func polygonContainsSphere(polygon *[][]float64, point []float64) bool {
	b, on := rayCastSphere(point, (*polygon)[0], (*polygon)[len(*polygon)-1])
	if on {
		return true
	}

	for i := 0; i < len(*polygon)-1; i++ {
		inter, on := rayCastSphere(point, (*polygon)[i], (*polygon)[i+1])
		if on {
			return true
		}
		if inter {
			b = !b
		}
	}
	return b
}

func isLandSphere(tree *boundingTree, point []float64, allCoastlines *[][][]float64) bool {
	land := false
	if boundingContains(&tree.boundingBox, point) {
		for _, child := range (*tree).children {
			land = isLandSphere(&child, point, allCoastlines)
			if land {
				return land
			}
		}
		if (*tree).id >= 0 {
			land = polygonContainsSphere(&(*allCoastlines)[(*tree).id], point)
		}
	}
	return land
}

func isLandSphereNBT(allBoundingBoxes *[]map[string]float64, point []float64, allCoastlines *[][][]float64) bool {
	for i, j := range *allBoundingBoxes {
		if boundingContains(&j, []float64{point[0], point[1]}) {
			if polygonContainsSphere(&(*allCoastlines)[i], []float64{point[0], point[1]}) {
				return true
			}
		}
	}
	return false
}
