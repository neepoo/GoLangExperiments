package main

import (
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type CodeDepth struct {
	file  string
	level int
}

func deepestNestedBlock(filename string) CodeDepth {
	code, _ := os.ReadFile(filename)
	max, level := 0, 0
	for _, c := range code {
		if c == '{' {
			level++
			max = int(math.Max(float64(max), float64(level)))
		} else if c == '}' {
			level--
		}
	}
	return CodeDepth{filename, max}
}

func forkIfNeeded(path string, info os.FileInfo, wg *sync.WaitGroup, results chan CodeDepth) {
	if !info.IsDir() && strings.HasSuffix(path, ".go") {
		wg.Add(1)
		go func() {
			results <- deepestNestedBlock(path)
			wg.Done()
		}()
	}
}

func joinResults(partialResults chan CodeDepth) chan CodeDepth {
	finalResult := make(chan CodeDepth)
	maxDepth := CodeDepth{}
	go func() {
		for pr := range partialResults {
			if pr.level > maxDepth.level {
				maxDepth = pr
			}
		}
		finalResult <- maxDepth
	}()
	return finalResult
}

func main() {
	dir := os.Args[1]
	partialResults := make(chan CodeDepth)
	wg := new(sync.WaitGroup)

	filepath.Walk(dir,
		func(path string, info fs.FileInfo, err error) error {
			forkIfNeeded(path, info, wg, partialResults)
			return err
		})
	finalResult := joinResults(partialResults)

	wg.Wait()
	close(partialResults)
	res := <-finalResult
	fmt.Printf("%s has the deepest nested code block of %d\n", res.file, res.level)
}
