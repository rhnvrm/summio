package main

import (
	"embed"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rhnvrm/summio/pkg/database"
	"github.com/rhnvrm/summio/pkg/database/migration"
	"github.com/rhnvrm/summio/pkg/summarizer"
)

type Envelope struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

var (
	//go:embed frontend/dist/*
	distFS embed.FS
)

func checkDir(dir string) {
	// Check if dir exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Directory does not exist, create it
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Printf("Failed to create directory: %v\n", err)
			return
		}
	}
}

func main() {
	var (
		dbPath    = getEnvDefault("DB_PATH", "summiodb.sqlite3")
		filesPath = getEnvDefault("FILES_PATH", "data")
		debugMode = getBoolEnv("DEBUG")
		address   = getEnvDefault("ADDRESS", ":1323")
		_         = mustGetEnv("OPENAI_API_KEY")
	)

	migration.RunMigration(dbPath)

	db, err := database.InitDB(dbPath)
	if err != nil {
		log.Fatalf("could not init db: %v", err)
	}

	checkDir(filesPath)

	e := echo.New()
	e.Debug = debugMode

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:       "frontend/dist",
		Index:      "index.html",
		HTML5:      true,
		Browse:     false,
		Filesystem: http.FS(distFS),
	}))

	e.Static("/api/static/docs/", filesPath)

	e.GET("/api/pdf", func(c echo.Context) error {
		summaries, err := db.GetPDFSummaries()
		if err != nil {
			debug(fmt.Errorf("could not get pdf summaries: %w", err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not get PDF summaries")
		}

		return c.JSON(200, newSuccessEnvelope(summaries))
	})

	e.GET("/api/pdf/:id", func(c echo.Context) error {
		id := c.Param("id")

		summary, err := db.GetPDFSummary(id)
		if err != nil {
			debug(fmt.Errorf("could not get pdf summary: %w", err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not get PDF summary")
		}

		return c.JSON(200, newSuccessEnvelope(summary))
	})

	e.POST("/api/pdf", func(c echo.Context) error {
		file, err := c.FormFile("file")
		if err != nil {
			debug(fmt.Errorf("could not get form file: %w", err))
			return echo.NewHTTPError(http.StatusBadRequest, "Could not get form file")
		}

		src, err := file.Open()
		if err != nil {
			debug(fmt.Errorf("could not open file: %w", err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not open file")
		}
		defer src.Close()

		// check mime type == pdf
		mtype, err := mimetype.DetectReader(src)
		if err != nil {
			debug(fmt.Errorf("could not detect mime type: %w", err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not detect mime type")
		}

		if mtype.String() != "application/pdf" {
			debug(fmt.Errorf("invalid mime type: %s", mtype.String()))
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid mime type")
		}

		// reset src to start
		src.Seek(0, io.SeekStart)

		fname := genRandomString(10) + ".pdf"
		fpath := filepath.Join(filesPath, fname)
		dst, err := os.Create(fpath)
		if err != nil {
			debug(fmt.Errorf("could not create file: %w", err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not create file")
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			debug(fmt.Errorf("could not copy file: %w", err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not copy file")
		}

		sum, err := summarizer.New()
		if err != nil {
			debug(fmt.Errorf("could not create summarizer: %w", err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not create summarizer")
		}

		debug("loading pdf")
		docs, err := summarizer.LoadPDF(fpath)
		if err != nil {
			debug(fmt.Errorf("could not load pdf: %w", err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not load pdf")
		}

		log.Printf("pdf loaded with %d sub-docs", len(docs))

		debug("summarizing")
		out, err := sum.SummarizeDocs(docs)
		if err != nil {
			debug(fmt.Errorf("could not summarize docs: %w", err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not summarize docs")
		}

		dbSummary := database.PDFSummary{
			File:                fname,
			Summary:             out.Summary,
			Title:               out.Title,
			IntermediateSummary: out.IntermediateSummary,
		}

		id, err := db.InsertPDFSummary(dbSummary)
		if err != nil {
			debug(fmt.Errorf("could not insert pdf summary: %w", err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not insert pdf summary")
		}

		dbSummary.ID = id

		return c.JSON(200, newSuccessEnvelope(dbSummary))
	})

	e.Logger.Fatal(e.Start(address))
}

func newSuccessEnvelope(data interface{}) Envelope {
	return Envelope{
		Status: "success",
		Data:   data,
	}
}

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

func getEnvDefault(key, defaultVal string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	return v
}

func genRandomString(size int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz"
	var b strings.Builder
	for i := 0; i < size; i++ {
		b.WriteByte(chars[rand.Intn(len(chars))])
	}

	// add timestamp to make it unique.
	t := time.Now().Unix()
	b.WriteString(fmt.Sprintf("-%d", t))

	return b.String()
}

func debug(err any) {
	if os.Getenv("DEBUG") == "true" {
		log.Printf("debug: %+v", err)
	}
}
