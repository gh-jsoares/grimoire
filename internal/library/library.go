// Package library handles loading and resolving .grim documents from a directory.
package library

import (
	"sort"

	"github.com/gh-jsoares/grimoire/internal/document"
)

// Library holds a collection of parsed documents and any errors encountered during loading.
type Library struct {
	Documents []document.Document
	Errors    []error
}

// Sort orders documents by their Order field (ascending), falling back to alphabetical title.
func Sort(docs []document.Document) {
	sort.SliceStable(docs, func(i, j int) bool {
		oi, oj := docs[i].Order, docs[j].Order
		if oi != nil && oj != nil {
			if *oi != *oj {
				return *oi < *oj
			}
			return docs[i].Title < docs[j].Title
		}
		if oi != nil {
			return true
		}
		if oj != nil {
			return false
		}
		return docs[i].Title < docs[j].Title
	})
}
