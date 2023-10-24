package database

import (
	"fmt"
)

func (m *Manager) InsertPDFSummary(s PDFSummary) (int64, error) {
	result, err := m.db.Exec(
		"INSERT INTO pdf_summary (file, summary, title, intermediate_summary) VALUES ($1, $2, $3, $4)",
		s.File,
		s.Summary,
		s.Title,
		s.IntermediateSummary,
	)
	if err != nil {
		return -1, fmt.Errorf("could not insert pdf summary: %v", err)
	}

	return result.LastInsertId()
}

func (m *Manager) SearchPDFSummaries(q string) ([]PDFSummary, error) {
	var summaries []PDFSummary
	if err := m.db.Select(&summaries, "SELECT * FROM pdf_summary WHERE title LIKE $1 or summary LIKE $1", "%"+q+"%"); err != nil {
		return nil, fmt.Errorf("could not select pdf summaries: %v", err)
	}

	return summaries, nil
}

func (m *Manager) GetPDFSummaries() ([]PDFSummary, error) {
	var summaries []PDFSummary
	if err := m.db.Select(&summaries, "SELECT * FROM pdf_summary limit 5"); err != nil {
		return nil, fmt.Errorf("could not select pdf summaries: %v", err)
	}

	return summaries, nil
}

func (m *Manager) GetPDFSummary(id string) (PDFSummary, error) {
	var summary PDFSummary
	if err := m.db.Get(&summary, "SELECT * FROM pdf_summary WHERE id = $1", id); err != nil {
		return PDFSummary{}, fmt.Errorf("could not select pdf summary: %v", err)
	}

	return summary, nil
}
