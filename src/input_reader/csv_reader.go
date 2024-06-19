package inputreader

import (
	"encoding/csv"
	"os"

	"github.com/molinama/timescale/src/repository"
)

type CSVReader struct {
	reader *csv.Reader
	file   *os.File
}

func NewCSVReader(csvFilePath string) (*CSVReader, error) {
	file, err := os.Open(csvFilePath)
	if err != nil {
		return nil, err
	}

	return &CSVReader{
		file:   file,
		reader: csv.NewReader(file),
	}, nil

}
func (r *CSVReader) Read() (*repository.QueryParams, error) {
	data, err := r.reader.Read()
	if err != nil {
		return nil, err
	}

	return repository.NewQueryParams(data)
}

func (r *CSVReader) Close() error {
	return r.file.Close()
}
