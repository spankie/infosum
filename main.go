package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/spankie/infosum/algorithms/usebitset"
	"github.com/spankie/infosum/results"
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
	chunkSize := flag.Int("chunksize", 1000, "Size of the chunks to process at a time")
	filenameA := flag.String("fileA", "", "filename of the first file to compare")
	filenameB := flag.String("fileB", "", "filename of the second file to compare")
	flag.Parse()
	// Check if the correct number of arguments are passed
	if *filenameA == "" || *filenameB == "" {
		fmt.Println("Usage: go run main.go [--chunksize=1000] <file1.csv> <file2.csv>")
		return
	}

	// Print the file names
	fmt.Printf("File A: %s\n", *filenameA)
	fmt.Printf("File B: %s\n", *filenameB)

	fileA := mustGetCSVFIle(*filenameA)
	defer fileA.Close()
	fileB := mustGetCSVFIle(*filenameB)
	defer fileB.Close()

	var comparator Comparator

	/*
		comparator = usemaps.NewComparator(chunksize)
		result, err := comparator.Compare(fileA, fileB)
		if err != nil {
			fmt.Printf("err getting result: %v", err)
			os.Exit(1)
		}
	*/

	// 0.00001, 0.8 works well with the given sample data
	comparator = usebitset.NewComparator(*chunkSize, 0.00001, 0.8)
	result, err := comparator.Compare(fileA, fileB)
	if err != nil {
		fmt.Printf("err getting result: %v", err)
		os.Exit(1)
	}

	result.Print(os.Stdout)

	// printMemoryUsage()
}
