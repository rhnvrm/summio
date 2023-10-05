package database

type PDFSummary struct {
	ID                  int64    `db:"id"`
	File                string   `db:"file"`
	Summary             string   `db:"summary"`
	Title               string   `db:"title"`
	IntermediateSummary []string `db:"intermediate_summary"`
}
