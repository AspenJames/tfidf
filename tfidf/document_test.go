package tfidf

import (
	"bytes"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/google/uuid"
)

var meta Meta = make(Meta)

// Shortcut for creating a Document, byo frequency map.
func documentFactory(tfmap TF) *Document {
	return &Document{
		ID:    uuid.Must(uuid.NewRandom()),
		Meta:  Meta{},
		tfmap: tfmap,
	}
}

func termListsEqual(expected, given []string) bool {
	if len(expected) != len(given) {
		return false
	}
	sort.Strings(expected)
	sort.Strings(given)
	for i, str := range given {
		if str != expected[i] {
			return false
		}
	}
	return true
}
func TestProcessAssignsId(t *testing.T) {
	input := strings.NewReader("input")
	document, err := Process(input, meta)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(document.ID).String() != "uuid.UUID" {
		t.Errorf("expected type UUID, got %T", document.ID)
	}
}

func TestProcessSetsMeta(t *testing.T) {
	input := strings.NewReader("input")
	expectedMeta := Meta{
		"integer": 123,
		"key":     "value",
		"nested": Meta{
			"hi": "there",
		},
	}

	document, err := Process(input, expectedMeta)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expectedMeta, document.Meta) {
		t.Errorf("expected %v, got %v", expectedMeta, document.Meta)
	}
}

func TestProcessContent(t *testing.T) {
	input := []byte("input")
	inputReader := bytes.NewReader(input)
	document, err := Process(inputReader, meta)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(input, document.content) {
		t.Logf("size expected: %d, size got: %d", len(input), len(document.content))
		t.Errorf("\nexpected: %v\ngot: %v", input, document.content)
	}
}

func TestProcessContentFromFile(t *testing.T) {
	f, err := os.Open("testdata/lorem.txt")
	if err != nil {
		t.Error(err)
	}
	expected, err := ioutil.ReadFile("testdata/lorem.txt")
	if err != nil {
		t.Error(err)
	}
	// Remove final newline.
	expected = expected[:len(expected)-1]
	defer f.Close()
	document, err := Process(f, meta)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(expected, document.content) {
		t.Logf("size expected: %d, size got: %d", len(expected), len(document.content))
		t.Errorf("\nexpected: %v\ngot: %v\n", expected, document.content)
	}
}

func TestProcessDDT(t *testing.T) {
	type ddt struct { // Data-driven test input format.
		data     string
		expected TF
	}
	type testSuite map[string][]ddt
	suite := testSuite{
		"countsRepeated": {
			{
				data:     "word word repeated",
				expected: TF{"word": 2.0 / 3.0, "repeated": 1.0 / 3.0},
			},
			{
				data:     "word separated word",
				expected: TF{"word": 2.0 / 3.0, "separated": 1.0 / 3.0},
			},
		},
		"ignoresCase": {
			{
				data:     "test Test TEST",
				expected: TF{"test": 3.0 / 3.0},
			},
			{
				data:     "tEsT Test one TWO oNe TwO",
				expected: TF{"test": 2.0 / 6.0, "one": 2.0 / 6.0, "two": 2.0 / 6.0},
			},
		},
		"removesSymbols": {
			{
				data:     "this && that",
				expected: TF{"this": 1.0 / 2.0, "that": 1.0 / 2.0},
			},
			{
				data:     "Process removes %#$@! symbols",
				expected: TF{"process": 1.0 / 3.0, "removes": 1.0 / 3.0, "symbols": 1.0 / 3.0},
			},
			{
				data:     "Also removes, punctuation!!",
				expected: TF{"also": 1.0 / 3.0, "removes": 1.0 / 3.0, "punctuation": 1.0 / 3.0},
			},
			{
				data:     "'Removes' \"quotations\"",
				expected: TF{"removes": 1.0 / 2.0, "quotations": 1.0 / 2.0},
			},
		},
	}

	for label, tests := range suite {
		t.Run(label, func(t *testing.T) {
			for i, test := range tests {
				t.Run(strconv.Itoa(i), func(t *testing.T) {
					input := strings.NewReader(test.data)
					document, err := Process(input, meta)
					if err != nil {
						t.Error(err)
					}

					if !reflect.DeepEqual(test.expected, document.tfmap) {
						t.Errorf("incorrect tfmap, expected %v, got %v", test.expected, document.tfmap)
					}
				})
			}
		})
	}
}

func TestGetTF(t *testing.T) {
	type ddt struct {
		tfmap    TF
		term     string
		expected float64
	}
	type testSuite map[string]ddt
	suite := testSuite{
		"termInMap": {
			tfmap:    TF{"termA": 0.7, "termB": 0.3},
			term:     "termA",
			expected: 0.7,
		},
		"termNotInMap": {
			tfmap:    TF{"termA": 0.7, "termB": 0.3},
			term:     "termC",
			expected: 0.0,
		},
	}
	for label, test := range suite {
		t.Run(label, func(t *testing.T) {
			document := &Document{
				tfmap: test.tfmap,
			}
			tf := document.GetTF(test.term)
			if tf != test.expected {
				t.Errorf("expected %f, got %f", test.expected, tf)
			}
		})
	}
}

func TestGetTerms(t *testing.T) {
	type ddt struct {
		tfmap    TF
		expected []string
	}
	type testSuite map[string]ddt
	suite := testSuite{
		"singleTerm": {
			tfmap:    TF{"term": 1.0},
			expected: []string{"term"},
		},
		"manyTerms": {
			tfmap:    TF{"termA": 0.4, "termB": 0.2, "termC": 0.5},
			expected: []string{"termA", "termB", "termC"},
		},
		"noTerms": {
			tfmap:    TF{},
			expected: []string{},
		},
	}
	for label, test := range suite {
		t.Run(label, func(t *testing.T) {
			document := documentFactory(test.tfmap)
			terms := document.GetTerms()
			if !termListsEqual(test.expected, terms) {
				t.Errorf("expected %v, got %v", test.expected, terms)
			}
		})
	}
}
