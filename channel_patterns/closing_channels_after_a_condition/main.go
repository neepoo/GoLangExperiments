package main

import (
	"fmt"
	"github.com/neepoo/GoLangExperiments/channel_patterns/broadcasting_to_multiple_goroutine"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"

	"github.com/neepoo/GoLangExperiments/channel_patterns/fanning_in_and_out"
)

const MAXDOWNLOADER = 20

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

func longestWords(quit <-chan struct{}, words <-chan string) <-chan string {
	longWords := make(chan string)
	go func() {
		defer close(longWords)
		uniqueWordsMap := make(map[string]bool)
		uniqueWords := make([]string, 0)
		moreData, word := true, ""
		for moreData {
			select {
			case word, moreData = <-words:
				if moreData && !uniqueWordsMap[word] {
					uniqueWordsMap[word] = true
					uniqueWords = append(uniqueWords, word)
				}

			case <-quit:
				return
			}
		}
		sort.Slice(uniqueWords, func(i, j int) bool {
			return len(uniqueWords[i]) > len(uniqueWords[j]) || (len(uniqueWords[i]) == len(uniqueWords[j]) && uniqueWords[i] < uniqueWords[j])
		})
		longWords <- strings.Join(uniqueWords[:10], ", ")
	}()
	return longWords
}

func frequentWords(quit <-chan struct{}, words <-chan string) <-chan string {
	mostFrequentWords := make(chan string)
	go func() {
		defer close(mostFrequentWords)
		freqMap := make(map[string]int)
		freqList := make([]string, 0)
		moreData, word := true, ""
		for moreData {
			select {
			case word, moreData = <-words:
				if moreData {
					if freqMap[word] == 0 {
						freqList = append(freqList, word)
					}
					freqMap[word] += 1
				}
			case <-quit:
				return
			}
		}
		sort.Slice(freqList, func(i, j int) bool {
			return freqMap[freqList[i]] > freqMap[freqList[j]] ||
				(freqMap[freqList[i]] == freqMap[freqList[j]] && freqList[i] < freqList[j])
		})
		mostFrequentWords <- strings.Join(freqList[:10], ", ")
	}()
	return mostFrequentWords
}

func main() {
	quitWords := make(chan struct{})
	quit := make(chan struct{})
	defer close(quit)
	urls := generateUrls(quit)
	pages := make([]<-chan string, MAXDOWNLOADER)
	for i := 0; i < MAXDOWNLOADER; i++ {
		pages[i] = downloadPages(quit, urls)

	}
	words := Take(quitWords, 100000, extractWords(quit, fanning_in_and_out.FanIn(quit, pages...)))

	wordsMulti := broadcasting_to_multiple_goroutine.BroadCast(quit, words, 2)
	longestResults := longestWords(quit, wordsMulti[0])
	frequentResults := frequentWords(quit, wordsMulti[1])

	fmt.Println("Longest Words: ", <-longestResults)
	fmt.Println("Most Frequent Words: ", <-frequentResults)

}
