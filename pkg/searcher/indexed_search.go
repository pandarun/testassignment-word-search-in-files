package searcher

import (
	"bufio"
	"log"
	"strings"
	"word-search-in-files/pkg/internal/dir"
	"word-search-in-files/pkg/internal/text"
)

type LowercaseString struct {
	value string
}

func NewLowercasedString(s string) LowercaseString {
	return LowercaseString{value: strings.ToLower(s)}
}

func (s *LowercaseString) String() string {
	return s.value
}

type IndexedSearcher struct {
	Searcher *Searcher
	index    map[LowercaseString]map[dir.FileNameWithoutExtension]struct{}
}

func (s *IndexedSearcher) Search(word string) (files []string, err error) {

	var result []string
	if len(s.index) > 0 {
		result, err = s.searchFromIndex(word)
		if err != nil {
			return nil, err
		}
	}

	if len(result) > 0 {
		log.Println("Search from index")
		return result, nil
	}

	log.Println("Search from files")
	return s.Searcher.Search(word)
}

func (s *IndexedSearcher) Index() error {
	// creates an index of files and their contents
	// returns a map where the keys are words and the values are file names where the words were found
	// returns an error if the directory cannot be opened or read

	dirFiles, err := dir.FilesFS(s.Searcher.FS, s.Searcher.Directory)

	if err != nil {
		return err
	}

	s.index = map[LowercaseString]map[dir.FileNameWithoutExtension]struct{}{}

	for _, f := range dirFiles {
		err = s.indexFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *IndexedSearcher) searchFromIndex(word string) ([]string, error) {
	// searches for a word in the index
	// returns a list of file names where the word was found
	// returns nil if the searcherword is not found in the index

	validSearchString, err := NewValidSearchString(word)
	if err != nil {
		return nil, err
	}

	lowerCased := NewLowercasedString(validSearchString.value)

	var res []string
	files, ok := s.index[lowerCased]
	if !ok {
		return nil, nil
	}

	for file := range files {
		res = append(res, file.String())
	}

	return res, nil
}

func (s *IndexedSearcher) indexFile(f string) error {
	// indexes a file
	// returns an error if the file cannot be opened or read

	file, err := s.Searcher.FS.Open(f)

	if err != nil {
		return err
	}

	in := bufio.NewReader(file)

	words, err := text.Words(in, text.GetDefaultFilteredRunes()...)
	if err != nil {
		return err
	}

	for _, word := range words {
		lowerCased := NewLowercasedString(word)
		if _, ok := s.index[lowerCased]; !ok {
			s.index[lowerCased] = make(map[dir.FileNameWithoutExtension]struct{})
		}

		fileName := dir.NewFileNameWithoutExtension(f)
		s.index[lowerCased][fileName] = struct{}{}
	}

	err = file.Close()
	if err != nil {
		return err
	}
	return nil
}
