package results

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func validateRecord(record string) string {
	return strings.Trim(record, "\t ")
}

func StreamCSVInChunks(source io.Reader, chunkSize int, processChunk func([]string)) error {
	// defer source.Close()

	// Create a CSV reader
	reader := csv.NewReader(source)
	// reader.ReuseRecord = true
	reader.FieldsPerRecord = 1 //TODO: confirm what this does exactly

	// read the header from the csv
	_, err := reader.Read()
	if err != nil {
		return fmt.Errorf("error reading the header row: %w", err)
	}

	// Read the file in chunks
	chunk := make([]string, 0, chunkSize)
	for {
		// reuse the chunk by Reseting slice length and retain capacity
		chunk = chunk[:0]
		// Read up to chunkSize lines
		for i := 0; i < chunkSize; i++ {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("error reading CSV: %w", err)
			}
			key := validateRecord(record[0])
			if key == "" {
				continue
			}
			chunk = append(chunk, key)
		}

		// If the chunk is empty, we're done
		if len(chunk) == 0 {
			break
		}

		// Process the chunk
		processChunk(chunk)
	}

	return nil
}

type ComparisonResult struct {
	KeyCountA         uint
	KeyCountB         uint
	DistinctKeyCountA uint
	DistinctKeyCountB uint
	DistinctOverlap   uint
	TotalMaxOverlap   uint
}

type KeyCount struct {
	KeyCount         uint `json:"key_count"`
	DistinctKeyCount uint `json:"distinct_key_count"`
}

type Overlap struct {
	DistinctOverlap uint `json:"distinct_overlap"`
	TotalMaxOverlap uint `json:"total_max_overlap"`
}

func (c ComparisonResult) Print(w io.Writer) error {
	result := struct {
		FileA   KeyCount `json:"file_a"`
		FileB   KeyCount `json:"file_b"`
		Overlap Overlap  `json:"overlap"`
	}{
		FileA: KeyCount{
			KeyCount:         c.KeyCountA,
			DistinctKeyCount: c.DistinctKeyCountA,
		},
		FileB: KeyCount{
			KeyCount:         c.KeyCountB,
			DistinctKeyCount: c.DistinctKeyCountB,
		},
		Overlap: Overlap{
			DistinctOverlap: c.DistinctOverlap,
			TotalMaxOverlap: c.TotalMaxOverlap,
		},
	}
	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling json response: %w", err)
	}
	_, err = fmt.Fprintln(w, string(resultJSON))
	if err != nil {
		return fmt.Errorf("error printing to resource: %w", err)
	}

	return nil
}
