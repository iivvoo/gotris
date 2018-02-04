package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/gen2brain/raylib-go/raylib"
)

/*
 * TODO
 * Pause game
 * Restart after finished
 * "Animate" clearing of rows
 * Show next block
 * Increase linedrop speed when #lines/score increases
 * Projection where block will end when dropped
 * make Row a struct?
 *  -> keeps track if 'full', etc
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
	fcLastKey     int
	active        *ActiveBlock
	score         int
	lines         int
	full          []int
	board         [][]*Cell
}

func (g *Game) init(rows int, cols int, fps int, linesPs int, keysPs int) {
	g.framesCounter = 0
	g.fcLastKey = 0
	g.rows = rows
	g.cols = cols
	g.fps = fps
	g.linesPs = linesPs
	g.keysPs = keysPs

	g.score = 0
	g.lines = 0

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

func (g *Game) putBlock(block *ActiveBlock) bool {
	g.active = block
	return g.showBlock()
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

func (g *Game) setBlock(state bool) bool {
	// TODO: rotate
	var a = g.active
	var b = a.block
	var cRows, cCols = a.getDims(0)

	if state {
		for r := 0; r < cRows; r++ {
			for c := 0; c < cCols; c++ {
				if a.getCell(r, c, 0) && g.board[a.row+r][a.col+c].used {
					return false
				}
			}
		}
	}
	for r := 0; r < cRows; r++ {
		for c := 0; c < cCols; c++ {
			if a.getCell(r, c, 0) {
				// only if not already set, else block cannot be placed
				g.board[a.row+r][a.col+c].used = state
				g.board[a.row+r][a.col+c].color = b.color
			}
		}
	}
	return true
}

func (g *Game) showBlock() bool {
	return g.setBlock(true)
}

func (g *Game) hideBlock() {
	g.setBlock(false)
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

func (g *Game) checkFullRows() int {
	// this could be replaced entirely if Row would be a struct
	var full []int

	for i := 0; i < g.rows; i++ {
		var usedCount = 0
		for _, c := range g.board[i] {
			if c.used {
				usedCount++
			}
		}
		if usedCount == g.cols {
			full = append(full, i)
		}
	}
	g.full = full
	return len(full)
}

func (g *Game) clearFullRows() {
	for _, line := range g.full {
		for above := line - 1; above >= 0; above-- {
			for c := 0; c < g.cols; c++ {
				g.board[above+1][c] = g.board[above][c]
			}
		}
		for c := 0; c < g.cols; c++ {
			g.board[0][c] = &Cell{}
		}
	}
	// make sure g.board[0] is entirely cleared
}

func (g *Game) input() bool {
	var active = g.active
	var dCol, dRow, dRot = 0, 0, 0
	var changed = false

	g.framesCounter++

	// require a specific interval between keypresses
	if g.framesCounter-g.fcLastKey < (g.fps / g.keysPs) {
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
		g.showBlock()
		g.fcLastKey = g.framesCounter
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
	const keysPs = 10
	const cols = 10
	const rows = 10

	var gameOver = false

	rand.Seed(time.Now().UnixNano())

	raylib.InitWindow(600, 800, "Ivo's GO Tetris")

	raylib.SetTargetFPS(int32(fps))

	board := Game{}
	board.init(rows, cols, fps, linesPs, keysPs)

	ab := &ActiveBlock{block: &z}
	ab.random()
	ab.row = 4
	ab.col = 4
	board.putBlock(ab)
	// board.print()

	for !raylib.WindowShouldClose() {
		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.LightGray)
		raylib.DrawText("Next", 420, 100, 40, raylib.Black)
		raylib.DrawText("Lines", 420, 300, 40, raylib.Black)
		raylib.DrawText(strconv.Itoa(board.lines), 420, 350, 40, raylib.Black)
		raylib.DrawText("Score", 420, 500, 40, raylib.Black)
		raylib.DrawText(strconv.Itoa(board.score), 420, 550, 40, raylib.Black)
		if !gameOver {
			board.draw()
			board.input()
			if !board.blockDown() {
				ab := &ActiveBlock{block: &z}
				ab.random()
				ab.row = 0
				ab.col = 5
				fullLines := board.checkFullRows()
				board.lines += fullLines

				board.score += fullLines * fullLines * 10

				board.clearFullRows()
				if !board.putBlock(ab) {
					gameOver = true
					fmt.Println("Game Over?")
				}
			}
		} else {
			raylib.DrawText("Game Over!", 100, 400, 40, raylib.Black)
		}
		raylib.EndDrawing()
	}

	raylib.CloseWindow()
}
