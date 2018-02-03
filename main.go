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
 *
 * Blocks is pure definitie (met kleur),
 * ActiveBlock is block in game, met rotatie. Zou zelfs onder
 * Game kunnen hangen
 */

var block = []string{
	"XX",
	"XX",
}

// Block is a possible tetris figure
type Block struct {
	cells []string
	// color
}

var square = Block{cells: []string{"XX", "XX"}}
var line = Block{cells: []string{"XXXX"}}
var leftL = Block{cells: []string{"X", "XXXX"}}
var rightL = Block{cells: []string{"  X", "XXX"}}
var triangle = Block{cells: []string{" X", "XXX"}}

// ActiveBlock is the block currently played
type ActiveBlock struct {
	block    *Block
	row      int
	col      int
	rotation int // 1 2 3 4
}

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
	rows   int
	cols   int
	active *ActiveBlock
	board  [][]*Cell
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

func (g *Game) putBlock(block *ActiveBlock) {
	g.active = block
	g.insertBlock()
}

func (g *Game) canMoveBlock(dRow int, dCol int) bool {
	return true

	// insert the active block. May fail, handle that later
	// TODO: rotat
	var a = g.active
	var b = a.block

	fmt.Printf("%d %d", a.row, a.col)

	for r := 0; r < len(b.cells); r++ {
		for c := 0; c < len(b.cells[r]); c++ {
			if b.cells[r][c] == 'X' {
				// only if not already set, else block cannot be placed
				g.board[a.row+r][a.col+c].used = true
			}
		}
	}
	return false
}

func (g *Game) setBlock(state bool) {
	// insert the active block. May fail, handle that later
	// TODO: rotat
	var a = g.active
	var b = a.block

	fmt.Printf("%d %d", a.row, a.col)

	for r := 0; r < len(b.cells); r++ {
		for c := 0; c < len(b.cells[r]); c++ {
			if b.cells[r][c] == 'X' {
				// only if not already set, else block cannot be placed
				g.board[a.row+r][a.col+c].used = state
			}
		}
	}
}

func (g *Game) insertBlock() {
	g.setBlock(true)
}

func (g *Game) clearBlock() {
	g.setBlock(false)
}

func (g *Game) putSomeBlocks() {
	g.board[10][1].used = true
	g.board[10][2].used = true
	g.board[11][1].used = true
	g.board[11][2].used = true
}

func (g *Game) input() {
	var active = g.active
	var dCol, dRow = 0, 0

	if raylib.IsKeyDown(raylib.KeyRight) {
		fmt.Println("Right")
		dCol = 1
	}
	if raylib.IsKeyDown(raylib.KeyLeft) {
		fmt.Println("Left")
		dCol = -1
	}
	if raylib.IsKeyDown(raylib.KeyUp) {
		fmt.Println("Up")
		dRow = -1
	}
	if raylib.IsKeyDown(raylib.KeyDown) {
		fmt.Println("Down")
		dRow = 1
	}
	if g.canMoveBlock(dRow, dCol) {
		g.clearBlock()
		active.row += dRow
		active.col += dCol
		g.insertBlock()

	}
}
func (g *Game) draw() {
	// Draw a single block. Needs color
	// defer fmt.Println("Block drawn")

	for row := 0; row < len(g.board); row++ {
		for col := 0; col < len(g.board[row]); col++ {
			if g.board[row][col].used {
				raylib.DrawRectangle(int32(col*20), int32(row*20), 20, 20, raylib.Blue)
			} else {
				raylib.DrawRectangle(int32(col*20), int32(row*20), 20, 20, raylib.White)
			}
		}
	}
}

func main() {
	raylib.InitWindow(450, 800, "Ivo's GO Tetris")

	raylib.SetTargetFPS(6)
	fmt.Println("My favorite number is", rand.Intn(10))

	const cols = 8
	const rows = 12

	board := Game{}
	board.init(rows, cols)

	ab := &ActiveBlock{block: &triangle}
	ab.row = 4
	ab.col = 4
	board.putBlock(ab)
	board.putSomeBlocks()
	board.print()

	for !raylib.WindowShouldClose() {
		// fmt.Printf("loop..")
		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		// raylib.DrawText("Congrats! You created your first window!", 10, 200, 20, raylib.LightGray)

		// drawBlock(rand.Intn(20), rand.Intn(20))
		board.draw()
		board.input()
		raylib.EndDrawing()
	}

	raylib.CloseWindow()
}
