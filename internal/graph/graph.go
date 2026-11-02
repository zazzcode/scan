// Package graph defines source-backed route and string-use analysis artifacts.
//
// The next phase reads this document to judge validation and sanitization.
// It is syntax: route registrations and the handler chain attached to them.
package graph

// SchemaVersion is the JSON contract version.
const SchemaVersion = 1

// Graph is one scan of a repository.
type Graph struct {
	SchemaVersion int     `json:"schema_version"`
	Root          string  `json:"root"`
	Routes        []Route `json:"routes"`
	Gaps          []Gap   `json:"gaps"`
}

// Route is one framework registration with its source-backed handler.
type Route struct {
	ID           string `json:"id"`
	Method       string `json:"method"`
	PathPattern  string `json:"path_pattern"`
	Framework    string `json:"framework"`
	Registration Span   `json:"registration"`
	Handler      Symbol `json:"handler"`
	Steps        []Step `json:"steps"`
}

// Step is a syntactic handler or call record attached to a registration.
type Step struct {
	Kind   string `json:"kind"`
	Symbol Symbol `json:"symbol"`
	Caller string `json:"caller,omitempty"`
}

// Symbol is a named or anonymous function at a source span.
type Symbol struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
	Span Span   `json:"span"`
}

// Span is a source range. Lines are 1-based. Bytes are offsets into the file.
type Span struct {
	File      string `json:"file"`
	StartLine int    `json:"start_line"`
	EndLine   int    `json:"end_line"`
	StartByte int    `json:"start_byte"`
	EndByte   int    `json:"end_byte"`
}

// Gap is a route fact the extractor could not complete.
type Gap struct {
	Kind   string `json:"kind"`
	Detail string `json:"detail"`
	At     Span   `json:"at"`
}
