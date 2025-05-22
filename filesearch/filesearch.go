package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

// SearchOptions defines the search criteria
type SearchOptions struct {
	NamePattern    string
	SizeMin        int64
	SizeMax        int64
	ContentPattern string
	MaxDepth       int
	CurrentDepth   int
	IgnoreCase     bool
}

func main() {
	// Define command-line flags
	startPath := flag.String("path", ".", "Path to start searching from")
	namePattern := flag.String("name", "", "File name pattern (supports wildcards * and ?)")
	contentPattern := flag.String("content", "", "Content pattern to search for")
	sizeMin := flag.Int64("min-size", -1, "Minimum file size in bytes")
	sizeMax := flag.Int64("max-size", -1, "Maximum file size in bytes")
	maxDepth := flag.Int("depth", -1, "Maximum directory depth to search")
	ignoreCase := flag.Bool("ignore-case", false, "Ignore case in name and content matching")

	// Parse command-line flags
	flag.Parse()

	// Convert wildcards to regexp
	nameRegexp := ""
	if *namePattern != "" {
		nameRegexp = wildcardToRegexp(*namePattern)
	}

	// Create search options
	opts := SearchOptions{
		NamePattern:    nameRegexp,
		SizeMin:        *sizeMin,
		SizeMax:        *sizeMax,
		ContentPattern: *contentPattern,
		MaxDepth:       *maxDepth,
		CurrentDepth:   0,
		IgnoreCase:     *ignoreCase,
	}

	// Create a channel to receive matching files
	matches := make(chan string)

	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup
	wg.Add(1)

	// Start the search in a goroutine
	go func() {
		searchFiles(*startPath, opts, matches, &wg)
		wg.Done()
	}()

	// Start a goroutine to close the channel when all searches are done
	go func() {
		wg.Wait()
		close(matches)
	}()

	// Print matches as they come in
	count := 0
	for match := range matches {
		fmt.Println(match)
		count++
	}

	fmt.Printf("\nFound %d matching files\n", count)
}

// Convert wildcard pattern to regexp
func wildcardToRegexp(pattern string) string {
	pattern = regexp.QuoteMeta(pattern)
	pattern = strings.ReplaceAll(pattern, "\\*", ".*")
	pattern = strings.ReplaceAll(pattern, "\\?", ".")
	return "^" + pattern + "$"
}

// Search for files recursively
func searchFiles(path string, opts SearchOptions, matches chan<- string, wg *sync.WaitGroup) {
	// Check if we've reached the maximum depth
	if opts.MaxDepth > 0 && opts.CurrentDepth >= opts.MaxDepth {
		return
	}

	// List files in the current directory
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading directory %s: %v\n", path, err)
		return
	}

	// Create new options for subdirectories with increased depth
	subOpts := opts
	subOpts.CurrentDepth = opts.CurrentDepth + 1

	// Process each file/directory
	for _, file := range files {
		fullPath := filepath.Join(path, file.Name())

		// If it's a directory, search it recursively
		if file.IsDir() {
			wg.Add(1)
			go func(dirPath string) {
				searchFiles(dirPath, subOpts, matches, wg)
				wg.Done()
			}(fullPath)
			continue
		}

		// Check if the file matches our criteria
		if matchFile(fullPath, file, opts) {
			matches <- fullPath
		}
	}
}

// Check if a file matches the search criteria
func matchFile(path string, fileInfo fs.DirEntry, opts SearchOptions) bool {
	// Get file info
	info, err := fileInfo.Info()
	if err != nil {
		return false
	}

	// Check name pattern if specified
	if opts.NamePattern != "" {
		nameToMatch := info.Name()
		pattern := opts.NamePattern

		if opts.IgnoreCase {
			nameToMatch = strings.ToLower(nameToMatch)
			pattern = strings.ToLower(pattern)
		}

		matched, err := regexp.MatchString(pattern, nameToMatch)
		if err != nil || !matched {
			return false
		}
	}

	// Check file size if min or max is specified
	if opts.SizeMin >= 0 && info.Size() < opts.SizeMin {
		return false
	}
	if opts.SizeMax >= 0 && info.Size() > opts.SizeMax {
		return false
	}

	// Check content pattern if specified
	if opts.ContentPattern != "" {
		return matchContent(path, opts.ContentPattern, opts.IgnoreCase)
	}

	return true
}

// Check if file content matches the pattern
func matchContent(path, pattern string, ignoreCase bool) bool {
	// Read the file
	content, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	contentStr := string(content)
	patternToMatch := pattern

	if ignoreCase {
		contentStr = strings.ToLower(contentStr)
		patternToMatch = strings.ToLower(patternToMatch)
	}

	return strings.Contains(contentStr, patternToMatch)
}
