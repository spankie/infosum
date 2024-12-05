package maps

import (
	"fmt"
	"io"

	ds "github.com/spankie/infosum/dataset"
)

type dataset struct {
	resource io.Reader
	datamap  map[string]int
	count    uint
}

func newDataset(resource io.Reader) *dataset {
	// increase the capacity to reduce allocaton
	// using 1000000 increases the memory usage by about 100MB
	return &dataset{resource: resource}
}

func (d *dataset) processCSV(chunksize int) {
	if d.datamap == nil {
		d.datamap = make(map[string]int, chunksize)
	}
	err := ds.StreamCSVInChunks(d.resource, chunksize, func(chunk []string) {
		d.count += uint(len(chunk))
		for _, record := range chunk {
			if _, ok := d.datamap[record]; !ok {
				d.datamap[record] = 1
			} else {
				d.datamap[record] += 1
			}
		}
	})
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
	}
}
