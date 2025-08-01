// elvdoc is a Go package for reading and writing `.elv` files.
// It provides functions to parse, generate, and manipulate files with the `.elv` extension.
// This package is mainly used in the Elvoiz app for invoice management, but can be integrated into any Go project.
//
// Features:
//   - Read `.elv` files
//   - Write `.elv` files
//   - Parse and generate structured data
//   - Easy-to-use API

package elvdoc

import (
	"io"
	"os"
)

// Document represents the structured data of an .elv file.
type Document struct {
	// TODO: Define fields according to .elv file structure
	// Example:
	// ID      string
	// Content string
}

// Reader is the interface for reading .elv files.
type Reader interface {
	Read(r io.Reader) (*Document, error)
}

// Writer is the interface for writing .elv files.
type Writer interface {
	Write(w io.Writer, doc *Document) error
}

// ReadFile reads an .elv file from the given path and returns a Document.
func ReadFile(path string) (*Document, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Read(f)
}

// WriteFile writes a Document to an .elv file at the given path.
func WriteFile(path string, doc *Document) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return Write(f, doc)
}

// Read parses an .elv file from an io.Reader.
func Read(r io.Reader) (*Document, error) {
	// TODO: Implement parsing logic
	return &Document{}, nil
}

// Write serializes a Document to an io.Writer as an .elv file.
func Write(w io.Writer, doc *Document) error {
	// TODO: Implement serialization logic
	return nil
}
