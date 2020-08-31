package dataprocessing

// Polygon structure
type Polygon struct {
	Points [][]float64
	//LngTNorth []float64
	LngTNext []float64
	EoWNext  []int
	BtoX     []int
}

// Polygons structure
type Polygons []Polygon

func (p Polygons) Len() int {
	return len(p)
}

func (p Polygons) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p Polygons) Less(i, j int) bool {
	return len(p[i].Points) > len(p[j].Points)
}
