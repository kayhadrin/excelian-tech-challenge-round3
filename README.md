# Excelian Tech Challenge Round 3

## Challenge

### Introduction

As H L Mencken used to say, “For every complex problem there is an answer that is clear, simple, and wrong.”
Here, the problem looks simple but solving it efficiently needs a clever trick.

Hints: using a bitwise operator is advised.

### Challenge

There is an array of non-negative integers. A second array is formed by shuffling the elements of the first array and deleting a random element. Given these two arrays, find which element is missing in the second array.

The challenge is to write a small test program that will be as light as possible in terms of code and as fast as possible in terms of finding the solution. All the solutions will be ordered by number of characters used (=as light as possible) and performance (=as fast as possible) when run on a simple laptop (intel i5 2.6 Ghz and 8 GB RAM); 50% of the points will be attributed for each.

### Rules

The code should look like the following:

	// generate a list of random N Integers
	int N = 1000;
	List<Integer> first = generateList(N);
	// shuffle it and remove one element
	List<Integer> second = shuffleAnrRemoveElement(first);
	// call a function returning the missing element
	Integer missingElement = findMissingElement();
	System.out.println("Missing element is " + missingElement);

### Test

N (generateList()’s argument) will be set to a small value to check the algorithm is producing the expected result, and then set to 33554431 (2^25 - 1).

Good luck! :-) 

## Build

Assuming that you have a working Go development environment.

	cd <this git repo root>
	go install excelian/dhanszec/techChallenge/round3/multiCore
	go install excelian/dhanszec/techChallenge/round3/singleCore
	go install excelian/dhanszec/techChallenge/round3/multiCoreMin
	go install excelian/dhanszec/techChallenge/round3/singleCoreMin

	# run multi-core (code is using multiple CPUs by default)
	bin\multiCore.exe

	# run multi-core (minified) (code is using multiple CPUs by default)
	bin\multiCoreMin.exe

	# run single-core
	bin\singleCore.exe

	# run single-core (minified)
	bin\singleCoreMin.exe

## Performance test notes

Tested on:

- i7-2600K 3.7GHz
- 8GB RAM
- Win7 64b

Multi code bin version: