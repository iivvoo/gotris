package main

import "github.com/gen2brain/raylib-go/raylib"

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
}, color: raylib.DarkBlue}
var reverseZ = Block{cells: []string{
	" XX",
	"XX ",
}, color: raylib.Gray}
var triangle = Block{cells: []string{
	" X ",
	"XXX",
}, color: raylib.Purple}

var blocks = []Block{square, line, leftL, rightL, z, reverseZ, triangle}
