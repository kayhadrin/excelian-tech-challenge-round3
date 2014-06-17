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
	MAX_VALUE = 9223372036854775807
	N = 33554431
	TEST_COUNT = 4
)

var (
	CPU_NUM int
	sumChannel chan *big.Int
)

func generateList(count int) []int64 {
	ret := make([]int64, count)
	for i := 0; i < count; i++ {
		ret[i] = int64(rand.Int63())

	}
	return ret
}

func shuffleAndRemoveElement(array []int64) []int64 {
	arrayLen := len(array)
	randomIndexes := rand.Perm(arrayLen)

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
	defaultSliceCount := int(math.Floor(float64(secondLength / CPU_NUM)))
	sumArrays := func (_first, _second []int64, async bool) *big.Int {
		_sum := big.NewInt(0)
		_firstLen := len(_first)
		_secondLen := len(_second)
		for i := 0; i < _secondLen; i++ {
			_sum = _sum.Add(_sum, big.NewInt(int64(_first[i] - _second[i])))
		}
		if _firstLen > _secondLen {
			_sum = _sum.Add(_sum, big.NewInt(_first[_firstLen - 1]))
		}

		if async {
			sumChannel <- _sum
		}

		return _sum
	}

	cpuNumMinus := CPU_NUM - 1
	offset := 0
	for i := 0; i < cpuNumMinus; i++ {
		go sumArrays(first[offset:offset + defaultSliceCount], second[offset:offset + defaultSliceCount], true)
		offset += defaultSliceCount
	}
	sum := sumArrays(first[offset:], second[offset:], false)

	for i := 0; i < cpuNumMinus; i++ {
		receivedSum := <- sumChannel
		sum = sum.Add(sum, receivedSum)
	}

	for i, firstVal := range first {
		firstArrayMap[firstVal] = i
	}

	return firstArrayMap[sum.Int64()]
}

func main() {
	t0 := time.Now()
	rand.Seed(time.Now().UnixNano())
	fmt.Printf("N = %v\n", N)
	fmt.Printf("MAX_VALUE = %v\n", MAX_VALUE)
	CPU_NUM = runtime.NumCPU()
	runtime.GOMAXPROCS(CPU_NUM)
	fmt.Printf("CPU_NUM = %v\n", CPU_NUM)
	sumChannel = make(chan *big.Int)
	for i:=0; i < TEST_COUNT; i++ {
		first := generateList(N)
		second := shuffleAndRemoveElement(first)
		missingElementIndex := findMissingElement(first, second)
		fmt.Printf("Missing element is %v\n", missingElementIndex)
		fmt.Printf("end loop[%v] ----------\n", i)
	}
	t1 := time.Now()
	duration := t1.Sub(t0)
	averageTime := duration.Nanoseconds() / int64(TEST_COUNT) / 1000
	fmt.Printf("The call took %v to run. Average: %vus\n", duration, averageTime)
}
