package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// Fhash calculate specify filepath file sha256
func Fhash(filepath string) []byte {
	file, _ := os.Open(filepath)
	defer file.Close()
	sha := sha256.New()
	io.Copy(sha, file)
	return sha.Sum(nil)
}

func DirHash(dir string, wg *sync.WaitGroup) {
	files, _ := os.ReadDir(dir)
	for _, file := range files {
		if file.IsDir() {
			DirHash(filepath.Join(dir, file.Name()), wg)
		} else {
			wg.Add(1)
			go func(filename string) {
				fpath := filepath.Join(dir, filename)
				hash := Fhash(fpath)
				fmt.Printf("%s = %x\n", fpath, hash)
				wg.Done()
			}(file.Name())
		}
	}
}

func main() {
	dir := os.Args[1]
	wg := new(sync.WaitGroup)
	DirHash(dir, wg)
	wg.Wait()
}
