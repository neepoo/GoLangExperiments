package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

func generateUrls(quit <-chan struct{}) <-chan string {
	urls := make(chan string)
	go func() {
		defer close(urls)
		for i := 100; i < 130; i++ {
			url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
			select {
			case urls <- url:
			case <-quit:
				return
			}
		}
	}()
	return urls
}

func downloadPages(quit <-chan struct{}, urls <-chan string) <-chan string {
	pages := make(chan string)
	go func() {
		defer close(pages)
		moreData, url := true, ""
		for moreData {
			select {
			case url, moreData = <-urls:
				if moreData {
					resp, _ := http.Get(url)
					if resp.StatusCode != http.StatusOK {
						panic("Server's error: " + resp.Status)
					}
					body, _ := io.ReadAll(resp.Body)
					pages <- string(body)
					resp.Body.Close()
				}

			case <-quit:
				return
			}
		}
	}()
	return pages
}

func extractWords(quit <-chan struct{}, pages <-chan string) <-chan string {
	words := make(chan string)
	go func() {
		defer close(words)
		wordRegex := regexp.MustCompile(`[a-zA-Z]+`)
		moreData, pg := true, ""
		for moreData {
			select {
			case pg, moreData = <-pages:
				if moreData {
					for _, word := range wordRegex.FindAllString(pg, -1) {
						words <- strings.ToLower(word)
					}
				}
			case <-quit:
				return
			}
		}
	}()
	return words
}

type wordCount struct {
	word  string
	count int
}

type wordCounts []wordCount

func (w wordCounts) Len() int {
	return len(w)
}

func (w wordCounts) Less(i, j int) bool {
	return w[i].count > w[j].count || (w[i].count == w[j].count && w[i].word < w[j].word)
}

func (w wordCounts) Swap(i, j int) {
	w[i], w[j] = w[j], w[i]
}

func main() {
	quit := make(chan struct{})
	defer close(quit)
	results := extractWords(quit, downloadPages(quit, generateUrls(quit)))
	wordCountMap := make(map[string]int)
	for result := range results {
		wordCountMap[result] += 1
	}
	// 统计出现次数最多的单词
	var counts = make(wordCounts, 0, len(wordCountMap))
	for word, count := range wordCountMap {
		counts = append(counts, wordCount{word, count})
	}
	sort.Sort(counts)
	for i := 0; i < 10; i++ {
		fmt.Println(counts[i].word, counts[i].count)
	}

}
