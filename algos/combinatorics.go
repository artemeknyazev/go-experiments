package main

import (
	"fmt"
	"time"
)

func firstCandidate(candidates []bool) int {
	for candidate, used := range candidates {
		if !used {
			return candidate
		}
	}

	panic("No candidates!") // TODO: use assert before?
}

func nextCandidate(candidates []bool, k int) (int, bool) {
	// Check overflow
	if k >= len(candidates) {
		return 0, false
	}

	// Find first unused candidate after k
	k++
	for k < len(candidates) {
		if !candidates[k] {
			return k, true
		}
		k++
	}

	// No unused candidates
	return 0, false
}

// Ordered w/o repetition. Create the first n-choose-k permutation.
// xs is array where keys are position and values are object indices. n is number of objects.
//
//	xs := make([]int, 4)   // 4 is k, count of objects in permutation
//	InitPermitation(xs, 6) // 6 is n, total count of objects
func InitPermutation(xs []int, n int) {
	candidates := make([]bool, n) // Candidate array to mark used indices

	// From left to right fill each position with first available candidate w/o repetition
	for k := 0; k < len(xs); k++ {
		candidate := firstCandidate(candidates)
		candidates[candidate] = true
		xs[k] = candidate
	}
}

// Ordered w/o repetition. Next n-choose-k permutation. Wraps around.
// xs is array where keys are position and values are object indices. n is number of objects.
//
//	xs := make([]int, 4)   // 4 is k, count of objects in permutation
//	InitPermitation(xs, 6) // 6 is n, total count of objects
//	NextPermutation(xs, 6) // Next permutation
//	NextPermutation(xs, 6) // Next permutation again
func NextPermutation(xs []int, n int) {
	// Mark used indices in candidate array
	candidates := make([]bool, n)
	for _, v := range xs {
		candidates[v] = true
	}

	r := len(xs) - 1 // Rightmost position in xs

	k := r // Start from rightmost position

	// Do +1-analogue with carry choosing next not-used candidate
	for k >= 0 {
		previous := xs[k] // Previous value at current position
		candidate, ok := nextCandidate(candidates, previous)
		candidates[previous] = false // Unmark previous value

		if ok { // Candidate exists, overwrite at current
			candidates[candidate] = true
			xs[k] = candidate
			break
		}

		k-- // No candidate, overflow, move one position left
	}

	// Fill positions after incremented one with first available candidates
	if k < r {
		k++ // Point to first overflowed position
		for k <= r {
			candidate := firstCandidate(candidates)
			candidates[candidate] = true
			xs[k] = candidate
			k++
		}
	}
}

// Unordered w/o repetition. bs is array where keys are object ids, values are flags, if object is included in a subset. Works analogous to +1 for len(xs) digits binary number.
func NextSubset(bs []bool) {
	r := len(bs) - 1 // Rightmost position

	k := r // Start from rightmost position

	// Do +1 with carry
	for k >= 0 {
		value := bs[k]

		if !value {
			// No overflow, set current
			bs[k] = true
			break
		}

		k-- // Overflow, move one position left
	}

	// Exclude all elements after not overflowed one
	if k < r {
		k++ // Point to first overflowed position
		for k <= r {
			bs[k] = false
			k++
		}
	}
}

// Unordered with repetition. xs is array where key is object id, value is object count, at most n. n is max count of objects.
//
//	NextTuple([]int{0,1,3}, 4) // [0 1 4]
//	NextTuple([]int{0,1,4}, 4) // [0 2 0]
func NextTuple(xs []int, n int) {
	r := len(xs) - 1

	k := r

	// Do +1 with carry
	for k >= 0 {
		xs[k]++ // TODO: check int overflow

		if xs[k] > n {
			// Overflow, wrap to 0, carry
			xs[k] = 0
			k--
			continue
		}

		break
	}
}

// Unordered with repetition.
// xs is an array where key is object id, value is count of such object.
// m is sum of counts of objects.
//
// E.g. [0 1 0 2 3 0], m=6, object id=1 count is 1, object id=3 count is 2, object id=4 count is 3
//
//	NextUnorderedRepeatTotalLimit([]int{0,1,0,2,3,0}, 6) // [0 1 0 3 0 2]
//
// TODO: find better name
func NextUnorderedRepeatTotalLimit(xs []int, m int) {
	r := len(xs) - 1 // Right boundary

	// Find first non-zero from the right
	k := r
	for k >= 0 {
		if xs[k] == 0 {
			k--
			continue
		}
		break
	}

	// Zero xs initialization case
	if k < 0 {
		xs[r] = m
		return
	}

	// Wrapping case
	if xs[0] == m {
		xs[0] = 0
		xs[r] = m
		return
	}

	k-- // Move to the left of it

	// Calculate how much of m is left for k-th position
	y := m
	for i := 1; i < k; i++ {
		y = y - xs[i]
	}

	xs[k]++ // Increment left neighbour of first non-zero; TODO: check int overflow
	if xs[k] > y {
		xs[k] = 0 // Prevent overflow
	}

	z := y - xs[k] // Leftovers
	k++
	for k < r {
		xs[k] = 0 // Zero array to the right
		k++
	}
	xs[r] = z // Put leftovers into the last place
}

/*
Permutations of 3 elems:

- 0 - 123 - 000
- 1 - 132 - 010
- 2 - 213 - 100
- 3 - 231 - 110
- 4 - 312 - 200
- 5 - 321 - 210

Permutations of 4 elems:
- ix - perm - ixs of next elem
- 00 - 1234 - 0000
- 01 - 1243 - 0010 - period 1=1! of last-1
- 02 - 1324 - 0100
- 03 - 1342 - 0110 - period 2=2! of last-2
- 04 - 1423 - 0200
- 05 - 1432 - 0210
- 06 - 2134 - 1000 - period 6=3! of last-3
- 07 - 2143 - 1010
- 08 - 2314 - 1100
- 09 - 2341 - 1110
- 10 - 2413 - 1200
- 11 - 2431 - 1210
- 12 - 3124 - 2000
- 13 - 3142 - 2010
- 14 - 3214 - 2100
- 15 - 3241 - 2110
- 16 - 3412 - 2200
- 17 - 3421 - 2210
- 18 - 4123 - 3000
- 19 - 4132 - 3010
- 20 - 4213 - 3100
- 21 - 4231 - 3110
- 22 - 4312 - 3200
- 23 - 4321 - 3210
- 24 - 1234 - 0000 - period 24=4! of last-4

Look at 15-th 3241 permutation:
- 15 = 2 * 3! + 3; use 2nd (0-based!) elem in 1234 string which is 3
- 03 = 1 * 2! + 1; use 1st elem in 124 (3 removed) which is 2
- 01 = 1 * 1! + 0; use 1st elem in 14 (2 removed) which is 4
- use last left elem which is 1

3241 permutation has 2110 index string, meaning use 2nd of 1234, then
1st of 124, then 1st of 14, then last 1 left. To convert 2110 to actual
indexes can do the following:
- (2)110 -> 2110, nothing behind, no changes
- 2(1)10 -> 2110, 2 behind current 1, 2>=1, no changes
- 21(1)0 -> 2130, 1 (1=1) and 2 (2>=1) behind current 1, increase 1 to 3
- 211(0) -> 2130, no 0s behind current 0, no changes

Same for 09 - 2341 - 1110:
- (1)110 -> 1110, nothing behind, no changes
- 1(1)10 -> 1210, 1 behind current 1, increase to 2
- 12(1)0 -> 1230, 1 and 2 behind current 1, increase to 3
- 123(0) -> 1230, no 0s behind, no changes

Same for 19 - 4132 - 3010:
- (3)010 -> 3010, no changes
- 3(0)10 -> 3010, no changes
- 30(1)0 -> 3020, 0 behind current 1, increase to 3
- 302(0) -> 3021, 0 behind current 0, increase to 1

Or again for 19 - 4132 - 3010:
- init,   [0:0, 1:1, 2:2, 3:3]
- step 0, take 3, [0:0, 1:1, 2:2, 3:_]
- step 1, take 0, [0:1, 1:2, 2:_, 3:_]
- step 2, take 1, [0:1, 1:_, 2:_, 3:_]
- step 3, take 0, [0:_, 1:_, 2:_, 3:_]
- i.e. for each step j do a[i]=a[i+1] where i is from ix taken to 4-j, where 4 is num of elems
*/

func main() {
	if false {
		n, k := 4, 3
		xs := make([]int, k)
		InitPermutation(xs, n)
		fmt.Printf("%v %v\n", n, xs)
		for i := 0; i < 20; i++ {
			NextPermutation(xs, n)
			fmt.Printf("%v %v\n", n, xs)
		}
		return
	}

	if false {
		xs := make([]bool, 3)
		fmt.Printf("%v\n", xs)
		for i := 0; i < 10; i++ {
			NextSubset(xs)
			fmt.Printf("%v\n", xs)
		}
		return
	}

	if false {
		n := 3
		xs := make([]int, 3)
		fmt.Printf("%v\n", xs)
		for i := 0; i < 20; i++ {
			NextTuple(xs, n)
			fmt.Printf("%v\n", xs)
		}
		return
	}

	if false {
		dt := 1 * time.Millisecond
		n := 6
		m := 6
		xs := make([]int, n)
		for {
			NextUnorderedRepeatTotalLimit(xs, m)
			fmt.Printf("%v\n", xs)
			time.Sleep(dt)
		}
		return
	}

}
