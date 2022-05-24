package tfidf

import (
	"bufio"
	"bytes"
	"io"
	"regexp"

	"github.com/google/uuid"
)

// Document represents a collection of unique words (tokens) in a logic
// grouping of text (document).
type Document struct {
	ID    uuid.UUID
	Meta  Meta
	tfmap map[string]float64 // [term]: frequency
}

// Document metadata map.
type Meta map[string]interface{}

// Process input into a Document.
func Process(input io.Reader, meta Meta) (*Document, error) {
	s := bufio.NewScanner(input)
	counts := make(map[string]int)
	total := 0
	for s.Scan() {
		// Finds words, removes non-word characters.
		wordRe := regexp.MustCompile(`\w+\b`)
		// Find words from line.
		words := wordRe.FindAll(s.Bytes(), -1)
		for _, word := range words {
			total += 1
			s := string(normalize(word))
			if tf, found := counts[s]; !found {
				counts[s] = 1
			} else {
				counts[s] = tf + 1
			}
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	tfmap := make(map[string]float64)
	for word, count := range counts {
		tfmap[word] = float64(count) / float64(total)
	}
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &Document{
		ID:    id,
		tfmap: tfmap,
		Meta:  meta,
	}, nil
}

func normalize(w []byte) []byte {
	return bytes.ToLower(w)
}
