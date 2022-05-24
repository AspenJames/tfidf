# TFIDF

This is an implementation of the TF-IDF algorithm for processing text. TF-IDF
stands for "term frequency - inverse document frequency", and is calculated by
multiplying the frequency of a term within a corpus of documents by the inverse
of the frequency within documents.

## License

This package is licensed under the [MIT License](./LICENSE).

## Approach

This package uses the algorithm as described in [this document][tfidf]:

    tfidf(term) = tf(term) * idf(term)
    tf(term) = (count term) / (total terms)
    idf(term) = log((count documents) / (count documents containing term))

## Concepts

This package is organized around two main entities/types -- Document & Corpus.

### Document

A Document represents a logical unit of text, e.g. a piece of user feedback.
Documents have a UUID identifier and may have arbitrary metadata. `Process`
accepts input as an `io.Reader` and a map of metadata, returning a structured
Document.

### Corpus

A Corpus is a collection of related Documents. As Documents are added to the
Corpus, the Document's terms are collected into a map containing the term's
frequency within the entire Corpus, as well as a list of IDs for Documents that
contain the term. Once we have Documents collected into a Corpus, we can
calculate the TFIDF for each term with `Corpus.Calculate()`, populating a map of
term -> tfidf.

[tfidf]: http://tfidf.com
