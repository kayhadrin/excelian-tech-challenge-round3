package main

import (
	"fmt"
	"time"
	"math"
	"math/rand"
	"math/big"
	"runtime"
)

const (
	M = 9223372036854775807
	N = 33554431
)

var (
	C int
	s chan *big.Int
)

func generateList(c int) []int64 {
	r := make([]int64, c)
	for i := 0; i < c; i++ {
		r[i] = int64(rand.Int63())
	}
	return r
}

func shuffleAndRemoveElement(a []int64) []int64 {
	L := len(a)
	r := rand.Perm(L)

	n := L - 1
	t := make([]int64, n)
	for i := 0; i < n; i++ {
		t[i] = a[ r[i] ]
	}
	return t
}

func findMissingElement(a []int64, b []int64) int {
	c := len(a)
	d := len(b)
	e := make(map[int64]int, c)
	f := int(math.Floor(float64(d / C)))
	g := func (h, k []int64, l bool) *big.Int {
		m := big.NewInt(0)
		n := len(h)
		o := len(k)
		for i := 0; i < o; i++ {
			m = m.Add(m, big.NewInt(int64(h[i] - k[i])))
		}
		if n > o {
			m = m.Add(m, big.NewInt(h[n - 1]))
		}

		if l {
			s <- m
		}

		return m
	}

	p := C - 1
	q := 0
	for i := 0; i < p; i++ {
		go g(a[q:q + f], b[q:q + f], true)
		q += f
	}
	r := g(a[q:], b[q:], false)

	for i := 0; i < p; i++ {
		t := <- s
		r = r.Add(r, t)
	}

	for i, u := range a {
		e[u] = i
	}

	return e[r.Int64()]
}

func main() {
	rand.Seed(time.Now().UnixNano())
	C = runtime.NumCPU()
	runtime.GOMAXPROCS(C)
	s = make(chan *big.Int)
	F := generateList(N)
	S := shuffleAndRemoveElement(F)
	x := findMissingElement(F, S)
	fmt.Printf("Missing element is %v", x)
}
