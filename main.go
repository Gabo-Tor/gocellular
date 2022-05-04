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

const SIZE = 30                        // >50 is too much for the 3d engine
const FRECUENCY = 5                    //hz
const INITIAL_ALIVE_PROBABILITY = 0.02 //0 - 1

// 3D automata rules:
//                       0  1  2  3  4  5  6  7  8  9  10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26
var survival = [27]uint8{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

//                    0  1  2  3  4  5  6  7  8  9  10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26
var spawn = [27]uint8{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

const states = 5

const Neighbour = 1 //1=Moore 0=Von Newmann

func populate(board [][][]uint8) {
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

func printBoard(board [][][]uint8) {
	// Print board to the std output... for debbuging, replaced by game engine
	fmt.Println(strings.Repeat("-", (SIZE+1)*SIZE))
	for i, c := range board {
		for j, r := range c {
			for k := range r {
				switch board[i][j][k] { //TODO:does this work with less than 3 states?? (nope, commented out)
				case 0:
					fmt.Print(" ")
					// fmt.Print(board[i][j][k])
				/*case 1:
					fmt.Print("░")
				case states - 2:
					fmt.Print("▓")
				case states - 1:
					fmt.Print("█")*/
				default:
					fmt.Print("▒")
				}
			}
			fmt.Print("|")
		}
		fmt.Println()
	}
	fmt.Println(strings.Repeat("-", (SIZE+1)*SIZE))

}

func create3dSlice(n int) [][][]uint8 {
	// creates a dyanmical square 3d slice
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

func one_if_positive(value uint8) int { // Does the compile make this inline? (yes) can I ask for this explicitly? (??)
	if value > 0 {
		return 1
	}
	return 0
}

func count_neigbours(board [][][]uint8, x int, y int, z int) int {
	//counts neigbous with either Moore o Von Neymann vecinty
	//This function is slooow, slices are slowww :(
	count := 0
	if Neighbour == 1 {

		count += one_if_positive(board[x-1][y][z])
		count += one_if_positive(board[x-1][y-1][z])
		count += one_if_positive(board[x-1][y+1][z])

		count += one_if_positive(board[x-1][y-1][z-1])
		count += one_if_positive(board[x-1][y][z-1])
		count += one_if_positive(board[x-1][y+1][z-1])

		count += one_if_positive(board[x-1][y-1][z+1])
		count += one_if_positive(board[x-1][y][z+1])
		count += one_if_positive(board[x-1][y+1][z+1])

		count += one_if_positive(board[x+1][y][z])
		count += one_if_positive(board[x+1][y-1][z])
		count += one_if_positive(board[x+1][y+1][z])

		count += one_if_positive(board[x+1][y-1][z-1])
		count += one_if_positive(board[x+1][y][z-1])
		count += one_if_positive(board[x+1][y+1][z-1])

		count += one_if_positive(board[x+1][y-1][z+1])
		count += one_if_positive(board[x+1][y][z+1])
		count += one_if_positive(board[x+1][y+1][z+1])

		//  count += one_if_positive(board[x  ][y  ][z  ]) leaving this here for completeness
		count += one_if_positive(board[x][y-1][z])
		count += one_if_positive(board[x][y+1][z])

		count += one_if_positive(board[x][y-1][z-1])
		count += one_if_positive(board[x][y][z-1])
		count += one_if_positive(board[x][y+1][z-1])

		count += one_if_positive(board[x][y-1][z+1])
		count += one_if_positive(board[x][y][z+1])
		count += one_if_positive(board[x][y+1][z+1])

	} else {
		fmt.Print("von newman neighborhood not implemented yet!")
	}
	return count
}

func update(board [][][]uint8) {

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

	oldBoard := make([][][]uint8, len(board)) // Is this memory safe???
	for i, c := range board {
		oldBoard[i] = make([][]uint8, len(board[i]))
		for j := range c {
			oldBoard[i][j] = make([]uint8, len(board[i][j]))
			copy(oldBoard[i][j], board[i][j])
		}
	}

	for i, c := range board {
		if i > 0 && i < SIZE-1 { //TODO: this is ugly
			for j, r := range c {
				if j > 0 && j < SIZE-1 {
					for k := range r {
						if k > 0 && k < SIZE-1 {
							if board[i][j][k] == 1 { //cell is alive but on its last state
								board[i][j][k] = survival[count_neigbours(oldBoard, i, j, k)]

							} else if board[i][j][k] > 1 { //cell is alive
								board[i][j][k]--

							} else { //cell is dead
								board[i][j][k] = (states - 1) * spawn[count_neigbours(oldBoard, i, j, k)]
							}

						}
					}
				}
			}
		}
	}
}

// Fuctions for 3D
func display_board(gBoard [][][]*graphic.Mesh, board [][][]uint8) {
	// Makes cells bigger or smaller depending on their state, dead cells have scale 0
	for i, c := range board {
		for j, r := range c {
			for k := range r {
				scale := float32(board[i][j][k]) / float32(states-1)
				gBoard[i][j][k].SetScale(scale, scale, scale)

			}
		}
	}
}

func create_board(scene *core.Node) [][][]*graphic.Mesh {
	//Creates graphical objestes and stores them on a 3D slice
	n := SIZE
	buf := make([]*graphic.Mesh, n*n*n) // uint8eger exponentiantion in go: uint8(math.Pow(float64(n), float64(m)) isn there a cleaner way?????? :(
	x := make([][][]*graphic.Mesh, n)
	for i := range x {
		x[i] = make([][]*graphic.Mesh, n)
		for j := range x[i] {
			x[i][j] = buf[:n:n]
			buf = buf[n:]
			for k := 0; k < SIZE; k++ {
				x[i][j][k] = create_box(float32(i), float32(j), float32(k))
				c := math32.NewColor("white")
				c.B = float32(i) / float32(SIZE) // Nice color gradient
				c.G = float32(j) / float32(SIZE)
				c.R = float32(k) / float32(SIZE)
				x[i][j][k].SetMaterial(material.NewStandard(c)) //
				scene.Add(x[i][j][k])
			}
		}
	}
	return x
}

func create_box(x float32, y float32, z float32) *graphic.Mesh {
	geom := geometry.NewBox(1.0/SIZE, 1.0/SIZE, 1.0/SIZE)
	mat := material.NewStandard(math32.NewColor("Blue")) // Its probably a good idea to create the box with the final color
	mesh := graphic.NewMesh(geom, mat)
	mesh.SetPosition(x/SIZE-0.5, y/SIZE-0.5, z/SIZE-0.5)
	return mesh
}

func main() {

	// //--For profiling code--
	// file, _ := os.Create("./cpu.pprof")
	// pprof.StartCPUProfile(file)
	// defer pprof.StopCPUProfile()

	board := create3dSlice(SIZE)
	populate(board)

	ticker := time.NewTicker((1000 / FRECUENCY) * time.Millisecond)
	done := make(chan bool) //TODO: delete this

	go func() { //TODO: make this more compact
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				update(board)
				_ = t //TODO: delete t, we are probably never using done either
			}
		}
	}()

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

		populate(board)

	})
	scene.Add(btn)

	// Create and add lights to the scene
	scene.Add(light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.8))
	pointLight := light.NewPoint(&math32.Color{1, 1, 1}, 5.0)
	pointLight.SetPosition(1, 0, 2)
	scene.Add(pointLight)

	// Create and add an axis helper to the scene
	scene.Add(helper.NewAxes(0.5))

	// Set background color to gray
	a.Gls().ClearColor(0.4, 0.4, 0.4, 1.0)

	// Run the application
	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(scene, cam)
		display_board(graphical_board, board)
	})
}
