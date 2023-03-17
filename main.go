package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

type Match = func(name string) bool

var match Match
var ignore Match

func initMatchers(payload string) {
	payload = strings.TrimSuffix(payload, "\n")
	patterns := strings.Split(payload, ",")

	// Preprocess patterns to separate match and ignore patterns
	matchPatterns := []string{}
	ignorePatterns := []string{}
	for _, pattern := range patterns {
		if strings.HasPrefix(pattern, "!") {
			ignorePatterns = append(ignorePatterns, strings.TrimPrefix(pattern, "!"))
		} else {
			matchPatterns = append(matchPatterns, pattern)
		}
	}

	// Preprocess patterns to remove "**/" prefix
	processedMatchPatterns := make([]string, len(matchPatterns))
	processedIgnorePatterns := make([]string, len(ignorePatterns))
	for i, pattern := range matchPatterns {
		processedMatchPatterns[i] = strings.Replace(pattern, "**/", "", 1)
	}
	for i, pattern := range ignorePatterns {
		processedIgnorePatterns[i] = strings.Replace(pattern, "**/", "", 1)
	}

	
	match = func(name string) bool {
		for _, pattern := range processedMatchPatterns {
			matched, err := filepath.Match(pattern, name)
			if err != nil {
				fail(err)
				return false
			}
			if matched {
				return true
			}
		}
		return false
	}

	ignore = func(name string) bool {
		for _, pattern := range processedIgnorePatterns {
			matched, err := filepath.Match(pattern, name)
			if err != nil {
				fail(err)
				return false
			}
			if matched {
				return true
			}
		}
		return false
	}
}
func scanFolder(folder string, deep bool, wg *sync.WaitGroup) {
	defer wg.Done()

	entrys, err := os.ReadDir(folder)
	if err != nil {
		fail(err)
		return
	}

	for _, entry := range entrys {
		p := path.Join(folder, entry.Name())
		if entry.IsDir() {
			if deep {
				wg.Add(1)
				go scanFolder(p, true, wg)
			}
		} else {
			if !ignore(p) && match(p) {
				success(p)
			}
		}
	}
}

func success(path string) {
	fmt.Println(path)
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	payload, err := reader.ReadString('\n')

	if err != nil {
		fail(err)
		return
	}

	initMatchers(payload)

	var wg sync.WaitGroup
	wg.Add(1)

	go scanFolder(".", strings.Contains(payload, "**/*"), &wg)

	wg.Wait()

	success("EOF")
}
