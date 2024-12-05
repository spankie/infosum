package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/spankie/infosum/algorithms/bitset"
	"github.com/spankie/infosum/results"
)

func mustGetCSVFIle(filePath string) io.ReadCloser {
	// Open the CSV file
	f, err := os.Open(filePath)
	if err != nil {
		panic(fmt.Errorf("cannot find file %v: %w", filePath, err))
	}
	return f
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

	/*
		comparator := usemaps.NewComparator(chunksize)
		result, err := comparator.Compare(fileA, fileB)
		if err != nil {
			fmt.Printf("err getting result: %v", err)
			os.Exit(1)
		}
	*/

	// 0.00001, 0.8 works well with the given sample data
	comparator := bitset.NewComparator(*chunkSize, 0.00001, 0.8)
	result, err := comparator.Compare(fileA, fileB)
	if err != nil {
		panic(fmt.Errorf("err getting result: %v", err))
	}

	result.Print(os.Stdout)
}
