package indexer

import (
	"io/fs"
)

// Indexer it is a struct that contains a map of words and files
type Indexer struct {
	Index map[string][]string
	FS    fs.FS
}
