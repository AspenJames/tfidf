package tfidf

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func isInCorpus(d *Document, c map[uuid.UUID]*Document) bool {
	return c[d.ID] == d
}

func TestNewCorpus(t *testing.T) {
	corpus := NewCorpus()
	if reflect.TypeOf(corpus).String() != "*tfidf.Corpus" {
		t.Errorf("expected type Corpus, got %T", corpus)
	}
}

func TestAddDocument(t *testing.T) {
	corpus := NewCorpus()
	document := documentFactory(TF{"term": 1})

	corpus.AddDocument(document)

	if len(corpus.corpus) != 1 {
		t.Errorf("unexpected corpus length %d", len(corpus.corpus))
	}
	if !isInCorpus(document, corpus.corpus) {
		t.Errorf("document id %s not in corpus %v", document.ID, corpus.corpus)
	}
}

func TestAddDocumentIdempotency(t *testing.T) {
	corpus := NewCorpus()
	document := documentFactory(TF{"term": 1})

	corpus.AddDocument(document)
	corpus.AddDocument(document)

	if len(corpus.corpus) != 1 {
		t.Errorf("unexpected corpus length %d", len(corpus.corpus))
	}
	if !isInCorpus(document, corpus.corpus) {
		t.Errorf("document id %s not in corpus %v", document.ID, corpus.corpus)
	}
}

func TestAddDocuments(t *testing.T) {
	corpus := NewCorpus()
	d1 := documentFactory(TF{"term": 1})
	d2 := documentFactory(TF{"term": 1})
	documents := []*Document{d1, d2}

	corpus.AddDocuments(documents)

	if len(corpus.corpus) != 2 {
		t.Errorf("unexpected corpus length %d", len(corpus.corpus))
	}

	for _, document := range documents {
		if !isInCorpus(document, corpus.corpus) {
			t.Errorf("document id %s not in corpus %v", document.ID, corpus.corpus)
		}
	}
}
