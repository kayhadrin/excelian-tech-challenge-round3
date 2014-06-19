// Multi-routine version
// It seems that it only becomes more efficient than the single-core version when N>10,000,000
package main

import (
	"fmt"
	"time"
	"math"
	"math/rand"
	//"math/big"
	"runtime"
)

const (
	MAX_VALUE = 9223372036854775807 // maximum 64-bit signed integer
	// MAX_VALUE = 100
	// N = 100
	// N = 10000000 // this is when multi-core starts to become noticeable
	// N = 10000
	N = 33554431 // math.Pow(2, 25) - 1
	// N = 1048575 // math.Pow(2, 20) - 1
	TEST_COUNT = 1
)

var (
	CPU_NUM int
	//randNbChannel chan []int64
	xorChannel chan int64
)

func generateList(count int) []int64 {
	ret := make([]int64, count)
	for i := 0; i < count; i++ {
		ret[i] = int64(rand.Int63()) // generate random positive 63-bit integer
		//ret[i] = int64(rand.Int63n(MAX_VALUE)) // generate random positive 63-bit integer
	}
	return ret
}

func shuffleAndRemoveElement(array []int64) []int64 {
	arrayLen := len(array)
	randomIndexes := rand.Perm(arrayLen)
	//fmt.Printf("randomIndexes = %v\n", randomIndexes)

	//DEBUG
	//fmt.Printf("CHEAT missing element index = %v\n", randomIndexes[len(randomIndexes) - 1])
	//fmt.Printf("CHEAT missing element value = %v\n", array[randomIndexes[len(randomIndexes) - 1]])

	retLen := arrayLen - 1
	ret := make([]int64, retLen)
	for i := 0; i < retLen; i++ {
		ret[i] = array[ randomIndexes[i] ]
	}
	return ret
}

func findMissingElement(first []int64, second []int64) int {
	//firstLength := len(first)
	secondLength := len(second)
	//////////////////////////////////////////////////////////////////////////
	defaultSliceCount := int(math.Floor(float64(secondLength / CPU_NUM)))

	//DEBUG
	// fmt.Printf("sliceCount = %v\n", defaultSliceCount)

	xorArrays := func (_first, _second []int64, async bool) int64 {
		var _sum int64
		_firstLen := len(_first)
		_secondLen := len(_second)

		_sum = _first[0]
		for i := 1; i < _firstLen; i++ {
			_sum ^= _first[i]
		}
		for i := 0; i < _secondLen; i++ {
			_sum ^= _second[i]
		}

		if async {
			// fmt.Printf("<- _sum = %v\n", _sum)
			xorChannel <- _sum
		//} else {
		//	fmt.Printf("return _sum = %v\n", _sum)
		}

		return _sum
	}

	// send off goroutines
	cpuNumMinus := CPU_NUM - 1
	offset := 0
	for i := 0; i < cpuNumMinus; i++ {
		go xorArrays(first[offset:offset + defaultSliceCount], second[offset:offset + defaultSliceCount], true)
		offset += defaultSliceCount
	}
	xorSum := xorArrays(first[offset:], second[offset:], false)

	//fmt.Printf("xorSum = %v\n", xorSum)

	// process received data
	for i := 0; i < cpuNumMinus; i++ {
		receivedXor := <- xorChannel
		//fmt.Printf("receivedXor = %v\n", receivedXor)
		xorSum ^= receivedXor
		//fmt.Printf("xorSum = %v\n", xorSum)
	}
	//////////////////////////////////////////////////////////////////////////
	for i, firstVal := range first {
		if firstVal == xorSum {
			return i
		}
	}

	return -1
}

func main() {
	//t0 := time.Now()

	// Create and seed the generator.
	// Typically a non-fixed seed should be used, such as time.Now().UnixNano().
	// Using a fixed seed will produce the same output on every run.
	rand.Seed(time.Now().UnixNano())

	//fmt.Printf("N = %v\n", N)
	//fmt.Printf("MAX_VALUE = %v\n", MAX_VALUE)

	// Prepare multi-core
	CPU_NUM = runtime.NumCPU()
	//////////////////////
	//DEBUG
	// Enable the line below to allow Go to run goroutines on multiple CPU cores.
	// But it's not very effiecient unless N is really large (10M)
	runtime.GOMAXPROCS(CPU_NUM)
	//////////////////////

	//fmt.Printf("CPU_NUM = %v\n", CPU_NUM)

	// prepare goroutines communication channel
	//randNbChannel = make(chan []int64)
	xorChannel = make(chan int64)

	for i:=0; i < TEST_COUNT; i++ {

		// generateList(N)
		first := generateList(N)
		// fmt.Printf("first = %v\n", first)

		second := shuffleAndRemoveElement(first)
		// fmt.Printf("second = %v\n", second)

		missingElementIndex := findMissingElement(first, second)
		fmt.Printf("Missing element is %v\n", missingElementIndex)

		//fmt.Printf("end loop[%v] ----------\n", i)
	}

	/*
	t1 := time.Now()
	duration := t1.Sub(t0)
	averageTime := duration.Nanoseconds() / int64(TEST_COUNT) / 1000 // in microseconds
	fmt.Printf("The call took %v to run. Average: %vus\n", duration, averageTime)
	*/
}
