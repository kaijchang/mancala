package main

import (
	"fmt"
)

var STARTING_STONES uint = 3 // 3 stones / cell
var CLOSE, FAR uint = 1, 2 // 0 = no bank

// mancala board as a linked list
type Cell struct {
	value uint
	bank uint
	next *Cell
}

func newBoard() Cell {
	closeBank, farBank := Cell{ bank: CLOSE }, Cell{ bank: FAR }
	var cur, next *Cell
	cur = &closeBank
	for i := 0; i < 6; i++ {
		next = &Cell { value: STARTING_STONES }
		cur.next = next
		cur = next
	}
	cur.next = &farBank
	cur = &farBank
	for i := 0; i < 6; i++ {
		next = &Cell { value: STARTING_STONES }
		cur.next = next
		cur = next
	}
	cur.next = &closeBank // close the loop
	return closeBank
}

func main() {
	board := newBoard()
	cur := board
	for i := 0; i < 14; i++ {
		fmt.Println(cur)
		cur = *cur.next
	}
}
