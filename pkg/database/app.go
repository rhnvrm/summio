package database

import (
	"encoding/json"
	"fmt"
)

func (m *Manager) InsertPDFSummary(s PDFSummary) error {
	intSummaryJSON, err := json.Marshal(s.IntermediateSummary)
	if err != nil {
		return fmt.Errorf("could not marshal intermediate summary: %v", err)
	}

	_, err = m.db.Exec(
		"INSERT INTO pdf_summary (file, summary, title, intermediate_summary) VALUES ($1, $2, $3, $4)",
		s.File,
		s.Summary,
		s.Title,
		string(intSummaryJSON),
	)
	if err != nil {
		return fmt.Errorf("could not insert pdf summary: %v", err)
	}

	return nil
}
