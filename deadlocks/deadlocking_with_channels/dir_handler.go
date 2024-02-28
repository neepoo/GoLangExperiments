package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

//Let’s look at an example of a deadlock involving two channels. Consider a simple
//program that needs to recursively output file details, such as the filename, file size,
//and last modified date of all files under a directory. One solution is to have one goroutine
//that handles files and another that deals with directories.

func handleDirectoriesDeadlock(dirs <-chan string, files chan<- string) {
	for fullpath := range dirs {
		fmt.Println("reading all files from", fullpath)
		filesInDir, _ := os.ReadDir(fullpath)
		fmt.Printf("Pushing %d files from %s\n", len(filesInDir), fullpath)
		for _, file := range filesInDir {
			files <- filepath.Join(fullpath, file.Name())
		}
	}
}

/*
如何修复
*/

func handleDirectories(dirs <-chan string, files chan<- string) {
	// Creates a slice to store files that need to be pushed to the file handler's channel
	toPush := make([]string, 0)
	appendAllFiles := func(path string) {
		fmt.Println("Reading all files from", path)
		filesInDir, _ := os.ReadDir(path)
		fmt.Printf("Pushing %d files from %s\n", len(filesInDir), path)
		for _, f := range filesInDir {
			toPush = append(toPush, filepath.Join(path, f.Name()))
		}
	}
	for {
		if len(toPush) == 0 {
			// If there are no files to push, reads directory from
			// the input channel and adds all files in the directory
			appendAllFiles(<-dirs)
		} else {
			select {
			case fullpath := <-dirs:
				appendAllFiles(fullpath)
			case files <- toPush[0]:
				toPush = toPush[1:]
			}
		}
	}
}

//The reverse happens in the file handler goroutine. When the file handler meets a new
//directory, it sends it to the directory handler’s channel. The file handler consumes
//items from an input channel if the item is a file, and it outputs information about it,
//such as the file size and last modified date. If the item is a directory, it forwards the
//directory to the directory handler.

func handleFiles(files chan string, dirs chan string) {
	for path := range files {
		file, _ := os.Open(path)
		fileInfo, _ := file.Stat()
		if fileInfo.IsDir() {
			fmt.Printf("Pushing %s directory\n", fileInfo.Name())
			dirs <- path
		} else {
			fmt.Printf("File %s, size: %dMB, last modified: %s\n",
				fileInfo.Name(), fileInfo.Size()/(1024*1024),
				fileInfo.ModTime().Format("15:04:05"))
		}
	}
}

func main() {
	filesChannel := make(chan string)
	dirsChannel := make(chan string)
	go handleDirectoriesDeadlock(dirsChannel, filesChannel)
	go handleFiles(filesChannel, dirsChannel)
	dirsChannel <- os.Args[1]
	time.Sleep(2 * time.Second)
}
