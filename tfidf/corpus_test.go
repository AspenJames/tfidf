package tfidf

import (
	"math"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

// Creates a Corpus from the given term frequency maps, returns *Corpus & the
// first Document's ID.
func createCorpus(tfmaps []TF) (*Corpus, uuid.UUID) {
	corpus := NewCorpus()
	documents := []*Document{}
	for _, tfmap := range tfmaps {
		documents = append(documents, documentFactory(tfmap))
	}
	corpus.AddDocuments(documents)
	return corpus, documents[0].ID
}

func isInCorpus(d *Document, c map[uuid.UUID]*Document) bool {
	return c[d.ID] == d
}

func tfidfsEqual(expected, given TermTFIDF) bool {
	if len(expected) != len(given) {
		return false
	}
	for term, freq := range given {
		if freq != expected[term] {
			return false
		}
	}
	return true
}

func tfidfMapsEqual(expected, given map[uuid.UUID]TermTFIDF) bool {
	if len(expected) != len(given) {
		return false
	}
	for id, tfidf := range given {
		if !tfidfsEqual(expected[id], tfidf) {
			return false
		}
	}
	return true
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

func TestTFIDFDocumentNotFound(t *testing.T) {
	document := documentFactory(TF{})
	corpus := NewCorpus()
	corpus.AddDocument(document)

	garbageId := uuid.Must(uuid.NewRandom())
	_, err := corpus.TFIDF("term", garbageId)
	if err == nil {
		t.Errorf("expected 'not found in corpus' error")
	}
}

func TestTFIDF(t *testing.T) {
	type ddt struct {
		tfmaps   []TF
		term     string
		expected TermTFIDF
	}
	type testSuite map[string]ddt
	suite := testSuite{
		"termNotInCorpus": {
			tfmaps:   []TF{{}},
			term:     "term",
			expected: TermTFIDF{"term": 0.0},
		},
		"termInAllDocuments": {
			tfmaps: []TF{
				{"term": 1.0},
			},
			term:     "term",
			expected: TermTFIDF{"term": 0.0},
		},
		"termIsUnique": {
			tfmaps: []TF{
				{"termA": 1.0},
				{"termB": 1.0},
			},
			term:     "termA",
			expected: TermTFIDF{"termA": (1.0 * math.Log(2.0/1.0))},
		},
		"termIsNotUnique": {
			tfmaps: []TF{
				{"termA": 1.0},
				{"termA": 1.0},
				{"termB": 1.0},
			},
			term:     "termA",
			expected: TermTFIDF{"termA": (1.0 * math.Log(3.0/2.0))},
		},
	}
	for label, test := range suite {
		t.Run(label, func(t *testing.T) {
			corpus, targetID := createCorpus(test.tfmaps)

			tfidf, err := corpus.TFIDF(test.term, targetID)
			if err != nil {
				t.Error(err)
			}
			if !tfidfsEqual(test.expected, tfidf) {
				t.Errorf("expected %v, got %v", test.expected, tfidf)
			}
		})
	}
}

func TestTFIDFs(t *testing.T) {
	type ddt struct {
		tfmaps   []TF
		terms    []string
		expected TermTFIDF
	}
	type testSuite map[string]ddt
	suite := testSuite{
		"termsNotInCorpus": {
			tfmaps:   []TF{{}},
			terms:    []string{"termA", "termB"},
			expected: TermTFIDF{"termA": 0.0, "termB": 0.0},
		},
		"termsInCorpus": {
			tfmaps: []TF{
				{"termA": 0.5, "termB": 0.5},
				{"termB": 1.0},
			},
			terms: []string{"termA", "termB"},
			expected: TermTFIDF{
				"termA": (0.5 * math.Log(2.0/1.0)),
				"termB": (0.5 * math.Log(2.0/2.0)),
			},
		},
		"mixedPresence": {
			tfmaps: []TF{
				{"termA": 1.0},
				{"termB": 1.0},
			},
			terms: []string{"termA", "termC"},
			expected: TermTFIDF{
				"termA": (1.0 * math.Log(2.0/1.0)),
				"termC": 0.0},
		},
	}
	for label, test := range suite {
		t.Run(label, func(t *testing.T) {
			corpus, targetID := createCorpus(test.tfmaps)

			tfidf, err := corpus.TFIDFs(test.terms, targetID)
			if err != nil {
				t.Error(err)
			}
			if !tfidfsEqual(test.expected, tfidf) {
				t.Errorf("expected %v, got %v", test.expected, tfidf)
			}
		})
	}
}

func TestCalculateEmptyCorpus(t *testing.T) {
	corpus := NewCorpus()
	tfidf, err := corpus.Calculate()
	if err != nil {
		t.Error(err)
	}
	if !tfidfMapsEqual(map[uuid.UUID]TermTFIDF{}, tfidf) {
		t.Errorf("expected empty map, got %v", tfidf)
	}
}

func TestCalculateSingleDocumentNoTerms(t *testing.T) {
	corpus := NewCorpus()
	document := documentFactory(TF{})
	corpus.AddDocument(document)

	expected := make(map[uuid.UUID]TermTFIDF)
	expected[document.ID] = TermTFIDF{}

	tfidf, err := corpus.Calculate()
	if err != nil {
		t.Error(err)
	}

	if !tfidfMapsEqual(expected, tfidf) {
		t.Errorf("expected %v, got %v", expected, tfidf)
	}
}

func TestCalculateSingleDocumentWithTerms(t *testing.T) {
	corpus := NewCorpus()
	document := documentFactory(TF{"term": 1.0})
	corpus.AddDocument(document)

	expected := make(map[uuid.UUID]TermTFIDF)
	expected[document.ID] = TermTFIDF{"term": 0.0}

	tfidf, err := corpus.Calculate()
	if err != nil {
		t.Error(err)
	}

	if !tfidfMapsEqual(expected, tfidf) {
		t.Errorf("expected %v, got %v", expected, tfidf)
	}
}

func TestCalculateMultipleDocumentsWithTermsSimple(t *testing.T) {
	corpus := NewCorpus()
	documentA := documentFactory(TF{"termA": 1.0})
	documentB := documentFactory(TF{"termB": 1.0})
	corpus.AddDocument(documentA)
	corpus.AddDocument(documentB)

	expected := make(map[uuid.UUID]TermTFIDF)
	expected[documentA.ID] = TermTFIDF{"termA": 1.0 * math.Log(2.0/1.0)}
	expected[documentB.ID] = TermTFIDF{"termB": 1.0 * math.Log(2.0/1.0)}

	tfidf, err := corpus.Calculate()
	if err != nil {
		t.Error(err)
	}

	if !tfidfMapsEqual(expected, tfidf) {
		t.Errorf("expected %v, got %v", expected, tfidf)
	}
}

func TestCalculateMultipleDocumentsWithTermsComplex(t *testing.T) {
	corpus := NewCorpus()
	documentA := documentFactory(TF{"termA": 0.3, "termB": 0.7})
	documentB := documentFactory(TF{"termB": 0.2, "termC": 0.8})
	documentC := documentFactory(TF{"termD": 1.0})
	corpus.AddDocument(documentA)
	corpus.AddDocument(documentB)
	corpus.AddDocument(documentC)

	expected := make(map[uuid.UUID]TermTFIDF)
	expected[documentA.ID] = TermTFIDF{
		"termA": 0.3 * math.Log(3.0/1.0),
		"termB": 0.7 * math.Log(3.0/2.0),
	}
	expected[documentB.ID] = TermTFIDF{
		"termB": 0.2 * math.Log(3.0/2.0),
		"termC": 0.8 * math.Log(3.0/1.0),
	}
	expected[documentC.ID] = TermTFIDF{
		"termD": 1.0 * math.Log(3.0/1.0),
	}

	tfidf, err := corpus.Calculate()
	if err != nil {
		t.Error(err)
	}

	if !tfidfMapsEqual(expected, tfidf) {
		t.Errorf("expected %v, got %v", expected, tfidf)
	}
}
