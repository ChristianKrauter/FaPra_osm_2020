package dataprocessing

import (
	"fmt"
	"math"
)

func transformLon(newNorth, point []float64) float64 {
	var transformedLon float64

	var dtr = math.Pi / 180.0

	// New north is already the north pole
	if newNorth[1] == 90.0 {
		transformedLon = point[0]
	} else {
		var t = math.Sin((point[0]-newNorth[0])*dtr) * math.Cos(point[1]*dtr)
		var b = math.Sin(dtr*point[1])*math.Cos(newNorth[1]*dtr) - math.Cos(point[1]*dtr)*math.Sin(newNorth[1]*dtr)*math.Cos((point[0]-newNorth[0])*dtr)
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
	if point[1] == -90 {
		fmt.Printf("Tried to check point antipodal to the north pole.")
		return true
	}

	// Point is the north-pole
	if point[1] == 90 {
		return false
	}

	for i := 0; i < len(*polygon); i++ {
		var a = (*polygon)[i]
		var b = (*polygon)[(i+1)%len(*polygon)]

		if point[0] == a[0] {
			strike = true
		} else {
			var aToB = eastOrWest(a[0], b[0])
			var aToP = eastOrWest(a[0], point[0])
			var pToB = eastOrWest(point[0], b[0])
			if aToP == aToB && pToB == aToB {
				strike = true
			}
		}

		if strike {
			if point[1] == a[1] && point[0] == a[0] {
				return true
			}

			// Possible to calculate once at polygon creation
			var northPoleLonTransformed = transformLon(a, []float64{0, 90})
			var bLonTransformed = transformLon(a, b)
			// Not possible
			var pLonTransformed = transformLon(a, point)

			if bLonTransformed == pLonTransformed {
				return true
			}

			var bToX = eastOrWest(b[0], northPoleLonTransformed)
			var bToP = eastOrWest(b[0], point[0])
			if bToX == -bToP {
				inside = !inside
			}
		}
	}

	return inside
}

func isLandSphere(tree *boundingTree, point []float64, allCoastlines *[][][]float64) bool {
	land := false
	if boundingContains(&tree.boundingBox, point) {
		if (*tree).id >= 0 {
			land = pointInPolygonSphere(&(*allCoastlines)[(*tree).id], point)
			if land {
				return land
			}
		}
		for _, child := range (*tree).children {
			land = isLandSphere(&child, point, allCoastlines)
			if land {
				return land
			}
		}
	}
	return land
}

func isLandSphereNBT(allBoundingBoxes *[]map[string]float64, point []float64, allCoastlines *[][][]float64) bool {
	for i, j := range *allBoundingBoxes {
		if boundingContains(&j, point) {
			if pointInPolygonSphere(&(*allCoastlines)[i], point) {
				return true
			}
		}
	}
	return false
}
