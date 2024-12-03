package main

import (
	"fmt"
	"io"
	"os"
	"runtime"

	_ "net/http/pprof"

	"github.com/spankie/infosum/results"
	"github.com/spankie/infosum/usebitset"
)

func mustGetCSVFIle(filePath string) io.ReadCloser {
	// Open the CSV file
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("cannot find file %v\n", filePath)
		os.Exit(1)
	}
	return f
}

func printMemoryUsage() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	fmt.Printf("Allocated memory: %v MB\n", memStats.Alloc/1024/1024.00)
	fmt.Printf("Total allocated: %v MB\n", memStats.TotalAlloc/1024/1024.00)
	fmt.Printf("Heap allocated: %v MB\n", memStats.HeapAlloc/1024/1024.00)
}

type Comparator interface {
	Compare(resource1, resource2 io.Reader) (*results.ComparisonResult, error)
}

func main() {
	// Check if the correct number of arguments are passed
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <file1.csv> <file2.csv>")
		return
	}

	// get the filenames from the arguments
	filenameA := os.Args[1]
	filenameB := os.Args[2]

	// Print the file names
	fmt.Printf("File A: %s\n", filenameA)
	fmt.Printf("File B: %s\n", filenameB)

	fileA := mustGetCSVFIle(filenameA)
	defer fileA.Close()
	fileB := mustGetCSVFIle(filenameB)
	defer fileB.Close()

	var comparator Comparator
	chunksize := 1000

	/*
		comparator = usemaps.NewComparator(chunksize)
		result, err := comparator.Compare(fileA, fileB)
		if err != nil {
			fmt.Printf("err getting result: %v", err)
			os.Exit(1)
		}
	*/

	// 0.00001, 0.8 works well with the given sample data
	comparator = usebitset.NewComparator(chunksize, 0.00001, 0.8)
	result, err := comparator.Compare(fileA, fileB)
	if err != nil {
		fmt.Printf("err getting result: %v", err)
		os.Exit(1)
	}

	result.Print(os.Stdout)

	// printMemoryUsage()
}
