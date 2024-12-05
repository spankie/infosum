BIG_FILE_A=data/bigger/A_f_2.csv
BIG_FILE_B=data/bigger/B_f_2.csv

FILE_A=data/A_f.csv
FILE_B=data/B_f.csv

BINARY_NAME=infosum

test:
	@go test -v ./...

lint:
	@golangci-lint run

bench:
	@go test -bench=. -benchmem -benchtime=4s ./...

build:
	@go build -o $(BINARY_NAME)

run: build
	@./$(BINARY_NAME) --chunksize=2000 --fileA=$(FILE_A) --fileB=$(FILE_B)

bench-prof-bs:
	@go test -bench=. -benchmem -count=1 -benchtime=3s ./bitset -cpuprofile=cpu_bs.prof -memprofile=mem_bs.prof
	# @go tool pprof -http=:8080 mem_bs.prof

bench-prof-hm:
	@go test -bench=. -benchmem -count=1 -benchtime=3s ./hashmaps -cpuprofile=cpu_hm.prof -memprofile=mem_hm.prof
	# @go tool pprof -http=:8080 mem_hm.prof
	# @go test -bench=. -benchmem -benchtime=1x ./hashmaps -cpuprofile=cpu_hm.prof -memprofile=mem_hm.prof
