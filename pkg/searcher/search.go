package searcher

import (
	"bufio"
	"errors"
	"io/fs"
	"regexp"
	"runtime"
	"strings"
	"word-search-in-files/pkg/internal/dir"
	"word-search-in-files/pkg/internal/text"
)

type Searcher struct {
	FS          fs.FS
	Directory   string
	WorkerCount int
}

type SearcherInterface interface {
	Search(word string) (files []string, err error)
}

const (
	SearchWordRegex = "^[a-zA-Z0-9_а-яА-Я-]+$"
)

var searchWordRegexp = regexp.MustCompile(SearchWordRegex)

type ValidSearchString struct {
	value string
}

func NewValidSearchString(s string) (ValidSearchString, error) {
	// creates a new valid search string
	// returns an error if the string is empty or invalid
	// returns a valid search string if the string is valid
	// the string is valid if it contains only letters, numbers, and underscores
	// the string is invalid if it contains any other characters
	// the string is invalid if it is empty

	if len(s) == 0 {
		return ValidSearchString{}, errors.New("empty string")
	}

	trimmed := strings.Trim(s, " ")
	if !searchWordRegexp.MatchString(trimmed) {
		return ValidSearchString{}, errors.New("invalid string")
	}

	return ValidSearchString{value: trimmed}, nil

}

func (v *ValidSearchString) String() string {
	return v.value
}

func (s *Searcher) Search(word string) (files []string, err error) {

	validSearchString, err := NewValidSearchString(word)
	if err != nil {
		return nil, err
	}

	dirFiles, err := dir.FilesFS(s.FS, s.Directory)

	if err != nil {
		return nil, err
	}

	workerInput := make(chan string, 2)
	workerOutput := make(chan SearchResult, len(dirFiles))

	for i := 0; i < checkWorkerCount(s); i++ {
		go startWorker(workerInput, workerOutput, validSearchString, s)
	}

	for _, f := range dirFiles {
		workerInput <- f
	}

	close(workerInput)

	for i := 0; i < len(dirFiles); i++ {
		res := <-workerOutput

		if res.Error != nil {
			return nil, res.Error
		}

		if res.Found {
			files = append(files, dir.NewFileNameWithoutExtension(res.FileName).String())
		}
	}

	return files, nil
}

func checkWorkerCount(s *Searcher) int {

	if s.WorkerCount <= 0 {
		return runtime.NumCPU()
	}

	return s.WorkerCount
}

type SearchResult struct {
	FileName string
	Found    bool
	Error    error
}

func startWorker(workerInput <-chan string, workerOutput chan<- SearchResult, word ValidSearchString, s *Searcher) {
	// starts a worker that searches for a word in a file
	// receives a file name from the worker input channel
	// sends the result to the worker output channel
	// the result contains the file name and a boolean value that indicates if the word was found in the file

	for file := range workerInput {
		found, err := s.contains(file, word)
		if err != nil {
			workerOutput <- SearchResult{
				FileName: file,
				Error:    err,
				Found:    false,
			}
			continue
		}

		workerOutput <- SearchResult{
			FileName: file,
			Found:    found,
		}

		// this is a hint to the runtime that the current goroutine is willing to yield its current thread of execution
		runtime.Gosched()
	}
}

func (s *Searcher) contains(f string, searchString ValidSearchString) (found bool, err error) {
	// checks if a file contains a search string
	// returns a boolean value that indicates if the search string was found in the file
	// returns an error if the file cannot be opened or read

	file, err := s.FS.Open(f)
	if err != nil {
		return false, err
	}

	in := bufio.NewReader(file)

	expr := `(^|\s)(?i)` + searchString.String() + `(\s|$)`
	searchRegex := regexp.MustCompile(expr)
	if err != nil {
		return false, err
	}

	found = searchRegex.MatchReader(text.NewCustomRuneReader(in, text.GetDefaultFilteredRunes()...))

	err = file.Close()
	if err != nil {
		return false, err
	}

	return found, nil
}
