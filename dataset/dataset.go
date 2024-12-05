package dataset

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

func validateRecord(record string) string {
	return strings.Trim(record, "\t ")
}

func StreamCSVInChunks(source io.Reader, chunkSize int, processChunk func([]string)) error {
	// Create a CSV reader
	reader := csv.NewReader(source)

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
		n := chunkSize
		// using this kind of loop to make sure I get all the maximum chunk size
		for n > 0 {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("error reading CSV: %w", err)
			}
			// NOTE: this would be good to be in a go routine but it is trivial check so i'll leave it this way.
			key := validateRecord(record[0])
			if key == `` || key == `""` {
				continue
			}
			chunk = append(chunk, key)
			n--
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
