package bitset

import (
	"io"
	"strconv"
	"sync"

	"github.com/bits-and-blooms/bitset"
	"github.com/shenwei356/countminsketch"
	"github.com/spankie/infosum/results"
)

type bitsetComparator struct {
	chunksize      int
	epsilon, delta float64
}

func NewComparator(chunksize int, epsilon, delta float64) bitsetComparator {
	return bitsetComparator{
		chunksize: chunksize,
		epsilon:   epsilon,
		delta:     delta,
	}
}

func (bsc bitsetComparator) Compare(fileA, fileB io.Reader) (*results.ComparisonResult, error) {
	datasetA, err := newDataset(fileA, bsc.epsilon, bsc.delta)
	if err != nil {
		return nil, err
	}

	datasetB, err := newDataset(fileB, bsc.epsilon, bsc.delta)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		datasetA.processCSVFile(bsc.chunksize)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		datasetB.processCSVFile(bsc.chunksize)
		wg.Done()
	}()

	wg.Wait()

	result := &results.ComparisonResult{
		KeyCountA:         uint(datasetA.totalKeys),
		KeyCountB:         uint(datasetB.totalKeys),
		DistinctKeyCountA: uint(datasetA.distinctCount()),
		DistinctKeyCountB: uint(datasetB.distinctCount()),
	}
	result.DistinctOverlap, result.TotalMaxOverlap = bsc.distinctOverlap(datasetA, datasetB)

	return result, nil
}

func (bsc bitsetComparator) distinctOverlap(d1, d2 *dataset) (uint, uint) {
	intersection := d1.bitset.Intersection(d2.bitset)
	distinctOverlap := intersection.Count()
	totalMaxOveralp := getTotalOverlapMany(intersection, d1.cms, d2.cms)
	return uint(distinctOverlap), uint(totalMaxOveralp)
}

func getTotalOverlapMany(bs *bitset.BitSet, cmsA, cmsB *countminsketch.CountMinSketch) int64 {
	var totalOverlap int64

	buffer := make([]uint, 1024)
	j := uint(0)
	j, buffer = bs.NextSetMany(j, buffer)

	// for each set bit, get the number converted to string and check for it's aproximated frequency
	// in the countminsketch
	for ; len(buffer) > 0; j, buffer = bs.NextSetMany(j, buffer) {
		for _, v := range buffer {
			record := strconv.Itoa(int(v))
			totalOverlap += int64(cmsA.EstimateString(record) * cmsB.EstimateString(record))
		}
		j += 1
	}

	return totalOverlap
}
