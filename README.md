# Efficient Data Analysis with Hashmaps, Bitset, and Count-Min Sketch

This project is a Go-based project designed to analyze datasets stored in files,
specifically focusing on counting keys and overlaps between them.

In this project, I demonstrate achieving both correctness and efficiency through the use of hashmaps and advanced
data structures like Bitset and Count-Min Sketch.

## Objectives

This solution aims to compute the following metrics for datasets in the `data` and `data/bigger` directories:

1. **Count of Keys in Each File**: Total number of keys, including duplicates, in each dataset.
2. **Count of Distinct Keys in Each File**: Number of unique keys in each dataset.
3. **Count of Maximum Possible Overlap of All Keys Between the Files**: Total overlap considering duplicate keys.
4. **Count of Overlap of Distinct Keys Between the Files**: Overlap considering only unique keys.

## Approach

### Correctness: Hashmaps

Hashmaps are used to generate accurate results for all the objectives.
While this approach is not the most efficient for large datasets,
it serves as a benchmark to verify the accuracy of more optimized methods.

### Efficiency: Bitset and Count-Min Sketch

To improve efficiency and scalability:
- **[Bitset](https://github.com/bits-and-blooms/bitset)**: Used for calculating distinct overlaps and counting unique keys.
- **[Count-Min Sketch](https://github.com/shenwei356/countminsketch)**: A probabilistic data structure used to approximate counts and overlaps
efficiently. While it provides accurate estimates for most objectives, it cannot fully
compute the total overlap of all keys between files.

## Sample Data

- **Small datasets**: Located in the `data` folder for testing and validation.
- **Large datasets**: Found in the `data/bigger` folder to stress-test the algorithms. (this is randomly generated values to test the algorithms)

## Installation

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/spankie/infosum.git
   cd infosum
   ```

2. **Install Dependencies**:
   Ensure you have [Go](https://golang.org/dl/) installed. Run:
   ```bash
   go mod tidy
   ```

## Usage

1. **Prepare Datasets**:
   there are dataset files in the `data` or `data/bigger` directories as needed.

2. **Run the Program**:
   Execute the code with:
   
   ```bash
   make run
   ```
   
   if you don't have make installed, then you can run the program like this (ensure to replace the
       placeholder filenames with your filenames):
   
   ```bash
   go run main.go --chunksize=<chunksize> --fileA=<fileA> --fileB=<fileB>
   ```
   
   e.g:
   
   ```bash
   go run main.go --chunksize=1000 --fileA=data/A_f.csv --fileB=data/B_f.csv
   ```

   You can also run the test using make:
   ```
   make test
   ```

   of if you don't have make installed:

   ```
   go test ./...
   ```

4. **View Results**:
   The program outputs the results for each file and their comparisons, like this:
   
   ```
   File A: data/A_f.csv
   File B: data/B_f.csv
    
   Count of keys in file A: 86535
   Count of distinct keys file A: 72798
    
   Count of keys in file B: 72846
   Count of distinct keys in file B: 72814
    
   Count of distinct overlap: 58221
   Count of total max overlap: 66549
   ```

> [!NOTE]
> By default the code use the `bitset` and `Count Min Sketch` to get the values. But if you want to
> run it using hashmap, you can uncomment the code in `main.go` file to use the hashmap instead

## Design Highlights

- **Correctness Verification**: Hashmaps ensure the correctness of all metrics for testing purposes.
- **Efficiency for Large Datasets**: Bitset and Count-Min Sketch drastically reduce memory
and computational overhead especially for really large dataset, providing near-accurate results efficiently.
- **Test Coverage**: Includes both small and large datasets to validate and stress-test the implementation.

## folder structure

- **`data` Folder**: Contains small sample datasets for validation.
- **`data/bigger` Folder**: Contains larger datasets for performance testing.
- **`bitset`**: implements the solution using [Bitset](https://github.com/bits-and-blooms/bitset)
and [Count-Min Sketch](https://github.com/shenwei356/countminsketch) for calculating the results
- **`hashmaps`**: implements the solution using hashmap for calculating the results.
- **`dataset`**: contains function to read the dataset from the file
- **`results`**: is a package that defines a datastructure that defines the expected result and prints it out.

## Limitations

- The **Count-Min Sketch** is highly efficient but cannot accurately compute the total overlap
of all keys between files due to its probabilistic nature.

By combining the reliability of hashmaps and the efficiency of advanced data structures,
this project demonstrates the trade-offs between correctness and performance in large-scale data analysis.
```

