package tfidf

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
)

var meta Meta = make(Meta)

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

func TestProcessDDT(t *testing.T) {
	type ddt struct { // Data-driven test input format.
		data     string
		expected map[string]float64
	}
	type testSuite map[string][]ddt
	suite := testSuite{
		"countsRepeated": {
			{
				data:     "word word repeated",
				expected: map[string]float64{"word": 2.0 / 3.0, "repeated": 1.0 / 3.0},
			},
			{
				data:     "word separated word",
				expected: map[string]float64{"word": 2.0 / 3.0, "separated": 1.0 / 3.0},
			},
		},
		"ignoresCase": {
			{
				data:     "test Test TEST",
				expected: map[string]float64{"test": 3.0 / 3.0},
			},
			{
				data:     "tEsT Test one TWO oNe TwO",
				expected: map[string]float64{"test": 2.0 / 6.0, "one": 2.0 / 6.0, "two": 2.0 / 6.0},
			},
		},
		"removesSymbols": {
			{
				data:     "this && that",
				expected: map[string]float64{"this": 1.0 / 2.0, "that": 1.0 / 2.0},
			},
			{
				data:     "Process removes %#$@! symbols",
				expected: map[string]float64{"process": 1.0 / 3.0, "removes": 1.0 / 3.0, "symbols": 1.0 / 3.0},
			},
			{
				data:     "Also removes, punctuation!!",
				expected: map[string]float64{"also": 1.0 / 3.0, "removes": 1.0 / 3.0, "punctuation": 1.0 / 3.0},
			},
			{
				data:     "'Removes' \"quotations\"",
				expected: map[string]float64{"removes": 1.0 / 2.0, "quotations": 1.0 / 2.0},
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
