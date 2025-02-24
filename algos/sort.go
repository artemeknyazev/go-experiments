package main

import (
	"cmp"
	"fmt"
	"reflect"
	"runtime"
)

// https://stackoverflow.com/a/7053871
func GetFunctionName(i any) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

// https://codebase64.org/doku.php?id=base:16bit_xorshift_random_generator (<<7, >>9, <<8)
// https://en.wikipedia.org/wiki/Linear-feedback_shift_register#Xorshift_LFSRs (>>7, <<9, >>13)
func XorShift16(x uint16) uint16 {
	x ^= x >> 7
	x ^= x << 9
	x ^= x >> 13
	return x
}

// https://en.wikipedia.org/wiki/Xorshift
func XorShift32(x uint32) uint32 {
	x ^= x << 13
	x ^= x >> 17
	x ^= x << 5
	return x
}

func FillUint16Array(xs []uint16, seed uint16) uint16 {
	for i := 0; i < len(xs); i++ {
		seed = XorShift16(seed)
		xs[i] = seed
	}
	return XorShift16(seed)
}

func IsSorted[T cmp.Ordered](xs []T) bool {
	for i := 0; i < len(xs)-1; i++ {
		if cmp.Compare(xs[i], xs[i+1]) == 1 {
			return false
		}
	}
	return true
}

func Exchange[T any](A, B *T) {
	*A, *B = *B, *A
}

func CompareExchange[T cmp.Ordered](A, B *T) {
	if cmp.Less(*B, *A) {
		Exchange(A, B)
	}
}

// Based on Sedgewick, Algorithms in C++, prog. 6.1.
func InsertionSort[T cmp.Ordered](xs []T) {
	for i := 1; i < len(xs); i++ {
		// "Sink" the current element to the left while it is smaller
		for j := i; j > 0; j-- {
			CompareExchange(&xs[j-1], &xs[j])
		}
	}
}

// Based on Sedgewick, Algorithms in C++, prog. 6.3.
func InsertionSort2[T cmp.Ordered](xs []T) {
	// Make the leftmost element the "signal" element -- it is less than or equal to other elements
	for i := len(xs) - 1; i > 0; i-- {
		CompareExchange(&xs[i-1], &xs[i])
	}

	// Start from index 2 because the first element is already the smallest one
	for i := 2; i < len(xs); i++ {
		// Remember the current element
		j, v := i, xs[i]

		// Move elements that are to the left of the current and are greater than it one position to the right
		for v < xs[j-1] {
			xs[j] = xs[j-1]
			j--
		}

		// Insert the current element at the freed position
		xs[j] = v
	}
}

// Based on Sedgewick, Algorithms in C++, prog. 6.2.
func SelectionSort[T cmp.Ordered](xs []T) {
	for i := 0; i < len(xs); i++ {
		// Find the minimal element to the right of the current
		min := i
		for j := i + 1; j < len(xs); j++ {
			if cmp.Less(xs[j], xs[min]) {
				min = j
			}
		}

		// Swap the current element with the minimal one
		Exchange(&xs[i], &xs[min])
	}
}

// Based on Sedgewick, Algorithms in C++, prog. 6.4.
func BubbleSort[T cmp.Ordered](xs []T) {
	for i := 0; i < len(xs); i++ {
		// "Sink" the least element to the right of the current one down to the current position
		for j := len(xs) - 1; j > i; j-- {
			CompareExchange(&xs[j-1], &xs[j])
		}
	}
}

// Based on Sedgewick, Algorithms in C++, prog. 6.4.
// Uses an optimization similar to Insertion2: finds the smallest element in the right side, shifts elements one position to the right from current and places the smallest one into the current position.
func BubbleSort2[T cmp.Ordered](xs []T) {
	for i := 0; i < len(xs); i++ {
		// Find the index of the least element to the right of the current one
		min := i
		for j := i + 1; j < len(xs); j++ {
			if cmp.Less(xs[j], xs[min]) {
				min = j
			}
		}

		// Remember the smallest element
		v := xs[min]
		// Move all elements between the current and the smallest one one position to the right
		for j := min; j > i; j-- {
			Exchange(&xs[j-1], &xs[j])
		}

		// Place the smallest element into the current position
		xs[i] = v
	}
}

// Based on Sedgewick, Algorithms in C++, prog. 6.5.
func ShellSort[T cmp.Ordered](xs []T) {
	// Find the suitable stride using the 1, 4, 13, 40, 121, 364, … sequence
	var h int = 1
	for h <= len(xs)/9 {
		h = 3*h + 1
	}

	// From the largest stride down to the smallest one
	for ; h > 0; h /= 3 {
		for i := h; i < len(xs); i++ {
			j, v := i, xs[i]
			for j >= h && cmp.Less(v, xs[j-h]) {
				xs[j] = xs[j-h]
				j = j - h
			}
			xs[j] = v
		}
	}
}

// From golang.org/x/exp/constraints package

type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Integer interface {
	Signed | Unsigned
}

// Based on Sedgewick, Algorithms in C++, prog. 6.17.
func CountSort[T Integer](xs []T) {
	// TODO: implement
}

func BucketSort[T Integer](xs []T) {
	// TODO: implement
}

// Based on Sedgewick, Algorithms in C++, prog. 7.2.
func Partition[T cmp.Ordered](xs []T) int {
	i, j, v := 0, len(xs)-2, xs[len(xs)-1]

	for {
		for i < len(xs) && cmp.Less(xs[i], v) {
			i++
		}
		// xs[i] is the first element that is >= v

		for j >= 0 && cmp.Less(v, xs[j]) {
			j--
			if j == 0 {
				break
			}
		}
		// xs[j] is the first element that is <= v

		if i >= j {
			break
		}

		Exchange(&xs[i], &xs[j])
	}

	Exchange(&xs[i], &xs[len(xs)-1])
	return i
}

// Based on Sedgewick, Algorithms in C++, prog. 7.1.
func QuickSort[T cmp.Ordered](xs []T) {
	if len(xs) <= 1 {
		return
	}

	p := Partition(xs)
	QuickSort(xs[:p])
	QuickSort(xs[p+1:])
}

// Based on Sedgewick, Algorithms in C++, prog. 7.3.
func NonRecursiveQuickSort[T cmp.Ordered](xs []T) {
	stack := make([]int, 0, 50)      // Max 25-level stack containing `l, r` pairs
	stack = append(stack, 0)         // Push left index
	stack = append(stack, len(xs)-1) // Push right index, inclusive

	for len(stack) > 0 {
		r := stack[len(stack)-1]     // Pop right index, inclusive
		l := stack[len(stack)-2]     // Pop left
		stack = stack[:len(stack)-2] // Fix stack length

		if r <= l {
			continue
		}

		// `r+1` because right border of a slice is not included
		// `l+` because `Partition` returns an offset from `l`
		i := l + Partition(xs[l:r+1])

		// Compare lengths of left (l,i-1) and right (i+1,r) parts
		// to push the smallest part to the top to limit stack growth
		if i-l > r-i { // Right part is smaller
			stack = append(stack, l)   // Push left left
			stack = append(stack, i-1) // Push left right
			stack = append(stack, i+1) // Push right left
			stack = append(stack, r)   // Push right right
		} else { // Left part is smaller
			stack = append(stack, i+1) // Push right left
			stack = append(stack, r)   // Push right right
			stack = append(stack, l)   // Push left left
			stack = append(stack, i-1) // Push left right
		}
	}
}

const hybridQuickSortMinArrayLength = 10

// Based on Sedgewick, Algorithms in C++, prog. 7.4.
func medianOfThreeQuickSort[T cmp.Ordered](xs []T) {
	// Skip small subarrays, they are sorted on the next step
	ln := len(xs)
	if ln <= hybridQuickSortMinArrayLength {
		return
	}

	// Median-of-three: compare three elements and move median to the first position
	r := ln - 1
	Exchange(&xs[r/2], &xs[r-1])
	CompareExchange(&xs[0], &xs[r-1])
	CompareExchange(&xs[0], &xs[r])
	CompareExchange(&xs[r-1], &xs[r])

	p := Partition(xs[1:ln])
	QuickSort(xs[:p])
	QuickSort(xs[p+1:])
}

// Based on Sedgewick, Algorithms in C++, prog. 7.4.
func HybridQuickSort[T cmp.Ordered](xs []T) {
	medianOfThreeQuickSort(xs)
	InsertionSort2(xs)
}

func TestSort(buf []uint16, sort func(xs []uint16), seed uint16, pow int, iters int) {
	_, file, line, _ := runtime.Caller(1)

	nextSeed := seed
	for p := 0; p <= pow; p++ {
		length := 1 << p
		xs := buf[:length]

		for i := 0; i < iters; i++ {
			initialSeed := nextSeed
			nextSeed = FillUint16Array(xs, initialSeed)
			// fmt.Printf("  seed:   %v\n", initialSeed)
			// fmt.Printf("  before: %v\n", xs)
			sort(xs)
			// fmt.Printf("  after:  %v\n", xs)
			// fmt.Println("  sorted:", IsSorted(xs))
			if !IsSorted(xs) {
				fmt.Printf("%s failed for len=%d, seed=%d at %s:%d\n", GetFunctionName(sort), length, initialSeed, file, line)
				return
			}
		}
	}
}

func main() {
	// {
	// 	fmt.Printf("Selection sort:\n")
	// 	var length int = 4
	// 	var seed uint16 = 3508
	// 	xs := make([]uint16, length)
	// 	FillUint16Array(xs, seed)
	// 	fmt.Printf("  before: %v\n", xs)
	// 	// ix := Partition(xs)
	// 	// fmt.Printf("  ix:     %d\n", ix)
	// 	QuickSortNonRecursive(xs)
	// 	fmt.Printf("  after:  %v\n", xs)
	// }
	// return

	pow := 12
	buf := make([]uint16, 1<<pow)

	TestSort(buf, InsertionSort, 1, pow, 100)
	TestSort(buf, SelectionSort, 1, pow, 100)
	TestSort(buf, InsertionSort2, 1, pow, 100)
	TestSort(buf, BubbleSort, 1, pow, 100)
	TestSort(buf, BubbleSort2, 1, pow, 100)
	TestSort(buf, ShellSort, 1, pow, 100)
	TestSort(buf, QuickSort, 1, pow, 10000)
	TestSort(buf, NonRecursiveQuickSort, 1, pow, 10000)
	TestSort(buf, HybridQuickSort, 1, pow, 10000)

	// TestSort(buf, CountSort, 1, pow, 100) // TODO:
	// TestSort(buf, BucketSort, 1, pow, 100) // TODO:
}
