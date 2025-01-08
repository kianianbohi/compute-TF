package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"sort"
	"strings"
)

var urls = [...]string{
	"https://go.dev/blog/chacha8rand",
	"https://go.dev/blog/randv2",
	"https://go.dev/blog/survey2024-h1-results",
}

func main() {

	result := make(map[string]map[string]float64)

	for _, url := range urls {
		text, err := fetchPage(url)
		if err != nil {
			fmt.Printf("Error fetching URL %s: %v\n", url, err)
			continue
		}

		tf := calculateTF(text)
		result[url] = tf
	}

	for url, tf := range result {
		fmt.Printf("Top 5 words for URL %s:\n", url)
		topWords := getTopWords(tf, 5)
		for _, word := range topWords {
			fmt.Printf("%s: %f\n", word, tf[word])
		}
		fmt.Println()
	}
}

func fetchPage(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch page: %s", url)
	}

	// Parse the HTML document
	z := html.NewTokenizer(resp.Body)
	var textContent strings.Builder

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			// End of the document
			return textContent.String(), nil
		case html.TextToken:
			textContent.WriteString(z.Token().Data)
		}
	}
}

func calculateTF(text string) map[string]float64 {
	words := strings.Fields(strings.ToLower(text))
	tf := make(map[string]float64)
	totalWords := len(words)

	for _, word := range words {
		word := strings.Trim(word, ".,;:!?\"'()[]{}")
		tf[word]++

	}

	for word := range tf {
		tf[word] = tf[word] / float64(totalWords)
	}

	return tf
}

func getTopWords(tf map[string]float64, topN int) []string {
	type wordCount struct {
		word string
		tf   float64
	}

	wordList := make([]wordCount, 0, len(tf))
	for word, tfs := range tf {
		wordList = append(wordList, wordCount{word, tfs})
	}

	sort.Slice(wordList, func(i, j int) bool {
		return wordList[i].tf > wordList[j].tf
	})

	topWords := make([]string, 0, topN)
	for i := 0; i < topN && i < len(wordList); i++ {
		topWords = append(topWords, wordList[i].word)
	}
	return topWords
}
