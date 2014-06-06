package main

import (
	"fmt"
	"time"
	"math/rand"
	"math/big"
)

func generateList(c int) []int64 {
	r := make([]int64, c)
	for i := 0; i < c; i++ {
		r[i] = int64(rand.Int63())
	}
	return r
}

func shuffleAndRemoveElement(a []int64) []int64 {
	l := len(a)
	n := rand.Perm(l)
	t := l - 1
	e := make([]int64, t)
	for i := 0; i < t; i++ {
		e[i] = a[ n[i] ]
	}
	return e
}

func findMissingElement(f []int64, s []int64) int {
	m := big.NewInt(0)
	a := len(f)
	b := a - 1
	c := make(map[int64]int, a)

	for i, d := range s {
		e := f[i]
		g := int64(e - d)
		c[e] = i
		m = m.Add(m, big.NewInt(g) )
	}
	
	o := int64(f[b])
	c[o] = b
	m = m.Add(m, big.NewInt(o))
	return c[m.Int64()]
}

func main() {
	rand.Seed(time.Now().UnixNano())
	N := 10000
	first := generateList(N)	
	second := shuffleAndRemoveElement(first)
	missingElementIndex := findMissingElement(first, second)
	fmt.Printf("Missing element is %v\n", missingElementIndex)
}
