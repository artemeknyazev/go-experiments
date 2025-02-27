package main

import (
	"math/bits"
)

// https://en.wikipedia.org/wiki/Euclidean_algorithm
func EuclidianGcd(a, b uint64) uint64 {
	if a == 0 {
		return b
	}
	if b == 0 {
		return a
	}
	if b > a {
		a, b = b, a
	}

	var q uint64 = 0
	var r uint64 = 1
	for r != 0 {
		q = a / b
		r = a - q*b
		a, b = b, r
	}
	return a
}

// https://en.wikipedia.org/wiki/Binary_GCD_algorithm
// 1: gcd(  u,   0) = u
// 2: gcd(2*u, 2*v) = 2 * gcd(u,   v)
// 3: gcd(  u, 2*v) =     gcd(u,   v) if u is odd
// 4: gcd(  u,   v) =     gcd(u, v-u) if u, v are odd, u <= v
func BinaryGcd(u, v uint64) uint64 {
	if u == 0 {
		return v // 1
	}
	if v == 0 {
		return u // 1
	}

	i := bits.TrailingZeros64(u)
	u >>= i
	j := bits.TrailingZeros64(v)
	v >>= j
	k := min(i, j) // 2

	for {
		if u > v {
			u, v = v, u // 4
		}
		v -= u // 4
		if v == 0 {
			return u << k // 1
		}
		v >>= bits.TrailingZeros64(v) // 3
	}

}

func main() {
	if false {
		println(EuclidianGcd(5, 3))
	}
	if false {
		println(BinaryGcd(15, 6))
	}
}
