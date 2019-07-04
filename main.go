package main

import (
	"fmt"
	"sync"
)

const startingSeeds uint = 4 // 4 seeds / cell
const fieldSize = 6 // 6 cells / direction
const c, f uint = 1, 2  // 1 = (c)lose store, 2 = (f)ar store, 0 = not a store
var simulationGroup sync.WaitGroup

// mancala board as a linked list
type Cell struct {
	value uint
	store uint
	next *Cell
}

// map to the stores for convenience & organization
type Field map[uint]*Cell

// simulation result
type Result struct {
	cell *Cell
	field Field
	extraTurn bool
}

// recursive sow from an origin
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

// returns a new field
func newField() Field {
	closeStore, farStore := &Cell{ store: c }, &Cell{ store: f }
	cur := closeStore
	// close => far row
	for i := 0; i < fieldSize; i++ {
		cur.next = &Cell { value: startingSeeds }
		cur = cur.next
	}
	// insert the far store
	cur.next = farStore
	cur = farStore
	// far => close row
	for i := 0; i < fieldSize; i++ {
		cur.next = &Cell { value: startingSeeds }
		cur = cur.next
	}
	cur.next = closeStore // close the loop
	return Field{ c: closeStore, f: farStore }
}

func (field Field) runSim(origin *Cell, player uint, ch chan Result) {
	ch <- Result{ origin, field, origin.sow(player) }
	simulationGroup.Done()
}

func (field Field) runSims(player uint, ch chan Result) {
	simulationGroup.Add(fieldSize)
	for i := 0; i < fieldSize; i++ {
		// clone the field
		simField := field.clone()
		cur := simField[player].next
		// advance to the right cell
		for j := 0; j < i; j++ {
			cur = cur.next
		}
		go simField.runSim(cur, player, ch)
	}
}

// clones field to allow for multiple simulations on the same board state
func (field Field) clone() Field {
	closeStore, farStore := &Cell{ value: field[c].value, store: field[c].store }, &Cell{ value: field[f].value, store: field[f].store }
	oldCur, newCur := field[c], closeStore
	// close => far row
	for i := 0; i < fieldSize; i++ {
		newCur.next = &Cell { value: oldCur.next.value, store: oldCur.next.store }
		newCur = newCur.next
		oldCur = oldCur.next
	}
	// insert the far store
	newCur.next = farStore
	newCur = farStore
	// far => close row
	for i := 0; i < fieldSize; i++ {
		newCur.next = &Cell { value: oldCur.next.value, store: oldCur.next.store }
		newCur = newCur.next
		oldCur = oldCur.next
	}
	newCur.next = closeStore // close the loop
	return Field{ c: closeStore, f: farStore }
}


// formats the field into a readable format
func (field Field) String() string {
	// iterate through linked list to grab values
	closeToFar, farToClose := make([]*Cell, 0, fieldSize), make([]*Cell, 0, fieldSize)
	cur := field[c].next
	for i := 0; i < fieldSize; i++ {
		closeToFar = append(closeToFar, cur)
		cur = cur.next
	}
	cur = field[f].next
	for i := 0; i < fieldSize; i++ {
		farToClose = append(farToClose, cur)
		cur = cur.next
	}

	// make list of values for cleaner string formatting
	vals := make([]interface{}, 0, 2 * fieldSize + 2)

	vals = append(vals, field[f].value)
	for i := 0; i < fieldSize; i++ {
		vals = append(vals, farToClose[i].value)
		vals = append(vals, closeToFar[fieldSize - i - 1].value)
	}
	vals = append(vals, field[c].value)

	return fmt.Sprintf(`
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

// formats the result into a readable format
func (result Result) String() string {
	return fmt.Sprintf(`
Simulation Result:
Board State: %v
Extra Turn: %v
`, result.field, result.extraTurn)
}

func main() {
	field := newField()
	ch := make(chan Result, fieldSize)
	go field.runSims(c, ch)
	simulationGroup.Wait() // wait for simulations to finish
	for result := range ch {
		fmt.Println(result)
		if len(ch) == 0 {
			break
		}
	}
}
