package results

import (
	"encoding/json"
	"fmt"
	"io"
)

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
