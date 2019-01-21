package grid

// import (

// )

type Grid [][]int

type Coor struct {
	X, Y int
}

func New(x, y int) *Grid {
	arr := make([][]int, y)
	for i := range arr {
		arr[i] = make([]int, x)
	}
	grid := Grid(arr)
	return &grid
}