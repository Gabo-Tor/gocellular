package main

/*
3d cellular automata

Rules:

survival / spawn / states / neighbour
ex: 4 / 4 / 5 / M

*/

import (
	"fmt"
	"math/rand"
	"strings"
)

const SIZE = 5

//                     0  1  2  3  4  5  6  7  8  9  10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26
var Survival = [27]uint8{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

//                  0  1  2  3  4  5  6  7  8  9  10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26
var Spawn = [27]uint8{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

const States = 5
const Neighbour = 1

func populate(board [][][]uint8) {

	for i, c := range board {
		for j, r := range c {
			for k := range r {
				board[i][j][k] = uint8(rand.Intn(5)) //TODO: magic number
			}
		}
	}
}

func printBoard(board [][][]uint8) {

	fmt.Println(strings.Repeat("-", (SIZE+1)*SIZE))
	for i, c := range board {
		for j, r := range c {
			for k := range r {
				fmt.Print(board[i][j][k])
			}
			fmt.Print("|")
		}
		fmt.Println()

	}
	fmt.Println(strings.Repeat("-", (SIZE+1)*SIZE))

}

func make3D(n uint8) [][][]uint8 {

	buf := make([]uint8, n*n*n) // uint8eger exponentiantion in go uint8(math.Pow(float64(n), float64(m))?????? :(
	x := make([][][]uint8, n)
	for i := range x {
		x[i] = make([][]uint8, n)
		for j := range x[i] {
			x[i][j] = buf[:n:n]
			buf = buf[n:]
		}
	}
	return x
}

func update(board [][][]uint8) {

	//oldBoard := board

	for i, c := range board {
		for j, r := range c {
			for k := range r {
				board[i][j][k] = States //TODO: magic number
			}
		}
	}
}

/*
func neighbours(board [][][]uint8, x uint8, y uint8, z uint8) uint8 {
	// Count neigbours

	// uint8(board[i][j][k] != 0) // How cart bool to int in golang???
}
*/
func main() {

	board := make3D(SIZE)
	populate(board)
	printBoard(board)

}
