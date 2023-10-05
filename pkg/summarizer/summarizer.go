package summarizer

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
)

const (
	templateConciseMd = `Write a concise summary of the following in markdown format:

"{{.context}}"

MARKDOWN RESULT:`

	templateConcise = `Write a concise one-line summary of the following: 

"{{.context}}"

CONCISE ONELINE SUMMARY:`

	templateTitle = `Write a very concise title for the following summary:

"{{.context}}"

TITLE:`
)

func LoadPDF(path string) ([]schema.Document, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}

	finfo, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("could not stat file: %w", err)
	}

	p := documentloaders.NewPDF(f, finfo.Size())

	docs, err := p.LoadAndSplit(
		context.Background(),
		textsplitter.NewTokenSplitter(),
	)
	if err != nil {
		return nil, fmt.Errorf("could not load pdf: %w", err)
	}

	return docs, nil
}

func LoadLLM() (llms.LanguageModel, error) {
	llm, err := openai.New(openai.WithModel("gpt-3.5-turbo"))
	if err != nil {
		return nil, fmt.Errorf("could not open openai: %w", err)
	}

	return llm, nil
}

type Summarizer struct {
	llm llms.LanguageModel
}

func New() (*Summarizer, error) {
	llm, err := LoadLLM()
	if err != nil {
		return nil, fmt.Errorf("could not load llm: %w", err)
	}

	return &Summarizer{
		llm: llm,
	}, nil
}

type Summary struct {
	Summary             string
	IntermediateSummary []string
	Title               string
}

func (s *Summarizer) MakeSummaryChain() (chains.MapReduceDocuments, error) {
	mapChain := chains.NewLLMChain(s.llm,
		prompts.NewPromptTemplate(templateConcise, []string{"context"}),
	)

	combineChain := chains.NewStuffDocuments(chains.NewLLMChain(s.llm,
		prompts.NewPromptTemplate(templateConciseMd, []string{"context"}),
	))

	out := chains.NewMapReduceDocuments(mapChain, combineChain)
	out.MaxNumberOfConcurrent = 10
	out.ReturnIntermediateSteps = true

	return out, nil
}

func (s *Summarizer) SummarizeDocs(docs []schema.Document) (Summary, error) {
	ctx := context.Background()

	chain, err := s.MakeSummaryChain()
	if err != nil {
		return Summary{}, fmt.Errorf("could not make chain: %w", err)
	}

	outVals, err := chains.Call(ctx, chain, map[string]any{
		"input_documents": docs,
	})
	if err != nil {
		return Summary{}, fmt.Errorf("could not call chain: %w", err)
	}

	intermediateSteps := []string{}
	steps, ok := outVals["intermediateSteps"].([]map[string]interface{})
	if ok {
		for _, step := range steps {
			intermediateSteps = append(intermediateSteps, step["text"].(string))
		}
	} else {
		log.Println("could not convert intermediate steps")
		log.Println("type of intermediate steps: ", fmt.Sprintf("%T", outVals["intermediateSteps"]))
	}

	summary, ok := outVals["text"].(string)
	if !ok {
		return Summary{}, fmt.Errorf("could not convert summary output to string: %w", err)
	}

	titleChain := chains.NewLLMChain(s.llm,
		prompts.NewPromptTemplate(templateTitle, []string{"context"}),
	)

	summaryBuf := strings.NewReader(summary)
	sl := documentloaders.NewText(summaryBuf)
	summaryDoc, err := sl.LoadAndSplit(ctx, textsplitter.NewTokenSplitter())
	if err != nil {
		return Summary{}, fmt.Errorf("could not load summary doc: %w", err)
	}

	titleVals, err := chains.Call(ctx, titleChain, map[string]any{
		"context": summaryDoc,
	})

	title, ok := titleVals["text"].(string)
	if !ok {
		return Summary{}, fmt.Errorf("could not convert title to string: %w", err)
	}

	return Summary{
		Summary:             summary,
		IntermediateSummary: intermediateSteps,
		Title:               title,
	}, nil
}
