package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/labstack/echo/v4"
	"github.com/rhnvrm/summio/pkg/database"
	"github.com/rhnvrm/summio/pkg/summarizer"
)

func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("missing env var %q", key)
	}
	return v
}

func getBoolEnv(key string) bool {
	v := os.Getenv(key)
	return v != ""
}

func genRandomString(size int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz"
	var b strings.Builder
	for i := 0; i < size; i++ {
		b.WriteByte(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

func main() {
	var (
		dbPath    = mustGetEnv("DB_PATH")
		filesPath = mustGetEnv("FILES_PATH")
		debugMode = getBoolEnv("DEBUG")
	)

	db, err := database.InitDB(dbPath)
	if err != nil {
		log.Fatalf("could not init db: %v", err)
	}

	e := echo.New()
	e.Debug = debugMode
	e.GET("/api/pdf", func(c echo.Context) error {
		summaries, err := db.GetPDFSummaries()
		if err != nil {
			return fmt.Errorf("could not get pdf summaries: %w", err)
		}

		return c.JSON(200, summaries)
	})
	e.POST("/api/pdf", func(c echo.Context) error {
		file, err := c.FormFile("file")
		if err != nil {
			return fmt.Errorf("could not get form file: %w", err)
		}

		src, err := file.Open()
		if err != nil {
			return fmt.Errorf("could not open file: %w", err)
		}
		defer src.Close()

		// check mime type == pdf
		mtype, err := mimetype.DetectReader(src)
		if err != nil {
			return fmt.Errorf("could not detect mime type: %w", err)
		}

		if mtype.String() != "application/pdf" {
			return fmt.Errorf("invalid mime type: %s", mtype.String())
		}

		src.Seek(0, io.SeekStart)

		fname := genRandomString(10) + ".pdf"
		fpath := filepath.Join(filesPath, fname)
		dst, err := os.Create(fpath)
		if err != nil {
			return fmt.Errorf("could not create file: %w", err)
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return fmt.Errorf("could not copy file: %w", err)
		}

		sum, err := summarizer.New()
		if err != nil {
			return fmt.Errorf("could not create summarizer: %w", err)
		}

		log.Println("loading pdf")
		docs, err := summarizer.LoadPDF(fpath)
		if err != nil {
			return fmt.Errorf("could not load pdf: %w", err)
		}

		log.Printf("pdf loaded with %d sub-docs", len(docs))

		log.Println("summarizing")
		out, err := sum.SummarizeDocs(docs)
		if err != nil {
			return fmt.Errorf("could not summarize docs: %w", err)
		}

		dbSummary := database.PDFSummary{
			File:                fname,
			Summary:             out.Summary,
			Title:               out.Title,
			IntermediateSummary: out.IntermediateSummary,
		}

		id, err := db.InsertPDFSummary(dbSummary)
		if err != nil {
			return fmt.Errorf("could not insert pdf summary: %w", err)
		}

		dbSummary.ID = id

		return c.JSON(200, dbSummary)
	})

	e.Logger.Fatal(e.Start(":1323"))
}
