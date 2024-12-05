package dataset

import (
	"encoding/csv"
	"os"
	"testing"
)

func TestValidate(t *testing.T) {
	testcases := []struct {
		name     string
		record   string
		expected string
	}{
		{
			name:     "should validate tabs and spaces correctly",
			record:   "\t 12345",
			expected: "12345",
		},
		{
			name:     "should validate empty string correctly",
			record:   "",
			expected: "",
		},
		{
			name:     "should validate multiple spaces correctly",
			record:   "  ",
			expected: "",
		},
		{
			name:     "should validate multiple tabs correctly",
			record:   "\t\t\t\t\t12345",
			expected: "12345",
		},
		{
			name:     "should validate multiple tabs and spaces correctly",
			record:   "\t\t\t  \t   \t12345",
			expected: "12345",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			v := validateRecord(tc.record)
			if v != tc.expected {
				t.Errorf("validating %q should be %q but got %q", tc.record, tc.expected, v)
			}
		})
	}
}

func createTempCSV(t *testing.T, data [][]string) (string, error) {
	t.Helper()
	tempFile, err := os.CreateTemp("", "tempfile-*.csv")
	if err != nil {
		return "", err
	}

	writer := csv.NewWriter(tempFile)
	err = writer.WriteAll(data)
	if err != nil {
		tempFile.Close()
		os.Remove(tempFile.Name())

		return "", err
	}

	writer.Flush()
	tempFile.Close()

	t.Cleanup(func() {
		os.Remove(tempFile.Name())
	})

	return tempFile.Name(), nil
}

func TestStreamCSVInChunks(t *testing.T) {
	dataset := [][]string{
		{"udprn"},
		{"12345"},
		{"12347"},
		{"12347"},
		{"12348"},
		{`""`},
		{"12340"},
		{"12340"},
		{"12340"},
		{"12400"},
		{"12401"},
	}

	filename, err := createTempCSV(t, dataset)
	if err != nil {
		t.Errorf("expected nil error creating csv file but got : %v", err)
	}

	file, err := os.Open(filename)
	if err != nil {
		t.Errorf("expected nil error opening csv file but got : %v", err)
	}
	t.Cleanup(func() { file.Close() })

	expectedChunks := [][]string{
		{"12345", "12347", "12347", "12348", "12340"},
		{"12340", "12340", "12400", "12401"},
	}
	actualChunks := make([][]string, 0, 2)
	err = StreamCSVInChunks(file, 5, func(s []string) {
		chunk := make([]string, len(s))
		copy(chunk, s)
		actualChunks = append(actualChunks, chunk)
	})
	if err != nil {
		t.Fatalf("expected nil error but got: %v", err)
	}

	// make sure each chunksize is the maximum chunksize possible
	expentedChunksLength := len(expectedChunks)
	actualChunksLength := len(actualChunks)
	if actualChunksLength != expentedChunksLength {
		t.Fatalf("number of chunks should be %v but got %v", expentedChunksLength, actualChunksLength)
	}

	for i, e := range expectedChunks {
		for j, a := range e {
			if a != actualChunks[i][j] {
				t.Fatalf("expectedChunks[%v][%v] (%q) is not equal to actualChunks[%v][%v] (%q)", i, j, a, i, j, actualChunks[i][j])
			}
		}
	}
}
