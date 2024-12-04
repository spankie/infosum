package usemaps

import (
	"io"
	"sync"

	"github.com/spankie/infosum/results"
)

type comparator struct {
	chuksize int
}

// NewComparator returns a comparator that uses maps to compare two files
// @param chunksize specifies how many items to read from the files at a time for processing
func NewComparator(chunksize int) comparator {
	return comparator{chuksize: chunksize}
}

func (mc comparator) Compare(fileA, fileB io.Reader) (*results.ComparisonResult, error) {
	datasetA := newDataset(fileA)
	datasetB := newDataset(fileB)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		datasetA.processCSV(mc.chuksize)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		datasetB.processCSV(mc.chuksize)
		wg.Done()
	}()

	wg.Wait()

	var distinctOverlapCount, totalOverlap int = getDistinctData(datasetA, datasetB)

	return &results.ComparisonResult{
		KeyCountA:         datasetA.count,
		KeyCountB:         datasetB.count,
		DistinctKeyCountA: uint(len(datasetA.datamap)),
		DistinctKeyCountB: uint(len(datasetB.datamap)),
		DistinctOverlap:   uint(distinctOverlapCount),
		TotalMaxOverlap:   uint(totalOverlap),
	}, nil
}

func getDistinctData(dA, dB *dataset) (int, int) {
	var distinctOverlapCount int
	var totalMaxOverlap int

	for key := range dB.datamap {
		// check if the value is in fileAMap
		if _, ok := dA.datamap[key]; ok {
			distinctOverlapCount++
			totalMaxOverlap += dA.datamap[key] * dB.datamap[key]
		}
	}

	return distinctOverlapCount, totalMaxOverlap
}
