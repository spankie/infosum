package usemaps

import (
	"fmt"
	"io"
	"sync"

	"github.com/spankie/infosum/results"
)

type mapComparator struct {
	chuksize int
}

// NewComparator returns a comparator that uses maps to compare two files
// @param chunksize specifies how many items to read from the files at a time for processing
func NewComparator(chunksize int) mapComparator {
	return mapComparator{chuksize: chunksize}
}

func (mc mapComparator) Compare(fileA, fileB io.ReadCloser) (*results.ComparisonResult, error) {
	// increase the capacity to reduce allocaton
	// using 1000000 increases the memory usage by about 100MB
	fileAMap := make(map[string]int64, 100000)
	fileBMap := make(map[string]int64, 100000)

	var keyCountA, keyCountB int64

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		mc.processCSV(fileA, &keyCountA, fileAMap)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		mc.processCSV(fileB, &keyCountB, fileBMap)
		wg.Done()
	}()

	wg.Wait()

	var distinctOverlapCount int64
	var totalOverlap int64

	for key := range fileBMap {
		// check if the value is in fileAMap
		if _, ok := fileAMap[key]; ok {
			distinctOverlapCount++
			totalOverlap += fileAMap[key] * fileBMap[key]
		}
	}

	return &results.ComparisonResult{
		KeyCountA:         uint64(keyCountA),
		KeyCountB:         uint64(keyCountB),
		DistinctKeyCountA: uint64(len(fileAMap)),
		DistinctKeyCountB: uint64(len(fileBMap)),
		DistinctOverlap:   uint64(distinctOverlapCount),
		TotalMaxOverlap:   uint64(totalOverlap),
	}, nil
}

func (mc mapComparator) processCSV(f1 io.ReadCloser, keyCount *int64, filemap map[string]int64) {
	err := results.StreamCSVInChunks(f1, mc.chuksize, func(chunk []string) {
		*keyCount += int64(len(chunk))
		for _, record := range chunk {
			if _, ok := filemap[record]; !ok {
				filemap[record] = 1
			} else {
				filemap[record] += 1
			}
		}
	})
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
	}
}
