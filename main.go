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
	"time"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/util/helper"
	"github.com/g3n/engine/window"
)

const SIZE = 20                       // >50 is too much for the 3d engine
const FRECUENCY = 5                   // hz
const INITIAL_ALIVE_PROBABILITY = 0.3 // 0 - 1

// 3D automata rules:
//                       0  1  2  3  4  5  6  7  8  9  10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26
var survival = [27]uint8{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

//                    0  1  2  3  4  5  6  7  8  9  10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26
var spawn = [27]uint8{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

const states = 5

const Neighbour = 1 //1=Moore 0=Von Newmann

func populate(board *[SIZE][SIZE][SIZE]uint8) {
	// Initialice board randomly with INITIAL_ALIVE_PROBABILITY
	for i, c := range board {
		for j, r := range c {
			for k := range r {
				if rand.Float32() > INITIAL_ALIVE_PROBABILITY {
					board[i][j][k] = 0
				} else {
					board[i][j][k] = uint8(rand.Intn(states-1) + 1)
				}
			}
		}
	}
}

func printBoard(board *[SIZE][SIZE][SIZE]uint8) {
	// Print board to the std output... for debbuging, replaced by game engine
	fmt.Println(strings.Repeat("-", (SIZE+1)*SIZE))
	for i, c := range board {
		for j, r := range c {
			for k := range r {
				if board[i][j][k] == 0 {
					fmt.Print(" ")
				} else {
					fmt.Print("█")
				}
			}
			fmt.Print("|")
		}
		fmt.Println()
	}
	fmt.Println(strings.Repeat("-", (SIZE+1)*SIZE))

}

func one_if_alive(value uint8) int { // Does the compile make this inline? (yes) can I ask for this explicitly? (??)
	if value == states-1 {
		return 1
	}
	return 0
}

func count_neigbours(board *[SIZE][SIZE][SIZE]uint8, x int, y int, z int) int {
	//counts neigbous with either Moore o Von Neumann vecinty
	//This function is slooow, slices are slowww :(
	count := 0
	if Neighbour == 1 {
		count = countNeighborsMoore(count, board, x, y, z)
	} else {
		count = countNeighborsVonNeumann(count, board, x, y, z)
	}
	return count
}

func countNeighborsMoore(count int, board *[SIZE][SIZE][SIZE]uint8, x int, y int, z int) int {
	count += one_if_alive(board[x-1][y][z])
	count += one_if_alive(board[x-1][y-1][z])
	count += one_if_alive(board[x-1][y+1][z])

	count += one_if_alive(board[x-1][y-1][z-1])
	count += one_if_alive(board[x-1][y][z-1])
	count += one_if_alive(board[x-1][y+1][z-1])

	count += one_if_alive(board[x-1][y-1][z+1])
	count += one_if_alive(board[x-1][y][z+1])
	count += one_if_alive(board[x-1][y+1][z+1])

	count += one_if_alive(board[x+1][y][z])
	count += one_if_alive(board[x+1][y-1][z])
	count += one_if_alive(board[x+1][y+1][z])

	count += one_if_alive(board[x+1][y-1][z-1])
	count += one_if_alive(board[x+1][y][z-1])
	count += one_if_alive(board[x+1][y+1][z-1])

	count += one_if_alive(board[x+1][y-1][z+1])
	count += one_if_alive(board[x+1][y][z+1])
	count += one_if_alive(board[x+1][y+1][z+1])

	//  count += one_if_positive(board[x  ][y  ][z  ]) leaving this here for completeness
	count += one_if_alive(board[x][y-1][z])
	count += one_if_alive(board[x][y+1][z])

	count += one_if_alive(board[x][y-1][z-1])
	count += one_if_alive(board[x][y][z-1])
	count += one_if_alive(board[x][y+1][z-1])

	count += one_if_alive(board[x][y-1][z+1])
	count += one_if_alive(board[x][y][z+1])
	count += one_if_alive(board[x][y+1][z+1])
	return count
}

func countNeighborsVonNeumann(count int, board *[SIZE][SIZE][SIZE]uint8, x int, y int, z int) int {
	count += one_if_alive(board[x-1][y][z])
	count += one_if_alive(board[x+1][y][z])

	count += one_if_alive(board[x][y-1][z])
	count += one_if_alive(board[x][y+1][z])

	count += one_if_alive(board[x][y][z-1])
	count += one_if_alive(board[x][y][z+1])

	return count
}

func update(board *[SIZE][SIZE][SIZE]uint8) {

	// Makes the map circular(tiled): Faces
	for i := 1; i < SIZE-1; i++ {
		for j := 1; j < SIZE-1; j++ {
			board[0][i][j] = board[SIZE-2][i][j]
			board[i][j][0] = board[i][j][SIZE-2]
			board[j][0][i] = board[j][SIZE-2][i]
			board[SIZE-1][i][j] = board[1][i][j]
			board[i][j][SIZE-1] = board[i][j][1]
			board[j][SIZE-1][i] = board[j][1][i]
		}
	}
	// Makes the map circular(tiled): Corners
	for i := 1; i < SIZE-1; i++ {
		board[0][i][0] = board[SIZE-2][i][SIZE-2]
		board[i][0][0] = board[i][SIZE-2][SIZE-2]
		board[0][0][i] = board[SIZE-2][SIZE-2][i]

		board[SIZE-1][i][SIZE-1] = board[1][i][1]
		board[i][SIZE-1][SIZE-1] = board[i][1][1]
		board[SIZE-1][SIZE-1][i] = board[1][1][i]

		board[SIZE-1][i][0] = board[1][i][SIZE-2]
		board[i][0][SIZE-1] = board[i][SIZE-2][1]
		board[0][SIZE-1][i] = board[SIZE-2][1][i]

		board[0][i][SIZE-1] = board[SIZE-2][i][1]
		board[i][SIZE-1][0] = board[i][1][SIZE-2]
		board[SIZE-1][0][i] = board[1][SIZE-2][i]
	}
	// Makes the map circular(tiled): Vertices
	board[0][0][0] = board[SIZE-2][SIZE-2][SIZE-2]
	board[SIZE-1][0][0] = board[1][SIZE-2][SIZE-2]
	board[0][SIZE-1][0] = board[SIZE-2][1][SIZE-2]
	board[SIZE-1][SIZE-1][0] = board[1][1][SIZE-2]
	board[0][0][SIZE-1] = board[SIZE-2][SIZE-2][1]
	board[SIZE-1][0][SIZE-1] = board[1][SIZE-2][1]
	board[0][SIZE-1][SIZE-1] = board[SIZE-2][1][1]
	board[SIZE-1][SIZE-1][SIZE-1] = board[1][1][1]

	var oldBoard [SIZE][SIZE][SIZE]uint8 // Is this memory safe???
	for i, c := range board {
		for j := range c {
			oldBoard[i][j] = board[i][j]
		}
	}

	for i := 1; i < SIZE-1; i++ {
		go update_rows(board, i, oldBoard) // is multithreading this simple in go? should we wait for every thread to finish? (this is probably not a problem in practice)
	}
}

func update_rows(board *[SIZE][SIZE][SIZE]uint8, i int, oldBoard [SIZE][SIZE][SIZE]uint8) {
	for j := 1; j < SIZE-1; j++ {
		for k := 1; k < SIZE-1; k++ {
			if board[i][j][k] == states-1 { //cell is alive but on its first state
				board[i][j][k] -= (1 - survival[count_neigbours(&oldBoard, i, j, k)])

			} else if board[i][j][k] > 0 { //cell is alive
				board[i][j][k]--

			} else { //cell is dead
				board[i][j][k] = (states - 1) * spawn[count_neigbours(&oldBoard, i, j, k)]
			}
		}
	}
}

// Fuctions for 3D
func display_board(gBoard [SIZE][SIZE][SIZE]*graphic.Mesh, board [SIZE][SIZE][SIZE]uint8) {
	// Makes cells bigger or smaller depending on their state, dead cells have scale 0
	for i, c := range board {
		for j, r := range c {
			for k := range r {
				if board[i][j][k] == 0 {
					gBoard[i][j][k].SetVisible(false)
				} else {
					scale := float32(board[i][j][k]+2) / float32(states)
					gBoard[i][j][k].SetScale(scale, scale, scale)
					gBoard[i][j][k].SetVisible(true)
				}
			}
		}
	}
}

func create_board(scene *core.Node) [SIZE][SIZE][SIZE]*graphic.Mesh {
	// Creates graphical objestes and stores them on a 3D slice
	var x [SIZE][SIZE][SIZE]*graphic.Mesh
	for i := range x {
		for j := range x[i] {
			for k := 0; k < SIZE; k++ {
				color := math32.NewColor("white")
				color.B = float32(i) / float32(SIZE) // Nice color gradient
				color.G = float32(j) / float32(SIZE)
				color.R = float32(k) / float32(SIZE)
				x[i][j][k] = create_cell(float32(i), float32(j), float32(k), color)
				x[i][j][k].SetVisible(false)
				scene.Add(x[i][j][k])
			}
		}
	}
	return x
}

func create_cell(x float32, y float32, z float32, color *math32.Color) *graphic.Mesh {
	// Creates one box of desires color in position xyz
	geom := geometry.NewBox(1.0/SIZE, 1.0/SIZE, 1.0/SIZE)
	mat := material.NewStandard(color)
	mesh := graphic.NewMesh(geom, mat)
	mesh.SetPosition(x/SIZE-0.5, y/SIZE-0.5, z/SIZE-0.5)
	return mesh
}

func main() {

	rand.Seed(time.Now().UnixNano()) // seed

	// //--For profiling code--
	// file, _ := os.Create("./cpu.pprof")
	// pprof.StartCPUProfile(file)
	// defer pprof.StopCPUProfile()

	var board [SIZE][SIZE][SIZE]uint8
	populate(&board)

	// for n := 0; n < 9; n++ { // Simple game loop for debugging purpouses
	// 	// printBoard(board)
	// 	update(&board)
	// }

	// Create application and scene
	a := app.App()
	scene := core.NewNode()

	// Set the scene to be managed by the gui manager
	gui.Manager().Set(scene)

	// Create perspective camera
	cam := camera.New(1)
	cam.SetPosition(0, 0, 1.7)
	scene.Add(cam)

	// Set up orbit control for the camera
	camera.NewOrbitControl(cam)

	// Set up callback to update viewport and camera aspect ratio when the window is resized
	onResize := func(evname string, ev interface{}) {
		// Get framebuffer size and update viewport accordingly
		width, height := a.GetSize()
		a.Gls().Viewport(0, 0, int32(width), int32(height))
		// Update the camera's aspect ratio
		cam.SetAspect(float32(width) / float32(height))
	}
	a.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)

	graphical_board := create_board(scene)

	// Create and add a button to the scene
	btn := gui.NewButton("Reset")
	btn.SetPosition(60, 40)
	btn.SetSize(40, 40)
	btn.Subscribe(gui.OnClick, func(name string, ev interface{}) {

		populate(&board)

	})
	scene.Add(btn)

	// Create and add lights to the scene
	scene.Add(light.NewAmbient(&math32.Color{R: 1, G: 1, B: 1}, 0.8))
	pointLight := light.NewPoint(&math32.Color{R: 1, G: 1, B: 1}, 5.0)
	pointLight.SetPosition(1, 0, 2)
	scene.Add(pointLight)

	// Create and add an axis helper to the scene
	scene.Add(helper.NewAxes(0.05))

	// Set background color to gray
	a.Gls().ClearColor(0.4, 0.4, 0.4, 1.0)

	// Asynchronous map updating
	ticker := time.NewTicker((1000 / FRECUENCY) * time.Millisecond)
	go func() {
		for range ticker.C {
			update(&board)
		}
	}()

	// Run the application
	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(scene, cam)
		display_board(graphical_board, board)
	})
}
