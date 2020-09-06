package algorithms

import (
    "../grids"
)

func neighboursBg(idx int, bg *grids.BasicGrid) []int {
    var neighbours []int
    var result []int

    neighbours = append(neighbours, idx-bg.XSize-1) // top left
    neighbours = append(neighbours, idx-bg.XSize)   // top
    neighbours = append(neighbours, idx-bg.XSize+1) // top right
    neighbours = append(neighbours, idx-1)          // left
    neighbours = append(neighbours, idx+1)          // right
    neighbours = append(neighbours, idx+bg.XSize-1) // bottom left
    neighbours = append(neighbours, idx+bg.XSize)   // bottom
    neighbours = append(neighbours, idx+bg.XSize+1) // bottom right

    for _, j := range neighbours {
        if j >= 0 && j < len(bg.VertexData) {
            if !bg.VertexData[j] {
                result = append(result, j)
            }
        }
    }
    return result
}

func extractRoute(prev *[]int, end int, bg *grids.BasicGrid) [][][]float64 {
    var route = make([][][]float64, 1)
    for {
        route[0] = append(route[0], bg.GridToCoord(bg.IDToGrid(end)))
        if (*prev)[end] == -1 {
            break
        }
        end = (*prev)[end]
    }
    return route
}

func extractNodes(nodesProcessed *[]int, bg *grids.BasicGrid) [][]float64 {
    var nodesExtended [][]float64
    for _, node := range *nodesProcessed {
        x := bg.IDToGrid(node)
        coord := bg.GridToCoord([]int{x[0], x[1]})
        nodesExtended = append(nodesExtended, coord)
    }
    return nodesExtended
}

// ExtractRouteBi ...
func ExtractRouteBi(prev *[][]int, meeting int, bg *grids.BasicGrid) [][][]float64 {
    var routes = make([][][]float64, 2)

    var secondMeeting = meeting
    for {
        routes[0] = append(routes[0], bg.GridToCoord(bg.IDToGrid(meeting)))
        if (*prev)[meeting][0] == -1 {
            break
        }
        meeting = (*prev)[meeting][0]
    }

    meeting = secondMeeting
    for {
        routes[1] = append(routes[1], bg.GridToCoord(bg.IDToGrid(meeting)))

        if (*prev)[meeting][1] == -1 {

            break
        }
        meeting = (*prev)[meeting][1]
    }
    return routes
}

func hBg(dir, node, from, to int, bg *grids.BasicGrid) float64 {
    if dir == 0 {
        return 0.5 * (distance(bg.GridToCoord(bg.IDToGrid(node)), bg.GridToCoord(bg.IDToGrid(to))) - distance(bg.GridToCoord(bg.IDToGrid(node)), bg.GridToCoord(bg.IDToGrid(from))))
    }
    return 0.5 * (distance(bg.GridToCoord(bg.IDToGrid(node)), bg.GridToCoord(bg.IDToGrid(from))) - distance(bg.GridToCoord(bg.IDToGrid(node)), bg.GridToCoord(bg.IDToGrid(to))))

}
