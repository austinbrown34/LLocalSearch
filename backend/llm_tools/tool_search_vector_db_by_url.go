package llm_tools

import (
	"context"
	"encoding/json"
	"os"

	"github.com/nilsherzig/localLLMSearch/utils"
	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/tools"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/chroma"
)

// SearchVectorDBByUrl searches through added files and websites by URL.
type SearchVectorDBByUrl struct {
	CallbacksHandler callbacks.Handler
	SessionString    string
}

var _ tools.Tool = SearchVectorDBByUrl{}

// type SearchResult struct {
// 	Text   string
// 	Source string
// }

func (c SearchVectorDBByUrl) Description() string {
	return `Searches through added files and websites for a specific URL. Returns documents whose metadata URL matches the provided input URL.`
}

func (c SearchVectorDBByUrl) Name() string {
	return "SearchVectorDBByUrl"
}

func (c SearchVectorDBByUrl) Call(ctx context.Context, input string) (string, error) {
	if c.CallbacksHandler != nil {
		c.CallbacksHandler.HandleToolStart(ctx, input)
	}

	// Initialize the vector store with the same setup as the original tool
	llm, err := utils.NewOllamaEmbeddingLLM()
	if err != nil {
		return "", err
	}

	ollamaEmbedder, err := embeddings.NewEmbedder(llm)
	if err != nil {
		return "", err
	}

	store, errNs := chroma.New(
		chroma.WithChromaURL(os.Getenv("CHROMA_DB_URL")),
		chroma.WithEmbedder(ollamaEmbedder),
		chroma.WithNameSpace(c.SessionString),
	)
	if errNs != nil {
		return "", errNs
	}

	// Retrieve all documents as we need to manually filter by URL
	amountOfResults := 5 // Fetch all documents to filter manually
	options := []vectorstores.Option{}

	retriever := vectorstores.ToRetriever(store, amountOfResults, options...)
	docs, err := retriever.GetRelevantDocuments(context.Background(), input)
	if err != nil {
		return "", err
	}

	var results []Result

	// Filter documents by matching metadata URL with input URL
	for _, doc := range docs {
		sourceURL, ok := doc.Metadata["url"].(string)
		if !ok {
			continue
		}
		if sourceURL == input {
			results = append(results, Result{
				Text:   doc.PageContent,
				Source: sourceURL,
			})
		}
	}

	if len(results) == 0 {
		results = append(results, Result{Text: "No matching documents found for the provided URL."})
	}

	resultJson, err := json.Marshal(results)
	if err != nil {
		return "", err
	}

	if c.CallbacksHandler != nil {
		c.CallbacksHandler.HandleToolEnd(ctx, input)
	}

	return string(resultJson), nil
}
