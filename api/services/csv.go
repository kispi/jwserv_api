package services

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io/ioutil"
	"os"

	"../core"
)

type CSVService struct {
	FileName string
	Writer   *csv.Writer
}

func (c *CSVService) NewCSV(fileName string) (*CSVService, error) {
	file, err := os.Create(fileName)
	if err != nil {
		core.Log.Warning(err)
		return c, errors.New("FAILED_TO_CREATE_CSV")
	}

	c.FileName = fileName
	wr := csv.NewWriter(bufio.NewWriter(file))
	c.Writer = wr
	wr.Flush()

	return c, nil
}

func (c *CSVService) AddRow(row []string) error {
	err := c.Writer.Write(row)
	if err != nil {
		return errors.New("FAILED_TO_ADD_ROW")
	}
	return nil
}

func (c *CSVService) AddRows(rows [][]string) error {
	err := c.Writer.WriteAll(rows)
	if err != nil {
		return errors.New("FAILED_TO_ADD_ROWS")
	}
	return nil
}

func (c *CSVService) SaveFileAsBytes() ([]byte, error) {
	fileAsByte, err := ioutil.ReadFile(c.FileName)
	if err != nil {
		return nil, errors.New("FAILED_TO_SAVE_CREATE_BYTE_STREAM")
	}
	c.Writer.Flush()

	err = os.Remove(c.FileName)
	if err != nil {
		core.Log.Warning("FAILED_TO_REMOVE_TEMPORARY_FILE", c.FileName)
	}
	return fileAsByte, nil
}
