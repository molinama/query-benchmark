# Query Benchmark Tool README

## High-Level Goal

The Query Benchmark tool is a command-line application designed to benchmark `SELECT` query performance across multiple workers/clients against a TimescaleDB instance. The tool accepts input either as a CSV-formatted file or standard input, specifying query parameters and the number of concurrent workers. The tool processes queries concurrently and outputs a summary of the following statistics after processing all queries:

- Number of queries processed
- Total processing time across all queries
- Minimum query time for a single query
- Median query time
- Average query time
- Maximum query time

## Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/molinama/query-benchmark.git
    
    cd query-benchmark
    ```

## Usage

To use the Query Benchmark tool, you should to specify the path to the CSV file containing query parameters and optionally the number of concurrent workers.

### Command-line Arguments

- `-csv` : The file path to the CSV file containing query parameters (default: query_params.csv).
- `-workers` : The number of workers for the pool (default: 10).

### Example Command

```sh
go run main.go -csv=query_params.csv -workers=20
```

### CSV Format

The CSV file should contain the necessary query parameters, including hostname and raw query strings.

```csv
hostname,start_time,end_time
host_000008,2017-01-01 08:59:22,2017-01-01 09:59:22
```

### Output

After processing the queries, the tool will output the following statistics:

```bash
Total Queries: 200
Number of queries successfully processed: 200
Total processing time: 1.764687507s
Minimum query time: 2.537792ms
Median query time: 5.300917ms
Average query time: 8.823437ms
Maximum query time: 94.558459ms
Total Errors: 0
```

### Usage Instructions

1. Start Timescaledb

    ```sh
    make start-timescaledb
    ```

2. Run the application

    ```sh
    make run
    ```

    > If you want to use your own csv file, use the "csv" variable to run.
    >
    > ```sh
    > make run csv=filepath

3. Stop Timescaledb

    ```sh
    make stop-timescaledb
    ```
