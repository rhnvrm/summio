package database

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type PDFSummary struct {
	ID                  int64    `db:"id" json:"id"`
	File                string   `db:"file" json:"file"`
	Summary             string   `db:"summary" json:"summary"`
	Title               string   `db:"title" json:"title"`
	IntermediateSummary JSONList `db:"intermediate_summary" json:"intermediate_summary"`
}

type JSONList []string

var (
	_ driver.Valuer = JSONList{}
	_ sql.Scanner   = (*JSONList)(nil)
)

func (jl *JSONList) Scan(src interface{}) error {
	var s string
	switch src := src.(type) {
	case []byte:
		s = string(src)
	case string:
		s = src
	default:
		return fmt.Errorf("could not scan JSONList: unknown type")
	}

	return json.Unmarshal([]byte(s), jl)
}

func (jl JSONList) Value() (driver.Value, error) {
	return json.Marshal(jl)
}
