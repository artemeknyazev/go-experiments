package main

import (
	"cmp"
	"fmt"
	"os"
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

// Based on Sedgewick, Algorithms in C++, prog. 7.5.
func ThreeWayPartition[T cmp.Ordered](xs []T) (int, int) {
	// TODO: implement
	return 0, 0
}

// Based on Sedgewick, Algorithms in C++, prog. 7.5.
func ThreeWayQuickSort[T cmp.Ordered](xs []T) {
	if len(xs) <= 1 {
		return
	}
	i, j := ThreeWayPartition(xs)
	ThreeWayQuickSort(xs[:i])
	ThreeWayQuickSort(xs[j:])
}

// Based on Sedgewick, Algorithms in C++, prog. 7.6.
func Select[T cmp.Ordered](xs []T, k int) {
	if len(xs) <= 1 {
		return
	}

	p := Partition(xs)
	if p > k {
		Select(xs[:p], k)
	}
	if p < k {
		Select(xs[p+1:], k-p-1)
	}
}

// Based on Sedgewick, Algorithms in C++, prog. 7.7.
func NonRecursiveSelect[T cmp.Ordered](xs []T, k int) {
	l, r := 0, len(xs)
	for r > l+1 {
		p := l + Partition(xs[l:r])
		if p >= k {
			r = p
		}
		if p <= k {
			l = p + 1
		}
	}
}

func Assert(flag bool, msg string) {
	if !flag {
		_, file, line, _ := runtime.Caller(1)
		fmt.Fprintf(os.Stderr, "%s at %s:%d", msg, file, line)
		fmt.Fprint(os.Stderr, "\n")
		os.Exit(1)
	}
}

// Based on Sedgewick, Algorithms in C++, prog. 8.1.
// Assumes out, xs, and ys do not intersect!
// TODO: write a test
func MergeInto[T cmp.Ordered](out []T, xs []T, ys []T) {
	N, M := len(xs), len(ys)

	Assert(len(out) >= N+M, "out array is smaller than sum of lengths of input arrays")
	Assert(!IsSorted(xs), "xs is not sorted")
	Assert(!IsSorted(ys), "ys is not sorted")

	for i, j, k := 0, 0, 0; k < N+M; k++ {
		if i == N {
			out[k] = ys[j]
			j++
			continue
		}

		if j == M {
			out[k] = xs[i]
			i++
			continue
		}

		if cmp.Less(xs[i], ys[j]) {
			out[k] = xs[i]
			i++
		} else {
			out[k] = ys[j]
			j++
		}
	}

	Assert(!IsSorted(xs[:N+M]), "out array must be sorted")
}

// Based on Sedgewick, Algorithms in C++, prog. 8.2.
// Does only one check inside the loop compared to MergeInto which does three checks.
// TODO: write a test
func MergeInside[T cmp.Ordered](xs []T, m int, aux []T) {
	// fmt.Printf("MergeInside:\n")
	// fmt.Printf("  m: %v\n", m)
	// fmt.Printf("  l: %v\n", xs[:m+1])
	// fmt.Printf("  r: %v\n", xs[m+1:])

	Assert(IsSorted(xs[:m+1]), "input array part before provided index must be sorted")
	Assert(IsSorted(xs[m+1:]), "input array part after provided index must be sorted")
	Assert(len(aux) >= len(xs), "buffer is smaller than the merged array")

	i, j, r := 0, 0, len(xs)-1

	// Copies left part of xs into left part of aux in forward order
	// aux[l:m] = xs[l:m]
	for i = m + 1; i > 0; i-- {
		aux[i-1] = xs[i-1]
	}
	// i = 0, left border of aux = left border of the left part of xs

	// Copies right part of xs into right part of aux in reverse order
	// aux[m:r] = reverse(xs[m+1:r])
	for j = m; j < r; j++ {
		aux[r+m-j] = xs[j+1]
	}
	// j = r, right border of aux = left border of the reversed right part of xs

	// Merges aux[l:m] with left part of xs and aux[m:r] with reversed right part of xs into xs
	for k := 0; k <= r; k++ {
		if cmp.Less(aux[j], aux[i]) {
			xs[k] = aux[j]
			j--
		} else {
			xs[k] = aux[i]
			i++
		}
	}

	// fmt.Printf("  xs: %v\n", xs)
	Assert(IsSorted(xs), "input array after merge must be sorted")
}

// Merges two slices into an output slice using an auxiliary array.
// out, xs, ys can intersect. Assumes aux doesn't intersect with other arrays.
//
// Assumes that:
//   - out, xs, ys can intersect
//   - aux doesn't intersect with other arrays
//   - xs and ys are sorted
//   - len(out) == len(xs)+len(ys)
//   - len(aux) >= len(xs)+len(ys)
func MergeInto2[T cmp.Ordered](out []T, xs []T, ys []T, aux []T) {
	// fmt.Printf("MergeInto2:\n")
	// fmt.Printf("  xs: %v\n", xs)
	// fmt.Printf("  ys: %v\n", ys)
	// fmt.Printf("  aux: %v\n", aux)
	// fmt.Printf("  out: %v\n", out)
	Assert(IsSorted(xs), "xs must be sorted")
	Assert(IsSorted(ys), "ys must be sorted")
	Assert(len(out) == len(xs)+len(ys), "out is smaller than xs and ys combined")
	Assert(len(aux) >= len(out), "aux is smaller than out")

	lx, ly := len(xs), len(ys)

	// Right boundary of aux and out
	r := lx + ly - 1

	// Copy xs forwards into aux
	for i := 0; i < lx; i++ {
		aux[i] = xs[i]
	}

	// Copy ys backwards into aux
	for j := lx; j <= r; j++ {
		aux[j] = ys[r-j] // r-j = lx + ly - 1 - (lx + offset) = ly - 1 - offset
	}

	// k counts total number of elements minus 1 which is r
	// i starts from the left boundary 0 and goes up to lx-1
	// j starts from the right boundary r and goes down to lx
	for k, i, j := 0, 0, r; k <= r; k++ {
		if cmp.Less(aux[j], aux[i]) {
			out[k] = aux[j]
			j--
		} else {
			out[k] = aux[i]
			i++
		}
	}

	Assert(IsSorted(out), "output array after merge must be sorted")
}

// Based on Sedgewick, Algorithms in C++, prog. 8.2.
// Uses MergeInto2 to simplify working with slices.
func TopDownMergeSort[T cmp.Ordered](xs []T, aux []T) {
	// fmt.Printf("TopDownMergeSort")
	// fmt.Printf("  xs: %v\n", xs)
	ln := len(xs)
	if ln <= 1 {
		return
	}

	// Using an index to split an array

	// m := (ln - 1) / 2 // =(l+r)/2 where l=0, r=len(xs)-1
	// fmt.Printf("  m: %v\n", m)
	// TopDownMergeSort(xs[:m+1], aux)
	// TopDownMergeSort(xs[m+1:], aux)
	// MergeInside(xs, m, aux)

	// Or using two slices, which looks nicer

	m := (ln + 1) / 2 // =1+(ln-1)/2
	ys, zs := xs[:m], xs[m:]
	TopDownMergeSort(ys, aux)
	TopDownMergeSort(zs, aux)
	MergeInto2(xs, ys, zs, aux)
}

func TestSort(buf []uint16, fn func(xs []uint16), seed uint16, pow int, iters int) {
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
			fn(xs)
			// fmt.Printf("  after:  %v\n", xs)
			// fmt.Println("  sorted:", IsSorted(xs))
			if !IsSorted(xs) {
				fmt.Printf("%s failed for len=%d, seed=%d at %s:%d\n", GetFunctionName(fn), length, initialSeed, file, line)
				return
			}
		}
	}
}

// https://stackoverflow.com/a/75435478
func All[T any](xs []T, predicate func(T) bool) bool {
	for _, x := range xs {
		if !predicate(x) {
			return false
		}
	}
	return true
}

func TestSelect(buf []uint16, fn func(xs []uint16, k int), seed uint16, pow int, iters int) {
	_, file, line, _ := runtime.Caller(1)

	nextSeed := seed
	for p := 0; p <= pow; p++ {
		length := 1 << p
		xs := buf[:length]

		for i := 0; i < iters; i++ {
			initialSeed := nextSeed
			nextSeed = FillUint16Array(xs, initialSeed)
			nextSeed = XorShift16(nextSeed)
			k := int(float32(len(xs)) * float32(nextSeed) / float32(1<<16))
			nextSeed = XorShift16(nextSeed)
			// fmt.Printf("  seed:   %v\n", initialSeed)
			// fmt.Printf("  k:      %v of %v\n", k, len(xs))
			// fmt.Printf("  before: %v\n", xs)
			fn(xs, k)
			// fmt.Printf("  after:  %v\n", xs)
			partitionedLeft := All(xs[:k], func(x uint16) bool { return x <= xs[k] })
			partitionedRight := All(xs[k+1:], func(x uint16) bool { return x >= xs[k] })
			// fmt.Println("  sorted:", partitionedLeft, partitionedRight)
			if !partitionedLeft || !partitionedRight {
				fmt.Printf("%s failed for len=%d, seed=%d, k=%d at %s:%d\n", GetFunctionName(fn), length, initialSeed, k, file, line)
				return
			}
		}
	}
}

func main() {
	pow := 12
	buf := make([]uint16, 1<<pow)
	aux := make([]uint16, 1<<pow)

	// {
	// 	var length int = 16
	// 	var seed uint16 = 3509
	// 	xs := buf[:length]
	// 	FillUint16Array(xs, seed)
	// 	fmt.Printf("a: %v\n", xs)
	// 	TopDownMergeSort(xs, aux)
	// 	fmt.Printf("b: %v\n", xs)
	// 	Assert(IsSorted(xs), "not sorted")
	// }
	// return

	// TestSort(buf, InsertionSort, 1, pow, 100)
	// TestSort(buf, SelectionSort, 1, pow, 100)
	// TestSort(buf, InsertionSort2, 1, pow, 100)
	// TestSort(buf, BubbleSort, 1, pow, 100)
	// TestSort(buf, BubbleSort2, 1, pow, 100)
	// TestSort(buf, ShellSort, 1, pow, 100)
	// TestSort(buf, QuickSort, 1, pow, 10000)
	// TestSort(buf, NonRecursiveQuickSort, 1, pow, 10000)
	// TestSort(buf, HybridQuickSort, 1, pow, 10000)
	// TestSelect(buf, Select, 1, pow, 10000)
	// TestSelect(buf, NonRecursiveSelect, 1, pow, 10000)
	TestSort(buf, func(xs []uint16) { TopDownMergeSort(xs, aux) }, 1, pow, 10000)

	// TestSort(buf, CountSort, 1, pow, 100) // TODO:
	// TestSort(buf, BucketSort, 1, pow, 100) // TODO:
	// TestSort(buf, ThreeWayQuickSort, 1, pow, 10000) // TODO:
}
