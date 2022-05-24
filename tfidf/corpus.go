package tfidf

import (
	"github.com/google/uuid"
)

// Corpus is a collection of Documents.
type Corpus struct {
	corpus     map[uuid.UUID]*Document
	termInDocs map[string][]uuid.UUID
	tfidfMap   map[uuid.UUID]TermTFIDF
}

type TermTFIDF map[string]float64

func NewCorpus() *Corpus {
	return &Corpus{
		corpus:     make(map[uuid.UUID]*Document),
		termInDocs: make(map[string][]uuid.UUID),
		tfidfMap:   make(map[uuid.UUID]TermTFIDF),
	}
}

func (c *Corpus) AddDocument(d *Document) {
	c.corpus[d.ID] = d
	for _, term := range d.GetTerms() {
		c.termInDocs[term] = append(c.termInDocs[term], d.ID)
	}
}

func (c *Corpus) AddDocuments(ds []*Document) {
	for _, document := range ds {
		c.AddDocument(document)
	}
}
