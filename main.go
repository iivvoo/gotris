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

// Block is a possible tetris figure
type Block struct {
	cells []string
	color raylib.Color
}

var square = Block{cells: []string{
	"XX",
	"XX",
}, color: raylib.Red}
var line = Block{cells: []string{
	"XXXX",
}, color: raylib.Green}
var leftL = Block{cells: []string{
	"X  ",
	"XXX",
}, color: raylib.Blue}
var rightL = Block{cells: []string{
	"  X",
	"XXX",
}, color: raylib.Orange}
var z = Block{cells: []string{
	"XX ",
	" XX",
}, color: raylib.Yellow}
var reverseZ = Block{cells: []string{
	" XX",
	"XX ",
}, color: raylib.Gray}
var triangle = Block{cells: []string{
	" X ",
	"XXX",
}, color: raylib.Purple}

var blocks = []Block{square, line, leftL, rightL, z, reverseZ, triangle}

// ActiveBlock is the block currently played
type ActiveBlock struct {
	block    *Block
	row      int
	col      int
	rotation int // 1 2 3 4
}

func (a *ActiveBlock) random() {
	a.block = &blocks[rand.Intn(len(blocks))]
}

func (a *ActiveBlock) getCell(row int, col int, dRot int) bool {
	// returns the state of the cell at (row, col), taking
	// rotation + change into account into account (eventually)
	var cRows, cCols = a.getDims(dRot)
	var mRow, mCol = row, col
	var rotation = (a.rotation + dRot) % 4

	cRows--
	cCols--

	switch rotation {
	case 1:
		mRow = cCols - col
		mCol = row
	case 2:
		mRow = cRows - row
		mCol = cCols - col
	case 3:
		mRow = col
		mCol = cRows - row
	}

	return a.block.cells[mRow][mCol] == 'X'
}

func (a *ActiveBlock) getDims(dRot int) (int, int) {
	var rotation = (a.rotation + dRot) % 4

	if rotation%2 == 0 {
		return len(a.block.cells), len(a.block.cells[0])
	}
	return len(a.block.cells[0]), len(a.block.cells)
}

// Cell in the board
type Cell struct {
	used  bool
	color raylib.Color
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
	fps           int
	linesPs       int
	keysPs        int
	rows          int
	cols          int
	framesCounter int
	allowMove     bool
	active        *ActiveBlock
	board         [][]*Cell
}

func (g *Game) init(rows int, cols int, fps int, linesPs int, keysPs int) {
	g.framesCounter = 0
	g.allowMove = true
	g.rows = rows
	g.cols = cols
	g.fps = fps
	g.linesPs = linesPs
	g.keysPs = keysPs

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

func (g *Game) canMoveBlock(dRow int, dCol int, dRot int) bool {
	var a = g.active

	var cRows, cCols = a.getDims(dRot)

	defer g.setBlock(true)
	g.setBlock(false)

	for r := 0; r < cRows; r++ {
		for c := 0; c < cCols; c++ {
			var newRow = a.row + r + dRow
			var newCol = a.col + c + dCol

			if newRow >= g.rows || newRow < 0 || newCol >= g.cols || newCol < 0 {
				return false
			}
			if a.getCell(r, c, dRot) && g.board[newRow][newCol].used {
				return false
			}
		}
	}
	return true
}

func (g *Game) setBlock(state bool) {
	// TODO: rotate
	var a = g.active
	var b = a.block
	var cRows, cCols = a.getDims(0)

	for r := 0; r < cRows; r++ {
		for c := 0; c < cCols; c++ {
			if a.getCell(r, c, 0) {
				// only if not already set, else block cannot be placed
				g.board[a.row+r][a.col+c].used = state
				g.board[a.row+r][a.col+c].color = b.color

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
	g.board[g.rows-2][1].used = true
	g.board[g.rows-2][2].used = true
	g.board[g.rows-1][1].used = true
	g.board[g.rows-1][2].used = true
}

func (g *Game) blockDown() bool {
	// move active block down 1 row every N frames
	if g.framesCounter%(g.fps/g.linesPs) == 0 {
		if g.canMoveBlock(1, 0, 0) {
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
	var dCol, dRow, dRot = 0, 0, 0
	var changed = false

	g.framesCounter++

	// allow move every 5 frame updates. Rather arbitrary,
	// should probably relate to framerate. In this case, 60/5 times per second
	if g.framesCounter%(g.fps/g.keysPs) == 0 {
		g.allowMove = true
	}

	if !g.allowMove {
		return false
	}

	if raylib.IsKeyDown(raylib.KeyRight) {
		dCol = 1
		changed = true
	}
	if raylib.IsKeyDown(raylib.KeyLeft) {
		dCol = -1
		changed = true
	}
	if raylib.IsKeyDown(raylib.KeyUp) {
		dRot = 1
		changed = true
	}
	if raylib.IsKeyDown(raylib.KeyDown) {
		dRow = 1
		changed = true
	}
	if changed && g.canMoveBlock(dRow, dCol, dRot) {
		g.hideBlock()
		active.row += dRow
		active.col += dCol
		active.rotation += dRot
		if dRot > 0 {
			fmt.Println("Rotation changed")
		}
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
				raylib.DrawRectangle(int32(col*40), int32(row*40), 40, 40, g.board[row][col].color)
			} else {
				raylib.DrawRectangle(int32(col*40), int32(row*40), 40, 40, raylib.White)
			}
		}
	}
}

func main() {
	const fps = 60
	const linesPs = 2
	const keysPs = 4
	const cols = 12
	const rows = 20

	rand.Seed(time.Now().UnixNano())

	raylib.InitWindow(480, 800, "Ivo's GO Tetris")

	raylib.SetTargetFPS(int32(fps))

	board := Game{}
	board.init(rows, cols, fps, linesPs, keysPs)

	ab := &ActiveBlock{block: &z}
	ab.random()
	ab.row = 4
	ab.col = 4
	board.putBlock(ab)
	// board.putSomeBlocks()
	// board.print()

	for !raylib.WindowShouldClose() {
		// fmt.Printf("loop..")
		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		board.draw()
		board.input()
		if !board.blockDown() {
			ab := &ActiveBlock{block: &z}
			ab.random()
			ab.row = 0
			ab.col = 5
			board.putBlock(ab)
		}
		raylib.EndDrawing()
	}

	raylib.CloseWindow()
}
