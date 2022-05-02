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

const SIZE = 7

//                       0  1  2  3  4  5  6  7  8  9  10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26
var survival = [27]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}

//                    0  1  2  3  4  5  6  7  8  9  10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26
var spawn = [27]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0}

const states = 2
const Neighbour = 1

func populate(board [][][]uint8) {

	for i, c := range board {
		for j, r := range c {
			for k := range r {
				if rand.Float32() > 0.9 {
					board[i][j][k] = 0
				} else {
					board[i][j][k] = uint8(rand.Intn(states))
					// board[i][j][k] = 1 //TODO: delete debug
				}

			}
		}
	}
}

func printBoard(board [][][]uint8) {

	fmt.Println(strings.Repeat("-", (SIZE+1)*SIZE))
	for i, c := range board {
		for j, r := range c {
			for k := range r {
				switch board[i][j][k] { //TODO: this does not work with less than 3 states??
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
					// fmt.Print(" ", board[i][j][k])
				}
			}
			fmt.Print("|")
		}
		fmt.Println()

	}
	fmt.Println(strings.Repeat("-", (SIZE+1)*SIZE))

}

func make3D(n int) [][][]uint8 {

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

func one_if_positive(value uint8) int { // Does the compile make this inline? can I make it explicitly?
	if value > 0 {
		return 1
	}
	return 0
}

func count_neigbours(board [][][]uint8, x int, y int, z int) int {
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
		fmt.Print("moore neighborhood not implemented yet!")
	}
	// fmt.Println(count) //TODO: debug
	return count
}

func update(board [][][]uint8) {

	for i := 0; i < SIZE; i++ {
		for j := 0; j < SIZE; j++ {
			board[0][i][j] = board[SIZE-2][i][j]
			board[i][j][0] = board[i][j][SIZE-2]
			board[j][0][i] = board[j][SIZE-2][i]
			board[SIZE-1][i][j] = board[1][i][j]
			board[i][j][SIZE-1] = board[i][j][1]
			board[j][SIZE-1][i] = board[j][1][i]
		}
	}

	oldBoard := board

	for i, c := range board {
		if i > 0 && i < SIZE-1 { //TODO: this is ugly
			for j, r := range c {
				if j > 0 && j < SIZE-1 {
					for k := range r {
						if k > 0 && k < SIZE-1 {
							if board[i][j][k] > 0 { //cell is alive
								board[i][j][k]--
								board[i][j][k] += survival[count_neigbours(oldBoard, i, j, k)]
							} else {
								board[i][j][k] = (states - 1) * spawn[count_neigbours(oldBoard, i, j, k)]
							}

						}
					}
				}
			}
		}
	}
}

func main() {

	board := make3D(SIZE)
	populate(board)

	for n := 1; n < 3; n++ { //magic number
		printBoard(board)
		update(board)
	}

	// Create application and scene
	a := app.App()
	scene := core.NewNode()

	// Set the scene to be managed by the gui manager
	gui.Manager().Set(scene)

	// Create perspective camera
	cam := camera.New(1)
	cam.SetPosition(0, 0, 3)
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
	btn := gui.NewButton("Make Red")
	btn.SetPosition(100, 40)
	btn.SetSize(40, 40)
	btn.Subscribe(gui.OnClick, func(name string, ev interface{}) {

		display_board(graphical_board, board)
		update(board)

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
	a.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)

	// Run the application
	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(scene, cam)

	})

}

func display_board(gBoard [][][]*graphic.Mesh, board [][][]uint8) {
	for i, c := range board {
		for j, r := range c {
			for k := range r {
				switch board[i][j][k] {
				case 0:

					gBoard[i][j][k].SetScale(0.001, 0.001, 0.001)
				/*case 1:

					gBoard[i][j][k].SetScale(0.5, 0.5, 0.5)
				case states - 3:

					gBoard[i][j][k].SetScale(0.8, 0.8, 0.8)
				case states - 2:

					gBoard[i][j][k].SetScale(0.9, 0.9, 0.9)
				case states - 1:
					//gBoard[i][j][k].SetMaterial(material.NewStandard(math32.NewColor("DarkBlue")))
					gBoard[i][j][k].SetScale(1, 1, 1)*/
				default:

					gBoard[i][j][k].SetScale(0.99, 0.99, 0.99)
				}
			}
		}
	}
}

func create_board(scene *core.Node) [][][]*graphic.Mesh {
	n := SIZE
	buf := make([]*graphic.Mesh, n*n*n) // uint8eger exponentiantion in go uint8(math.Pow(float64(n), float64(m))?????? :(
	x := make([][][]*graphic.Mesh, n)
	for i := range x {
		x[i] = make([][]*graphic.Mesh, n)
		for j := range x[i] {
			x[i][j] = buf[:n:n]
			buf = buf[n:]
			for k := 0; k < SIZE; k++ {
				x[i][j][k] = create_box(float32(i), float32(j), float32(k))
				c := math32.NewColor("white")
				c.B = float32(i) / float32(SIZE)
				c.G = float32(j) / float32(SIZE)
				c.R = float32(k) / float32(SIZE)
				x[i][j][k].SetMaterial(material.NewStandard(c))
				scene.Add(x[i][j][k])
			}
		}
	}
	return x

}

func create_box(x float32, y float32, z float32) *graphic.Mesh {
	geom := geometry.NewBox(1.0/SIZE, 1.0/SIZE, 1.0/SIZE)
	mat := material.NewStandard(math32.NewColor("Blue"))
	mesh := graphic.NewMesh(geom, mat)
	mesh.SetPosition(x/SIZE-0.5, y/SIZE-0.5, z/SIZE-0.5)
	//mesh.SetMaterial(material.NewStandard(math32.NewColor("Green")))
	return mesh
}
