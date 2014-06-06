package main

import (
	"fmt"
	"time"
	"math/rand"
	"math/big"
	"runtime"
)

const (
	MAX_VALUE = 9223372036854775807 // maximum 64-bit signed integer
	// N = 10
	N = 10000
)

//var (
//	// N int = 10000
//)

func generateList(count int) []int64 {
	ret := make([]int64, count)
	for i := 0; i < count; i++ {
		ret[i] = int64(rand.Int63()) // generate random positive 63-bit integer
	}
	return ret
}

func shuffleAndRemoveElement(array []int64) []int64 {
	arrayLen := len(array)
	randomIndexes := rand.Perm(arrayLen)
	//fmt.Printf("randomIndexes = %v\n", randomIndexes)

	retLen := arrayLen - 1
	ret := make([]int64, retLen)
	for i := 0; i < retLen; i++ {
		ret[i] = array[ randomIndexes[i] ]
	}
	return ret
}

func findMissingElement(first []int64, second []int64) int {
	sum := big.NewInt(0)
	firstLength := len(first)
	lastIndex := firstLength - 1
	secondArrayMap := make(map[int64]int, firstLength)

	for i, secondVal := range second {
		firstVal := first[i]
		diff := int64(firstVal - secondVal)

		// index the 'first array' value
		secondArrayMap[firstVal] = i

		sum = sum.Add(sum, big.NewInt(diff) )
		//fmt.Printf("%v| diff = %v - %v = %v, sum = %v\n", i, firstVal, secondVal, diff, sum)
	}

	// add the last value from first array
	lastVal := int64(first[lastIndex])
	// index the last value of the 'first array'
	secondArrayMap[lastVal] = lastIndex

	// finalise sum calculation
	sum = sum.Add(sum, big.NewInt(lastVal))
	// fmt.Printf("add last value from first array = %v\n", lastVal)

	// fmt.Printf("Final sum = %v\n", sum)
	//fmt.Printf("secondArrayMap = %v\n", secondArrayMap)

	return secondArrayMap[sum.Int64()]
}

func main() {
	// t0 := time.Now()

	// Create and seed the generator.
	// Typically a non-fixed seed should be used, such as time.Now().UnixNano().
	// Using a fixed seed will produce the same output on every run.
	rand.Seed(time.Now().UnixNano())

	// fmt.Printf("MAX_VALUE = %v\n", MAX_VALUE)

	// Prepare multi-core
	NUM_CPU := runtime.NumCPU()
	fmt.Printf("NumCPU = %v\n", NUM_CPU)
	runtime.GOMAXPROCS(NUM_CPU)

	first := generateList(N)
	// fmt.Printf("first = %v\n", first)

	second := shuffleAndRemoveElement(first)
	// fmt.Printf("second = %v\n", second)

	missingElementIndex := findMissingElement(first, second)
	fmt.Printf("Missing element is %v\n", missingElementIndex)

	// t1 := time.Now()
	// fmt.Printf("The call took %v to run.\n", t1.Sub(t0))
}