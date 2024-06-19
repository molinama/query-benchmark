## WorkerPool Tool README

### High-Level Goal
The WorkerPool tool is a command-line application designed to benchmark `SELECT` query performance across multiple workers/clients against a TimescaleDB instance. The tool accepts input either as a CSV-formatted file or standard input, specifying query parameters and the number of concurrent workers. The tool processes queries concurrently and outputs a summary of the following statistics after processing all queries:
- Number of queries processed
- Total processing time across all queries
- Minimum query time for a single query
- Median query time
- Average query time
- Maximum query time

### Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/yourusername/workerpool.git
    cd workerpool
    ```

2. Build the application:
    ```sh
    go build -o workerpool main.go
    ```

### Usage

To use the WorkerPool tool, you need to specify the path to the CSV file containing query parameters and optionally the number of concurrent workers and the database connection string.

#### Command-line Arguments:
- `-csv` : The file path to the CSV file containing query parameters.
- `-workers` : The number of workers for the pool (default: 10).
- `-db` : The database connection string.

#### Example Command:
```sh
./workerpool -csv=query_params.csv -workers=20 -db="postgres://user:password@localhost:5432/mydb"
```

### CSV Format
The CSV file should contain the necessary query parameters, including hostname and raw query strings.

### Output
After processing the queries, the tool will output the following statistics:
```
Number of queries processed: 100
Total processing time: 1m23.456s
Minimum query time: 123.456ms
Median query time: 234.567ms
Average query time: 345.678ms
Maximum query time: 456.789ms
```

### Example Usage

1. Prepare a CSV file `query_params.csv`:
    ```csv
    hostname,raw_query
    db1.example.com,SELECT * FROM table1
    db2.example.com,SELECT * FROM table2
    ...
    ```

2. Run the application:
    ```sh
    ./workerpool -csv=query_params.csv -workers=10 -db="postgres://user:password@localhost:5432/mydb"
    ```