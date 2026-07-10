package library

import (
	"sort"

	"github.com/gh-jsoares/grimoire/internal/document"
)

type Library struct {
	Documents []document.Document
	Errors    []error
}

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
