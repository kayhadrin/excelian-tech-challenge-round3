// Multi-routine version
// It seems that it only becomes more efficient than the single-core version when N>10,000,000
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
	MAX_VALUE = 9223372036854775807 // maximum 64-bit signed integer
	// MAX_VALUE = 100
	// N = 100
	// N = 10000000 // this is when multi-core starts to become noticeable
	// N = 10000
	N = 33554431 // math.Pow(2, 25) - 1
	// N = 1048575 // math.Pow(2, 20) - 1
	TEST_COUNT = 4
)

var (
	CPU_NUM int
	//randNbChannel chan []int64
	sumChannel chan *big.Int
)

func generateList(count int) []int64 {
	ret := make([]int64, count)
	for i := 0; i < count; i++ {
		ret[i] = int64(rand.Int63()) // generate random positive 63-bit integer
		//ret[i] = int64(rand.Int63n(MAX_VALUE)) // generate random positive 63-bit integer
	}
	return ret
}
/*
*/

/*
func generateList(count int) []int64 {
	defaultSliceCount := int(math.Floor(float64(count / CPU_NUM)))
	finalSliceCount := count - defaultSliceCount * CPU_NUM + defaultSliceCount

	//DEBUG
	// fmt.Printf("sliceCount = %v\n", defaultSliceCount)
	// fmt.Printf("finalSliceCount = %v\n", finalSliceCount)

	// generate a sub-array of random numbers
	genRandNumbers := func (len int) {
		numbers := make([]int64, len)
		for j := 0; j < len; j++ {
			numbers[j] = rand.Int63n(MAX_VALUE)
		}
		randNbChannel <- numbers
	}

	finalArray := make([]int64, count)

	// send off goroutines
	cpuNumMinus := CPU_NUM - 1
	for i := 0; i < cpuNumMinus; i++ {
		go genRandNumbers(defaultSliceCount)
	}
	go genRandNumbers(finalSliceCount)

	// process received data
	finalIndex := 0
	for i := 0; i < CPU_NUM; i++ {
		newNumbers := <- randNbChannel
		// fmt.Printf("newNumbers = %v\n", newNumbers)
		// finalArray[i] = newNumbers
		// finalArray = append(finalArray, newNumbers...)

		for _, val := range newNumbers {
			finalArray[finalIndex] = val
			finalIndex++
		}
	}
	// fmt.Printf("done loop: finalArray = %v\n", finalArray)

	return finalArray
//
//	ret := make([]int64, count)
//	for i := 0; i < count; i++ {
//		ret[i] = int64(rand.Int63()) // generate random positive 63-bit integer
//	}
//	return ret
}
*/

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
	firstLength := len(first)
	secondLength := len(second)
	firstArrayMap := make(map[int64]int, firstLength)
	//////////////////////////////////////////////////////////////////////////
	defaultSliceCount := int(math.Floor(float64(secondLength / CPU_NUM)))

	//DEBUG
	// fmt.Printf("sliceCount = %v\n", defaultSliceCount)

	sumArrays := func (_first, _second []int64, async bool) *big.Int {
		_sum := big.NewInt(0)
		_firstLen := len(_first)
		_secondLen := len(_second)

		// fmt.Printf("_first = %v\n", _first)
		// fmt.Printf("_second = %v\n", _second)

		for i := 0; i < _secondLen; i++ {
			_sum = _sum.Add(_sum, big.NewInt(int64(_first[i] - _second[i])))
			//fmt.Printf("%v| diff = %v - %v = %v, sum = %v\n", i, firstVal, secondVal, diff, sum)
		}
		if _firstLen > _secondLen {
			// fmt.Printf("last item: %v\n", _first[_firstLen - 1])
			_sum = _sum.Add(_sum, big.NewInt(_first[_firstLen - 1]))
		}

		if async {
			//			fmt.Printf("<- _sum = %v\n", _sum)
			sumChannel <- _sum
		} else {
			//			fmt.Printf("return _sum = %v\n", _sum)
		}

		return _sum
	}

	// send off goroutines
	cpuNumMinus := CPU_NUM - 1
	offset := 0
	for i := 0; i < cpuNumMinus; i++ {
		go sumArrays(first[offset:offset + defaultSliceCount], second[offset:offset + defaultSliceCount], true)
		offset += defaultSliceCount
	}
	sum := sumArrays(first[offset:], second[offset:], false)

	//fmt.Printf("sum = %v\n", sum)

	// process received data
	for i := 0; i < cpuNumMinus; i++ {
		receivedSum := <- sumChannel
		//fmt.Printf("receivedSum = %v\n", receivedSum)
		sum = sum.Add(sum, receivedSum)
		//fmt.Printf("sum = %v\n", sum)
	}
	//////////////////////////////////////////////////////////////////////////

	for i, firstVal := range first {
		firstArrayMap[firstVal] = i
	}

	//	for i, secondVal := range second {
	//		firstVal := first[i]
	//		//diff := int64(firstVal - secondVal)
	//
	//		// index the 'first array' value
	//		firstArrayMap[firstVal] = i
	//
	//		//sum = sum.Add(sum, big.NewInt(diff) )
	//		//fmt.Printf("%v| diff = %v - %v = %v, sum = %v\n", i, firstVal, secondVal, diff, sum)
	//	}
	//
	//	// add the last value from first array
	//	lastVal := int64(first[lastIndex])
	//	// index the last value of the 'first array'
	//	firstArrayMap[lastVal] = lastIndex
	//
	//	// finalise sum calculation
	//	sum = sum.Add(sum, big.NewInt(lastVal))
	// fmt.Printf("add last value from first array = %v\n", lastVal)

	// fmt.Printf("Final sum = %v\n", sum)
	//fmt.Printf("firstArrayMap = %v\n", firstArrayMap)

	return firstArrayMap[sum.Int64()]
}

func main() {
	t0 := time.Now()

	// Create and seed the generator.
	// Typically a non-fixed seed should be used, such as time.Now().UnixNano().
	// Using a fixed seed will produce the same output on every run.
	rand.Seed(time.Now().UnixNano())

	fmt.Printf("N = %v\n", N)
	fmt.Printf("MAX_VALUE = %v\n", MAX_VALUE)

	// Prepare multi-core
	CPU_NUM = runtime.NumCPU()
	//////////////////////
	//DEBUG
	// Enable the line below to allow Go to run goroutines on multiple CPU cores.
	// But it's not very effiecient unless N is really large (10M)
	runtime.GOMAXPROCS(CPU_NUM)
	//////////////////////

	fmt.Printf("CPU_NUM = %v\n", CPU_NUM)

	// prepare goroutines communication channel
	//randNbChannel = make(chan []int64)
	sumChannel = make(chan *big.Int)

	for i:=0; i < TEST_COUNT; i++ {

		// generateList(N)
		first := generateList(N)
		// fmt.Printf("first = %v\n", first)

		second := shuffleAndRemoveElement(first)
		// fmt.Printf("second = %v\n", second)

		missingElementIndex := findMissingElement(first, second)
		fmt.Printf("Missing element is %v\n", missingElementIndex)

		fmt.Printf("end loop[%v] ----------\n", i)
	}

	t1 := time.Now()
	duration := t1.Sub(t0)
	averageTime := duration.Nanoseconds() / int64(TEST_COUNT) / 1000 // in microseconds
	fmt.Printf("The call took %v to run. Average: %vus\n", duration, averageTime)
}
