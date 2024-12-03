package usebitset

import (
	"fmt"
	"io"
	"strconv"

	"github.com/bits-and-blooms/bitset"
	"github.com/shenwei356/countminsketch"
	cms "github.com/shenwei356/countminsketch"
	"github.com/spankie/infosum/results"
)

type dataset struct {
	resource  io.Reader
	bitset    *bitset.BitSet
	cms       *cms.CountMinSketch
	totalKeys int64
}

// epsilon, delta := 0.00001, 0.8 OR (D = 3, W = 200,000) // values that works well for the tasks files
// epsilon, delta := 0.0000001, 0.9 OR (D = 3, W = 20,000,000) // values that works well for the generated files
func newDataset(resource io.Reader, epsilon, delta float64) (*dataset, error) {
	cms, err := countminsketch.NewWithEstimates(epsilon, delta)
	if err != nil {
		return nil, err
	}
	// fmt.Printf("ε: %f, δ: %f -> d: %d, w: %d\n", epsilon, delta, cmsA.D(), cmsA.W())

	return &dataset{
		resource: resource,
		bitset:   bitset.New(100),
		cms:      cms,
	}, nil
}

func (d *dataset) distinctCount() uint {
	return d.bitset.Count()
}

func (d *dataset) processCSVFile(chunksize int) {
	err := results.StreamCSVInChunks(d.resource, chunksize, func(records []string) {
		// Process each record
		d.totalKeys += int64(len(records))
		for _, record := range records {
			key, err := strconv.Atoi(record)
			if err != nil {
				fmt.Printf("Error parsing key: %v\n", err)
				continue
			}

			// update the countminsketch for the current record
			d.cms.UpdateString(record, 1)

			// Update bitset
			d.bitset.Set(uint(key))
		}
	})
	if err != nil {
		fmt.Printf("unable to process file: %v", err)
	}
}
