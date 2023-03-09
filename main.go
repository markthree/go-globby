package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

type Match = func(name string) bool

var match Match

func initMatch(payload string) {
	payload = payload[:len(payload)-1]
	patterns := strings.Split(payload, ",")
	patternsLen := len(patterns)
	match = func(name string) bool {
		for i := 0; i < patternsLen; i++ {
			matched, err := filepath.Match(strings.Replace(patterns[i], "**/", "", 1), name)
			if err != nil {
				fail(err)
				return false
			}
			if !matched {
				return false
			}
		}
		return true
	}
}

func deepScan(folder string) {
	entrys, err := os.ReadDir(folder)
	if err != nil {
		fail(err)
	}

	var wg sync.WaitGroup
	entrysLen := len(entrys)

	if entrysLen == 0 {
		return
	}

	wg.Add(entrysLen)

	for i := 0; i < entrysLen; i++ {
		go func(entry fs.DirEntry) {
			defer wg.Done()
			p := path.Join(folder, entry.Name())
			if entry.IsDir() {
				deepScan(p)
				return
			}

			if match(p) {
				success(p)
			}
		}(entrys[i])
	}

	wg.Wait()
	return
}

func scan(folder string) {
	entrys, err := os.ReadDir(folder)
	if err != nil {
		fail(err)
	}

	var wg sync.WaitGroup
	entrysLen := len(entrys)

	if entrysLen == 0 {
		return
	}

	wg.Add(entrysLen)

	for i := 0; i < entrysLen; i++ {
		go func(entry fs.DirEntry) {
			defer wg.Done()
			if !entry.IsDir() {
				p := path.Join(folder, entry.Name())
				if match(p) {
					success(p)
				}
			}
		}(entrys[i])
	}

	wg.Wait()
	return
}

func success(path string) {
	fmt.Fprintf(os.Stdout, "%v\n", path)
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "%v", err)
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	payload, err := reader.ReadString('\n')

	if err != nil {
		fail(err)
		return
	}

	initMatch(payload)

	if strings.Contains(payload, "**/*") {
		deepScan(".")
	} else {
		scan(".")
	}

	success("EOF")
}
