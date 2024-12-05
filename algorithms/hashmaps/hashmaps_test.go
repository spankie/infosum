package maps

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/spankie/infosum/results"
)

/*
	// actual result with empty spaces
	expected := results.ComparisonResult{
		KeyCountA:         95_000,
		KeyCountB:         80_000,
		DistinctKeyCountA: 72_798,
		DistinctKeyCountB: 72_814,
		TotalMaxOverlap:   69_272, //60_627_882,
		DistinctOverlap:   58_221,
	}
*/

func BenchmarkCompareWithMaps(b *testing.B) {
	fileA, fileB := "../../data/A_f.csv", "../../data/B_f.csv"

	fA, err := os.Open(fileA)
	if err != nil {
		b.Fatalf("error reading file %s: error: %v", fileA, err)
	}
	b.Cleanup(func() { fA.Close() })

	fB, err := os.Open(fileB)
	if err != nil {
		b.Fatalf("error reading file %s: error: %v", fileA, err)
	}
	b.Cleanup(func() { fB.Close() })

	comparator := NewComparator(1000)

	for n := 0; n < b.N; n++ {
		// Reset the file position to the beginning
		_, err := fA.Seek(0, io.SeekStart)
		if err != nil {
			b.Fatalf("failed to reset file position: %v", err)
		}

		_, err = fB.Seek(0, io.SeekStart)
		if err != nil {
			b.Fatalf("failed to reset file position: %v", err)
		}

		_, err = comparator.Compare(io.NopCloser(fA), io.NopCloser(fB))
		if err != nil {
			b.Fatalf("expected nil error but got : %v", err)
		}
	}
}

func TestCompareWithMaps(t *testing.T) {
	comparator := NewComparator(1000)
	// Dataset 1: A B C D D E F F
	// Dataset 2: A C C D F F F X Y
	t.Run("using small dataset representing the tasks example", func(t *testing.T) {
		datasetA := `udprn
		12345
		12346

		12347
		12348
		12348
		12349
		12340

		12340`
		datasetB := `udprn
		12345
		12347
		12347
		12348
		12340
		12340
		12340
		12400

		12401`

		fileA := bytes.NewReader([]byte(datasetA))
		fileB := bytes.NewReader([]byte(datasetB))

		result, err := comparator.Compare(io.NopCloser(fileA), io.NopCloser(fileB))
		if err != nil {
			t.Fatalf("expected error to be nil but got: %v", err)
		}
		expected := &results.ComparisonResult{
			KeyCountA:         8,
			KeyCountB:         9,
			DistinctKeyCountA: 6,
			DistinctKeyCountB: 6,
			TotalMaxOverlap:   11,
			DistinctOverlap:   4,
		}
		assertResult(t, expected, result)
	})

	t.Run("using the actual tasks files", func(t *testing.T) {
		filenameA, filenameB := "../../data/A_f.csv", "../../data/B_f.csv"

		fileA, err := os.Open(filenameA)
		if err != nil {
			t.Fatalf("error reading file %s: error: %v", filenameA, err)
		}
		t.Cleanup(func() { fileA.Close() })

		fileB, err := os.Open(filenameB)
		if err != nil {
			t.Fatalf("error reading file %s: error: %v", filenameA, err)
		}
		t.Cleanup(func() { fileB.Close() })

		result, err := comparator.Compare(io.NopCloser(fileA), io.NopCloser(fileB))
		if err != nil {
			t.Fatalf("expected error to be nil but got: %v", err)
		}

		expected := &results.ComparisonResult{
			KeyCountA:         86_535,
			KeyCountB:         72_846,
			DistinctKeyCountA: 72_798,
			DistinctKeyCountB: 72_814,
			TotalMaxOverlap:   69_272, //60_627_882,
			DistinctOverlap:   58_221,
		}

		assertResult(t, expected, result)
	})
}

func assertResult(t *testing.T, expected, actual *results.ComparisonResult) {
	if actual.KeyCountA != expected.KeyCountA {
		t.Errorf("key count for file A should be %d; got %d", expected.KeyCountA, actual.KeyCountA)
	}
	if actual.KeyCountB != expected.KeyCountB {
		t.Errorf("key count for file A should be %d; got %d", expected.KeyCountB, actual.KeyCountB)
	}

	if actual.DistinctKeyCountA != expected.DistinctKeyCountA {
		t.Errorf("distinct count of file A should be %d but got %d", expected.DistinctKeyCountA, actual.DistinctKeyCountA)
	}
	if actual.DistinctKeyCountB != expected.DistinctKeyCountB {
		t.Errorf("distinct count of file B should be %d but got %d", expected.DistinctKeyCountB, actual.DistinctKeyCountB)
	}

	if actual.TotalMaxOverlap != expected.TotalMaxOverlap {
		t.Errorf("total overlap should be %d but got %d", expected.TotalMaxOverlap, actual.TotalMaxOverlap)
	}
	if actual.DistinctOverlap != expected.DistinctOverlap {
		t.Errorf("Distinct overlap should be %d but got %d", expected.DistinctOverlap, actual.DistinctOverlap)
	}
}
