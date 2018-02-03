package main

import (
	"fmt"
	"math/rand"

	"github.com/gen2brain/raylib-go/raylib"
)

var block = []string{
	"XX",
	"XX",
}

func drawBlock(x int, y int) {
	// Draw a single block. Needs color
	defer fmt.Println("Block drawn")

	raylib.DrawRectangle(int32(x), int32(y), 20, 20, raylib.Blue)
}

func initBoard(rows int, cols int) [][]string {
	board := make([][]string, rows)

	for i := 0; i < len(board); i++ {
		board[i] = make([]string, cols)
	}

	return board
}

func printBoard(board [][]string) {
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			fmt.Print(board[i][j])
		}
		fmt.Println()
	}
}

func putSomeBlocks(board [][]string) {
	board[1][1] = "X"
	board[1][2] = "X"
	board[2][1] = "X"
	board[2][2] = "X"
}

func main() {
	raylib.InitWindow(450, 800, "Ivo's GO Tetris")

	raylib.SetTargetFPS(6)
	fmt.Println("My favorite number is", rand.Intn(10))

	const cols = 10
	const rows = 20

	board := initBoard(rows, cols)

	putSomeBlocks(board)
	printBoard(board)

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
