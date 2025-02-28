package main

import (
	"fmt"
	"math/rand"
	"sort"
)

func RowColToIx(row, col, dim uint16) uint16 {
	return row*dim + col
}

func IxToRowCol(ix, dim uint16) (row, col uint16) {
	row = ix / dim
	col = ix - row*dim
	return row, col
}

func CellCandidates(constraints []bool, ix uint16, n uint16) []bool {
	return constraints[ix*n : ix*n+n]
}

func isRowSolved(board []uint16, row uint16, dim uint16, n uint16) bool {
	numbers := make([]bool, n)

	for col := uint16(0); col < dim; col++ {
		ix := RowColToIx(row, col, dim)

		v := board[ix]
		if v == 0 {
			return false // Not filled
		}

		num := &numbers[v-1]
		if *num {
			return false // Incorrect
		}
		*num = true
	}
	// fmt.Printf("row %d numbers = %v\n", row, numbers)

	for i := uint16(0); i < n; i++ {
		if !numbers[i] {
			return false
		}
	}
	return true
}

func isColumnSolved(board []uint16, col uint16, dim uint16, n uint16) bool {
	numbers := make([]bool, n)

	for row := uint16(0); row < dim; row++ {
		ix := RowColToIx(row, col, dim)

		v := board[ix]
		if v == 0 {
			return false // Not filled
		}

		num := &numbers[v-1]
		if *num {
			return false // Incorrect
		}
		*num = true
	}
	// fmt.Printf("col %d numbers = %v\n", col, numbers)

	for i := uint16(0); i < n; i++ {
		if !numbers[i] {
			return false
		}
	}
	return true
}

func isSubCellSolved(board []uint16, rowStart, rowEnd uint16, colStart, colEnd uint16, dim uint16, n uint16) bool {
	numbers := make([]bool, n)

	for row := rowStart; row <= rowEnd; row++ {
		for col := colStart; col <= colEnd; col++ {
			ix := RowColToIx(row, col, dim)

			v := board[ix]
			if v == 0 {
				return false // Not filled
			}

			num := &numbers[v-1]
			if *num {
				return false // Incorrect
			}
			*num = true
		}
	}
	// fmt.Printf("cell %d %d numbers = %v\n", rowStart, colStart, numbers)

	for i := uint16(0); i < n; i++ {
		if !numbers[i] {
			return false
		}
	}
	return true
}

func IsBoardSolved(board []uint16, dim uint16, n uint16) bool {
	for row := uint16(0); row < dim; row++ {
		if !isRowSolved(board, row, dim, n) {
			println("row", row, "not solved")
			return false
		}
	}
	for col := uint16(0); col < dim; col++ {
		if !isColumnSolved(board, col, dim, n) {
			println("col", col, "not solved")
			return false
		}
	}
	for row := uint16(0); row < dim; row += 3 {
		for col := uint16(0); col < dim; col += 3 {
			if !isSubCellSolved(board, row, row+2, col, col+2, dim, n) {
				println("cell", row, col, "not solved")
				return false
			}
		}
	}
	return true
}

func printSubCellRow(board []uint16) {
	for _, v := range board {
		if v == 0 {
			fmt.Printf(" ")
		} else {
			fmt.Printf("%d", v)
		}
	}
}

func printRow(board []uint16) {
	fmt.Printf("|")
	printSubCellRow(board[0:3])
	fmt.Printf("|")
	printSubCellRow(board[3:6])
	fmt.Printf("|")
	printSubCellRow(board[6:9])
	fmt.Printf("|\n")
}

func PrintBoard(board []uint16, dim uint16) {
	delim := "+---+---+---+\n"
	fmt.Printf(delim)
	printRow(board[0*dim:])
	printRow(board[1*dim:])
	printRow(board[2*dim:])
	fmt.Printf(delim)
	printRow(board[3*dim:])
	printRow(board[4*dim:])
	printRow(board[5*dim:])
	fmt.Printf(delim)
	printRow(board[6*dim:])
	printRow(board[7*dim:])
	printRow(board[8*dim:])
	fmt.Printf(delim)
}

func fillConstraintsForCell(constraints []bool, v uint16, row uint16, col uint16, dim uint16, n uint16) {
	// Update row
	for c := uint16(0); c < dim; c++ {
		ixx := RowColToIx(row, c, dim)
		// println("  row", row, c, ixx, ixx*n+v-1)
		constraints[ixx*n+v-1] = true
	}

	// Update column
	for r := uint16(0); r < dim; r++ {
		ixx := RowColToIx(r, col, dim)
		// println("  col", r, col, ixx, ixx*n+v-1)
		constraints[ixx*n+v-1] = true
	}

	// Update subboard
	rr, cc := row/3, col/3
	for r := uint16(0); r < 3; r++ {
		for c := uint16(0); c < 3; c++ {
			ixx := RowColToIx(rr*3+r, cc*3+c, dim)
			// println("  sub", rr*3+r, cc*3+c, ixx, ixx*n+v-1)
			constraints[ixx*n+v-1] = true
		}
	}
}

func FillConstraints(constraints []bool, board []uint16, dim uint16, n uint16) {
	for row := uint16(0); row < dim; row++ {
		for col := uint16(0); col < dim; col++ {
			ix := RowColToIx(row, col, dim)
			v := board[ix]
			if v == 0 {
				continue
			}
			// println("x", row, col, ix, v)

			fillConstraintsForCell(constraints, v, row, col, dim, n)
		}
	}
}

func Solve(board []uint16, dim uint16, n uint16, step uint16) bool {
	if int(step) == len(board) {
		return IsBoardSolved(board, dim, n)
	}

	// Find constraints for each cell
	var constraints = make([]bool, dim*dim*n)
	FillConstraints(constraints, board, dim, n)
	// fmt.Printf("constraints %v", constraints)

	// List free cells and how many candidate numbers they can have
	type Move struct {
		Ix    uint16
		Count uint16
	}
	moves := make([]Move, 0, dim*dim)
	for r := uint16(0); r < dim; r++ {
		for c := uint16(0); c < dim; c++ {
			ix := RowColToIx(r, c, dim)
			if board[ix] != 0 {
				continue // Already filled
			}

			cnt := uint16(0)
			candidates := CellCandidates(constraints, ix, n)
			for _, used := range candidates {
				if !used {
					cnt++
				}
			}
			if cnt == 0 {
				return false // No moves for cell
			}
			moves = append(moves, Move{Ix: ix, Count: cnt})
		}
	}

	// Sort move candidates so that cells w/ less candidate numbers count are first
	sort.Slice(moves, func(i, j int) bool {
		left, right := moves[i], moves[j]
		if left.Count == 0 {
			return false
		}
		if right.Count == 0 {
			return false
		}
		return left.Count < right.Count
	})

	// println("Solve step", step, ",", len(moveCandidates), "moves left")

	// Try making moves
	for _, move := range moves {
		candidates := CellCandidates(constraints, move.Ix, n)
		row, col := IxToRowCol(move.Ix, dim)
		// fmt.Printf("checking move %d %d %v\n", row, col, candidates)
		for candidate, used := range candidates {
			if used {
				continue
			}
			candidate++                        // Zero-based indexing for 1...n numbers
			board[move.Ix] = uint16(candidate) // Make move
			println("  make (", row, col, ") =", candidate, "at step", step)
			// println(int(step), len(board))
			// PrintBoard(board, dim)
			if Solve(board, dim, n, step+1) {
				return true // Solved
			}
			board[move.Ix] = 0 // Revert move
			println("revert (", row, col, ") =", candidate, "at step", step)
		}
	}

	return false // Not solved
}

// From Skiena, Algorithms, ed. 2, ch. 7.3.
func FillInit1(board []uint16, dim uint16) {
	x := uint16(0)

	copy(board[0*dim:], []uint16{x, x, x /**/, x, x, x /**/, x, 1, 2})
	copy(board[1*dim:], []uint16{x, x, x /**/, x, 3, 5 /**/, x, x, x})
	copy(board[2*dim:], []uint16{x, x, x /**/, 6, x, x /**/, x, 7, x})

	copy(board[3*dim:], []uint16{7, x, x /**/, x, x, x /**/, 3, x, x})
	copy(board[4*dim:], []uint16{x, x, x /**/, 4, x, x /**/, 8, x, x})
	copy(board[5*dim:], []uint16{1, x, x /**/, x, x, x /**/, x, x, x})

	copy(board[6*dim:], []uint16{x, x, x /**/, 1, 2, x /**/, x, x, x})
	copy(board[7*dim:], []uint16{x, 8, x /**/, x, x, x /**/, x, 4, x})
	copy(board[8*dim:], []uint16{x, 5, x /**/, x, x, x /**/, 6, x, x})
}

// From Skiena, Algorithms, ed. 2, ch. 7.3.
func FillSol1(board []uint16, dim uint16) {
	copy(board[0*dim:], []uint16{6, 7, 3 /**/, 8, 9, 4 /**/, 5, 1, 2})
	copy(board[1*dim:], []uint16{9, 1, 2 /**/, 7, 3, 5 /**/, 4, 8, 6})
	copy(board[2*dim:], []uint16{8, 4, 5 /**/, 6, 1, 2 /**/, 9, 7, 3})

	copy(board[3*dim:], []uint16{7, 9, 8 /**/, 2, 6, 1 /**/, 3, 5, 4})
	copy(board[4*dim:], []uint16{5, 2, 6 /**/, 4, 7, 3 /**/, 8, 9, 1})
	copy(board[5*dim:], []uint16{1, 3, 4 /**/, 5, 8, 9 /**/, 2, 6, 7})

	copy(board[6*dim:], []uint16{4, 6, 9 /**/, 1, 2, 8 /**/, 7, 3, 5})
	copy(board[7*dim:], []uint16{2, 8, 7 /**/, 3, 5, 6 /**/, 1, 4, 9})
	copy(board[8*dim:], []uint16{3, 5, 1 /**/, 9, 4, 7 /**/, 6, 2, 8})
}

func RemoveFromBoard(board []uint16, dim uint16, m int) {
	for i, v := range rand.Perm(int(dim * dim)) {
		if i == m {
			return
		}
		board[v] = 0
	}
}

func main() {
	// TODO: https://stackoverflow.com/questions/6924216/how-to-generate-sudoku-boards-with-unique-solutions

	const dim uint16 = 9 // Square side
	const n uint16 = 9   // Count of distinct numbers
	var board = make([]uint16, dim*dim)

	// FillInit1(board, dim)
	// board[RowColToIx(5, 4, dim)] = 1
	// GenerateBoard(board, dim)
	FillSol1(board, dim)
	RemoveFromBoard(board, dim, 40)
	PrintBoard(board, dim)

	step := uint16(0)
	for _, v := range board {
		if v != 0 {
			step++
		}
	}
	solved := Solve(board, dim, n, step)
	println("Solved =", solved)
	PrintBoard(board, dim)
}
