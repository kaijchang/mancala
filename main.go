package main

import "fmt"

const startingSeeds uint = 3 // 3 stones / cell
const fieldSize = 6 // 6 cells / direction
const c, f uint = 1, 2  // 1 = (c)lose store, 2 = (f)ar store, 0 = not a store

// mancala board as a linked list
type Cell struct {
	value uint
	store uint
	next *Cell
}

// recursive sow from an origin, returns whether player gets an extra turn or not
func (origin *Cell) sow(player uint) bool {
	// pick up all the stones
	hand := origin.value
	origin.value = 0
	cur := origin.next
	for {
		// skip opponent's store
		if cur.store != 0 && cur.store != player {
			// move on
			cur = cur.next
			continue
		}
		// drop a seed
		cur.value++
		hand--
		// if out of seeds, start sowing again from that cell
		if hand == 0 {
			// extra turn if you end at your store
			if cur.store == player {
				return true
			}
			// end if the cell only has the seed you just put in
			if cur.value == 1 {
				return false
			}
			return cur.sow(player)
		}
		// move on
		cur = cur.next
	}
}

// returns a new field with two entry points at either store
func newField() (*Cell, *Cell) {
	closeStore, farStore := &Cell{ store: c }, &Cell{ store: f }
	var cur, next *Cell
	cur = closeStore
	// close => far row
	for i := 0; i < fieldSize; i++ {
		next = &Cell { value: startingSeeds }
		cur.next = next
		cur = next
	}
	// insert the far store
	cur.next = farStore
	cur = farStore
	// far => close row
	for i := 0; i < fieldSize; i++ {
		next = &Cell { value: startingSeeds }
		cur.next = next
		cur = next
	}
	cur.next = closeStore // close the loop
	return closeStore, farStore
}

// prints the whole field from any origin
func printField(origin *Cell) {
	temp, closeToFar, farToClose := make([]*Cell, 0, fieldSize), make([]*Cell, 0, fieldSize), make([]*Cell, 0, fieldSize)
	var closeBank, farBank *Cell
	var to, tempDirection uint // the to variable keeps track of which direction we're going.
	                           // the tempDirection variable keeps track of which direction we were going originally
	                           // so we know where to append the leftover cells from when we didn't know the direction yet
	cur := origin
	for i := 0; i < fieldSize * 2 + 2; i++ {
		if cur.store != 0 {
			if cur.store == 1 {
				closeBank = cur
			} else if cur.store == 2 {
				farBank = cur
			}
			to = cur.store
			if tempDirection == 0 {
				tempDirection = cur.store
			}
		} else {
			switch to {
				case 0: // we don't know the direction yet
					temp = append(temp, cur)
				case 1: // close => far row
					closeToFar = append(closeToFar, cur)
				case 2: // far => close row
					farToClose = append(farToClose, cur)
			}
		}
		cur = cur.next
	}

	// append leftover cells
	if tempDirection == 1 {
		farToClose = append(farToClose, temp...)
	} else if tempDirection == 2 {
		closeToFar = append(closeToFar, temp...)
	}

	// make list of values for cleaner string formatting
	vals := make([]interface{}, 0, 14)

	vals = append(vals, farBank.value)
	for i := 0; i < fieldSize; i++ {
		vals = append(vals, farToClose[i].value)
		vals = append(vals, closeToFar[fieldSize - i - 1].value)
	}
	vals = append(vals, closeBank.value)

	fmt.Printf(`
  %v
____
|%v|%v|
|%v|%v|
|%v|%v|
----
|%v|%v|
|%v|%v|
|%v|%v|
----
  %v
`, vals...)
}

func main() {
	closeStore, _ := newField()
	printField(closeStore)
	fmt.Println(closeStore.next.sow(c)) // start sowing at first cell past close store
	printField(closeStore)
}
