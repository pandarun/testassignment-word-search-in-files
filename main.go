package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"runtime"
	"word-search-in-files/pkg/searcher"
)

func main() {

	s := searcher.IndexedSearcher{
		Searcher: &searcher.Searcher{
			FS:          os.DirFS("./examples"),
			WorkerCount: runtime.NumCPU(),
		},
	}

	log.Println("Indexing...")

	err := s.Index()
	if err != nil {
		log.Fatalf("Failed to index: %v", err)
	}

	log.Println("Indexed")

	log.Println("Server started")
	h := NewHandler(&s)
	http.HandleFunc("/files/search", h)
	log.Fatalln(http.ListenAndServe(":8080", nil))
	log.Println("Server stopped")
}

func NewHandler(s searcher.SearcherInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		word := r.URL.Query().Get("word")
		files, err := s.Search(word)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(files)
	}

}
