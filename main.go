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
 * Restart after finished
 * Increase linedrop speed when #lines/score increases
 * Projection where block will end when dropped
 * make Row a struct?
 *  -> keeps track if 'full', etc
 * "Animate" clearing of rows
 * "Next block" incorrect at game over (or more specific,
 *  block that failed to place should be next)
 *  multi-gui? https://godoc.org/github.com/golang-ui/nuklear/nk
 * can go-routines/channels be used to implement/abstract the main loop?
 */

// ActiveBlock is the block currently played
type ActiveBlock struct {
	block    *Block
	row      int
	col      int
	rotation int // 1 2 3 4
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

func (a *ActiveBlock) setBlock(g *Game, state bool) bool {
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

func (a *ActiveBlock) showBlock(g *Game) bool {
	return a.setBlock(g, true)
}

func (a *ActiveBlock) hideBlock(g *Game) {
	a.setBlock(g, false)
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
	next          *Block
	score         int
	lines         int
	full          []int
	paused        bool
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
	g.paused = false

	g.score = 0
	g.lines = 0

	g.board = make([][]*Cell, rows)

	for i := range g.board {
		g.board[i] = make([]*Cell, cols)
		for j := range g.board[i] {
			g.board[i][j] = &Cell{}
		}
	}
}

func (g *Game) print() {
	for i := range g.board {
		for j := range g.board[i] {
			g.board[i][j].print()
		}
		fmt.Println()
	}
}

func (g *Game) putBlock(block *ActiveBlock) bool {
	g.active = block
	return g.active.showBlock(g)
}

func (g *Game) canMoveBlock(dRow int, dCol int, dRot int) bool {
	var a = g.active

	var cRows, cCols = a.getDims(dRot)

	defer a.setBlock(g, true) // show
	a.setBlock(g, false)      // hide

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

func (g *Game) blockDown(immediate bool) bool {
	// move active block down 1 row every N frames
	if immediate || g.framesCounter%(g.fps/g.linesPs) == 0 {
		if g.canMoveBlock(1, 0, 0) {
			g.active.hideBlock(g)
			g.active.row++
			g.active.showBlock(g)
			// try to draw a ghost version further down

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
}

func (g *Game) input() bool {
	// refactor this into handling input and returning action
	// to be taken, if any?
	// returns true if action requires immediate update
	var active = g.active
	var dCol, dRow, dRot = 0, 0, 0
	var changed = false
	var immediate = false

	g.framesCounter++

	// require a specific interval between keypresses
	if g.framesCounter-g.fcLastKey < (g.fps / g.keysPs) {
		return false
	}
	if raylib.IsKeyDown(raylib.KeyRight) {
		dCol = 1
		changed = true
	}
	if raylib.IsKeyDown(raylib.KeySpace) {
		// test how far down we can go
		for g.canMoveBlock(dRow+1, dCol, dRot) {
			dRow++
		}
		changed = true
		immediate = true
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
		immediate = true
	}
	if raylib.IsKeyDown(raylib.KeyP) {
		g.paused = !g.paused
		g.fcLastKey = g.framesCounter
	}
	if !g.paused && changed && g.canMoveBlock(dRow, dCol, dRot) {
		// dropping means adding rows until we no longer can

		g.active.hideBlock(g)
		active.row += dRow
		active.col += dCol
		active.rotation += dRot
		g.active.showBlock(g)
		g.fcLastKey = g.framesCounter
	}
	return immediate
}

func (g *Game) draw() {
	// Draw a single block. Needs color
	// defer fmt.Println("Block drawn")

	for i, row := range g.board { // row := 0; row < len(g.board); row++ {
		for j, col := range row { //
			if col.used {
				raylib.DrawRectangle(int32(j*40), int32(i*40), 40, 40, col.color)
			} else {

				raylib.DrawRectangle(int32(j*40), int32(i*40), 40, 40, raylib.White)
			}
		}
	}
}

func (g *Game) drawNext(nextX int, nextY int) {
	var next = g.next

	for r, row := range next.cells {
		for c, cell := range row {
			if cell == 'X' {
				raylib.DrawRectangle(int32(nextX+c*30), int32(nextY+r*30), 30, 30, next.color)
			} else {
				raylib.DrawRectangle(int32(nextX+c*30), int32(nextY+r*30), 30, 30, raylib.LightGray)
			}
		}
	}
}

func (g *Game) shiftBlock() *ActiveBlock {
	// return next block as active, initialize next with fresh block
	if g.next == nil {
		// first time there won't be a next block ready
		g.next = &blocks[rand.Intn(len(blocks))]
	}
	ab := &ActiveBlock{block: g.next}
	ab.row = 0
	ab.col = g.cols/2 - 1
	g.next = &blocks[rand.Intn(len(blocks))]
	return ab
}

func main() {
	const fps = 60
	const linesPs = 2
	const keysPs = 6
	const cols = 10
	const rows = 20

	var gameOver = false

	rand.Seed(time.Now().UnixNano())

	raylib.InitWindow(600, 800, "Ivo's GO Tetris")

	raylib.SetTargetFPS(int32(fps))

	board := Game{}
	board.init(rows, cols, fps, linesPs, keysPs)

	ab := board.shiftBlock()
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
		board.draw()
		board.drawNext(420, 150)
		if !gameOver {
			immediate := board.input()
			if !board.paused {
				if !board.blockDown(immediate) {
					ab := board.shiftBlock()
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
				raylib.DrawText("Paused", 100, 400, 40, raylib.Black)
			}
		} else {
			raylib.DrawText("Game Over!", 100, 400, 40, raylib.Black)
		}
		raylib.EndDrawing()
	}

	raylib.CloseWindow()
}
