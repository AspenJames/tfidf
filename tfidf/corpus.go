package tfidf

import (
	"fmt"
	"math"

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

// Calculates the TF-IDF value of each term in each Document of the Corpus.
// Returns a mapping of Document ID -> mapping of term -> TF-IDF value.
func (c *Corpus) Calculate() (map[uuid.UUID]TermTFIDF, error) {
	for id, document := range c.corpus {
		documentTFIDF, err := c.TFIDFs(document.GetTerms(), id)
		if err != nil {
			return nil, err
		}
		c.tfidfMap[id] = documentTFIDF
	}
	return c.tfidfMap, nil
}

// Calculates TF-IDF for `term` in the Document specified by `docID`.
func (c *Corpus) TFIDF(term string, docID uuid.UUID) (TermTFIDF, error) {
	tfidf := make(TermTFIDF)
	doc, ok := c.corpus[docID]
	if !ok {
		return tfidf, fmt.Errorf("document %s not found in corpus", docID)
	}
	inDocCount := len(c.termInDocs[term])
	totalDocuments := len(c.corpus)

	if inDocCount == 0 {
		tfidf[term] = 0
	} else {
		tf := doc.GetTF(term)
		idf := math.Log(float64(totalDocuments) / float64(inDocCount))
		tfidf[term] = tf * idf
	}

	return tfidf, nil
}

// Calculates TF-IDF for `terms` in the Document specified by `docID`.
func (c *Corpus) TFIDFs(terms []string, docID uuid.UUID) (TermTFIDF, error) {
	tfidf := make(TermTFIDF)
	for _, term := range terms {
		freq, err := c.TFIDF(term, docID)
		if err != nil {
			return tfidf, err
		}
		tfidf[term] = freq[term]
	}
	return tfidf, nil
}
