package csvRead

import (
	"encoding/csv"
	"mime/multipart"
)

func ReadFile(file multipart.File) ([][]string, error) {
	reader := csv.NewReader(file)
	reader.Comma = ';'
	record, err := reader.ReadAll()

	if err != nil {
		return nil, err
	}
	return record, nil
}
