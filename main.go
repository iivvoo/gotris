package main

import (
	"fmt"
	"math/rand"
	"time"

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
var z = Block{cells: []string{"XX", " XX"}}
var reverseZ = Block{cells: []string{" XX", "XX"}}
var triangle = Block{cells: []string{" X", "XXX"}}

var blocks = []Block{square, line, leftL, rightL, z, reverseZ, triangle}

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
	rows          int
	cols          int
	framesCounter int
	allowMove     bool
	active        *ActiveBlock
	board         [][]*Cell
}

func (g *Game) init(rows int, cols int) {
	g.framesCounter = 0
	g.allowMove = true
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
	g.showBlock()
}

func (g *Game) canMoveBlock(dRow int, dCol int) bool {
	var a = g.active
	var b = a.block

	defer g.setBlock(true)
	g.setBlock(false)

	for r := 0; r < len(b.cells); r++ {
		for c := 0; c < len(b.cells[r]); c++ {
			var newRow = a.row + r + dRow
			var newCol = a.col + c + dCol

			if newRow >= g.rows || newRow < 0 || newCol >= g.cols || newCol < 0 {
				return false
			}
			if g.board[newRow][newCol].used == true {
				return false
			}
		}
	}
	return true
}

func (g *Game) setBlock(state bool) {
	// insert the active block. May fail, handle that later
	// TODO: rotat
	var a = g.active
	var b = a.block

	for r := 0; r < len(b.cells); r++ {
		for c := 0; c < len(b.cells[r]); c++ {
			if b.cells[r][c] == 'X' {
				// only if not already set, else block cannot be placed
				g.board[a.row+r][a.col+c].used = state
			}
		}
	}
}

func (g *Game) showBlock() {
	g.setBlock(true)
}

func (g *Game) hideBlock() {
	g.setBlock(false)
}

func (g *Game) putSomeBlocks() {
	g.board[18][1].used = true
	g.board[18][2].used = true
	g.board[19][1].used = true
	g.board[19][2].used = true
}

func (g *Game) blockDown() bool {
	// move active block down 1 row every N frames
	if g.framesCounter%60 == 0 {
		if g.canMoveBlock(1, 0) {
			g.hideBlock()
			g.active.row++
			g.showBlock()
			return true
		}
		return false
	}
	return true
}
func (g *Game) input() bool {
	var active = g.active
	var dCol, dRow = 0, 0

	g.framesCounter++

	// allow move every 5 frame updates. Rather arbitrary,
	// should probably relate to framerate. In this case, 60/5 times per second
	if g.framesCounter%5 == 0 {
		g.allowMove = true
	}

	if !g.allowMove {
		return false
	}

	if raylib.IsKeyDown(raylib.KeyRight) {
		dCol = 1
	}
	if raylib.IsKeyDown(raylib.KeyLeft) {
		dCol = -1
	}
	if raylib.IsKeyDown(raylib.KeyUp) {
		dRow = -1
	}
	if raylib.IsKeyDown(raylib.KeyDown) {
		dRow = 1
	}
	if g.canMoveBlock(dRow, dCol) {
		g.hideBlock()
		active.row += dRow
		active.col += dCol
		g.showBlock()
		g.allowMove = false
		return true
	}
	return false
}

func (g *Game) draw() {
	// Draw a single block. Needs color
	// defer fmt.Println("Block drawn")

	for row := 0; row < len(g.board); row++ {
		for col := 0; col < len(g.board[row]); col++ {
			if g.board[row][col].used {
				raylib.DrawRectangle(int32(col*40), int32(row*40), 40, 40, raylib.Blue)
			} else {
				raylib.DrawRectangle(int32(col*40), int32(row*40), 40, 40, raylib.White)
			}
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	raylib.InitWindow(480, 800, "Ivo's GO Tetris")

	raylib.SetTargetFPS(60)
	fmt.Println("My favorite number is", rand.Intn(10))

	const cols = 12
	const rows = 20

	board := Game{}
	board.init(rows, cols)

	ab := &ActiveBlock{block: &blocks[rand.Intn(len(blocks))]}
	ab.row = 4
	ab.col = 4
	board.putBlock(ab)
	board.putSomeBlocks()
	board.print()

	for !raylib.WindowShouldClose() {
		// fmt.Printf("loop..")
		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		board.draw()
		board.input()
		if !board.blockDown() {
			ab := &ActiveBlock{block: &blocks[rand.Intn(len(blocks))]}
			ab.row = 0
			ab.col = 5
			board.putBlock(ab)
		}
		raylib.EndDrawing()
	}

	raylib.CloseWindow()
}
