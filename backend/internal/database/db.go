package database

import (
	"encoding/csv"
	"os"

	"github.com/furqanmk/reviews-browser/config"
)

// DB holds paths to CSV files.
type DB struct {
	config *config.Config
}

func (db *DB) Close() {
	panic("unimplemented")
}

// NewDBConnection initializes DB with CSV file paths.
func NewDBConnection(config *config.Config) (*DB, error) {
	return &DB{config: config}, nil
}

func getReader(filePath string) (*csv.Reader, *os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	reader := csv.NewReader(file)
	return reader, file, nil
}

func getWriter(filePath string) (*csv.Writer, *os.File, error) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, nil, err
	}
	writer := csv.NewWriter(file)
	return writer, file, nil
}

func emptyFile(filePath string, header []string) (*csv.Writer, *os.File, error) {
	outFile, err := os.Create(filePath)
	if err != nil {
		return nil, nil, err
	}
	writer := csv.NewWriter(outFile)

	// Write header row
	if err := writer.Write(header); err != nil {
		return nil, nil, err
	}

	return writer, outFile, nil
}