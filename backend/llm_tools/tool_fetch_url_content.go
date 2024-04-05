package llm_tools

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/nilsherzig/localLLMSearch/utils"
	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/tools"
)

// FetchURLContent is a tool that fetches and processes content from a specific URL.
type FetchURLContent struct {
	CallbacksHandler callbacks.Handler
	SessionString    string
}

var _ tools.Tool = FetchURLContent{}

func (f FetchURLContent) Description() string {
	return `Fetches and processes content from a specific URL. Useful for retrieving specific information from the internet.`
}

func (f FetchURLContent) Name() string {
	return "FetchURLContent"
}

func (f FetchURLContent) Call(ctx context.Context, inputURL string) (string, error) {
	if f.CallbacksHandler != nil {
		f.CallbacksHandler.HandleToolStart(ctx, inputURL)
	}

	// Check if inputURL is a valid URL, this step is crucial to ensure safety and correctness.
	if !utils.IsValidURL(inputURL) {
		return "", fmt.Errorf("invalid URL")
	}

	// Fetching content from the URL
	resp, err := http.Get(inputURL)
	if err != nil {
		log.Println("Error making the request:", err)
		return "", err
	}
	defer resp.Body.Close()

	// Reading the response body
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading the response body:", err)
		return "", err
	}
	bodyString := string(bodyBytes)

	// Optionally, process the body content as needed here.
	// For simplicity, this example will just log the size of the content.
	log.Printf("Fetched content size: %d bytes\n", len(bodyString))

	// Here you might want to process the fetched content, similar to how the original code processed search results.
	// Since the specifics depend on your requirements, this part is left as an exercise.

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		// This example does not include the implementation of utils.DownloadWebsiteToVectorDB.
		// You would need to adapt or implement this function to handle content processing and storage as per your needs.
		err := utils.DownloadWebsiteToVectorDB(ctx, inputURL, f.SessionString)
		if err != nil {
			log.Printf("error processing the content: %s", err.Error())
			return
		}
	}()

	wg.Wait()

	result := fmt.Sprintf("Downloaded website to vector db. You dont know anything new from this tool, you have to search through the vector db to find anything about the downloaded website.")

	if f.CallbacksHandler != nil {
		f.CallbacksHandler.HandleToolEnd(ctx, result)
	}

	return result, nil
}
