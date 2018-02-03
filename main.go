package main

import (
	"fmt"
	"math/rand"

	"github.com/gen2brain/raylib-go/raylib"
)

/*
 * Board consists of Cells
 * Blocks define figures, but there's also an active block.
 * This can be the actual figure, just reset it afterwards.
 * Once a block is fixed, it becomes fixed lines on the board
 * full lines are marked as full which allows animation of their
 * removal
 */

// Cell in the board
type Cell struct {
	used bool
	// should probably be color, possibly connectleft, connecttop
}

func (c *Cell) print() {
	if c.used {
		fmt.Print("X")
	} else {
		fmt.Print(" ")
	}
}

// Game of cells
type Game struct {
	rows  int
	cols  int
	board [][]*Cell
}

func (g *Game) init(rows int, cols int) {
	g.rows = rows
	g.cols = cols

	g.board = make([][]*Cell, rows) // rows + 1?

	for i := 0; i < len(g.board); i++ {
		g.board[i] = make([]*Cell, cols) // cols + 1?
		for j := 0; j < cols; j++ {
			g.board[i][j] = &Cell{}
		}
	}
}

func (g *Game) print() {
	for i := 0; i < len(g.board); i++ {
		for j := 0; j < len(g.board[i]); j++ {
			g.board[i][j].print()
		}
		fmt.Println()
	}
}

func (g *Game) putSomeBlocks() {
	g.board[1][1].used = true
	g.board[1][2].used = true
	g.board[2][1].used = true
	g.board[2][2].used = true
}

var block = []string{
	"XX",
	"XX",
}

func drawBlock(x int, y int) {
	// Draw a single block. Needs color
	defer fmt.Println("Block drawn")

	raylib.DrawRectangle(int32(x), int32(y), 20, 20, raylib.Blue)
}

func main() {
	raylib.InitWindow(450, 800, "Ivo's GO Tetris")

	raylib.SetTargetFPS(1)
	fmt.Println("My favorite number is", rand.Intn(10))

	const cols = 10
	const rows = 20

	board := Game{}
	board.init(rows, cols)

	board.putSomeBlocks()
	board.print()

	for !raylib.WindowShouldClose() {
		// fmt.Printf("loop..")
		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		raylib.DrawText("Congrats! You created your first window!", 10, 200, 20, raylib.LightGray)

		drawBlock(rand.Intn(20), rand.Intn(20))
		raylib.EndDrawing()
	}

	raylib.CloseWindow()
}
