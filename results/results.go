package results

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

func validateRecord(record string) string {
	return strings.Trim(record, "\t ")
}

func StreamCSVInChunks(source io.ReadCloser, chunkSize int, processChunk func([]string)) error {
	defer source.Close()

	// Create a CSV reader
	reader := csv.NewReader(source)
	// reader.ReuseRecord = true
	reader.FieldsPerRecord = 1 //TODO: confirm what this does exactly

	_, err := reader.Read()
	if err != nil {
		return fmt.Errorf("error reading the header row: %w", err)
	}

	// Read the file in chunks
	chunk := make([]string, 0, chunkSize)
	for {
		chunk = chunk[:0] // Reset slice length, retain capacity
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
	KeyCountA         uint64
	KeyCountB         uint64
	DistinctKeyCountA uint64
	DistinctKeyCountB uint64
	DistinctOverlap   uint64
	TotalMaxOverlap   uint64
}

func (c ComparisonResult) Print(w io.Writer) {
	fmt.Fprintln(w)
	fmt.Fprintf(w, "Count of keys in file A: %v\n", c.KeyCountA)
	fmt.Fprintf(w, "Count of distinct keys file A: %v\n", c.DistinctKeyCountA)
	fmt.Fprintln(w)
	fmt.Fprintf(w, "Count of keys in file B: %v\n", c.KeyCountB)
	fmt.Fprintf(w, "Count of distinct keys in file B: %v\n", c.DistinctKeyCountB)
	// fmt.Println()
	// fmt.Printf("Count of empty keys in file B: %v\n", fileBMap[""])
	fmt.Fprintln(w)
	fmt.Fprintf(w, "Count of distinct overlap: %v\n", c.DistinctOverlap)
	fmt.Fprintf(w, "Count of total max overlap: %v\n", c.TotalMaxOverlap)
	fmt.Fprintln(w)
}
