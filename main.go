package main

import (
	"log"
	"os"

	"github.com/rhnvrm/summio/pkg/database"
	"github.com/rhnvrm/summio/pkg/summarizer"
)

func main() {
	log.Println("starting summio")
	sum, err := summarizer.New()
	if err != nil {
		log.Fatalf("could not create summarizer: %v", err)
	}

	log.Println("loading pdf")
	docs, err := summarizer.LoadPDF("testdata/1694170811233.pdf")
	if err != nil {
		log.Fatalf("could not load pdf: %v", err)
	}

	log.Printf("pdf loaded with %d sub-docs", len(docs))

	log.Println("summarizing")
	out, err := sum.SummarizeDocs(docs)
	if err != nil {
		log.Fatalf("could not summarize docs: %v", err)
	}

	dbPath := os.Getenv("DB_PATH")
	db, err := database.InitDB(dbPath)
	if err != nil {
		log.Fatalf("could not init db: %v", err)
	}

	if err := db.InsertPDFSummary(database.PDFSummary{
		File:                "testdata/1694170811233.pdf",
		Summary:             out.Summary,
		Title:               out.Title,
		IntermediateSummary: out.IntermediateSummary,
	}); err != nil {
		log.Fatalf("could not insert pdf summary: %v", err)
	}

	log.Printf("out: %#v", out)
}
