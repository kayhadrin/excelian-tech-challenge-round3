# Excelian Tech Challenge Round 3

# Build

Assuming that you have a working Go development environment.

	cd <this git repo root>
	go install excelian/dhanszec/techChallenge/round3/multiCore
	go install excelian/dhanszec/techChallenge/round3/singleCore
	go install excelian/dhanszec/techChallenge/round3/singleCoreMin

	# run multi-core (code is using multiple threads by default)
	bin\multiCore.exe

	# run single-core
	bin\singleCore.exe

	# run single-core (minified)
	bin\singleCoreMin.exe

## Performance test notes

Tested on:

- i7-2600K 3.7GHz
- 8GB RAM
- Win7 64b

| N | Single Core   | Multi threads | Multi core |  
| --|:--------------|:-------------:| ----------:|
| 10,000 | 3610 | 3600 | 8685 |
| 10,000,000 | 5291302 | 5443311 | 4788273 |
*Time in microseconds.*
