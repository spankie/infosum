package usebitset

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

/*
func generateRandomNumber(digits int) int {
	if digits < 7 || digits > 8 {
		panic("Only 7 or 8 digit numbers are supported")
	}

	// Compute the range for the number
	min := int(1e6) // 7-digit starts from 1,000,000
	if digits == 8 {
		min = int(1e7) // 8-digit starts from 10,000,000
	}
	max := min*10 - 1 // Max value for 7 or 8 digits (e.g., 9,999,999 or 99,999,999)

	return rand.Intn(max-min+1) + min
}

func main() {
	f, err := os.OpenFile("data/bigger/B_f_2.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	for i := 0; i < 10_000_000; i++ {
		num := generateRandomNumber(7)
		// log.Printf("%d\n", num)
		_, err := f.WriteString(fmt.Sprintf("%d\n", num))
		if err != nil {
			fmt.Printf("error writing to file: %v", err)
			return
		}
	}
}
*/
