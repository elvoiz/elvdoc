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
	"path/filepath"

	"github.com/elvoiz/elvdoc/core"
)

// ReadElvDoc reads all elvdoc files from a directory and returns their content as strings.
// It expects the directory to contain template.html, style.css, function.js, and config.yaml.
func ReadElvDoc(dir string) (templateHTML, styleCSS, functionJS, configYAML string, err error) {
	return core.ReadElvDocFiles(dir)
}

// WriteElvDoc writes the provided elvdoc files to a .elv (tar.gz) archive.
// The outputPath should end with .elv or .tar.gz.
func WriteElvDoc(outputPath, templateHTML, styleCSS, functionJS, configYAML string) error {
	// Accept .elv as extension, but store as .tar.gz
	ext := filepath.Ext(outputPath)
	if ext == ".elv" {
		outputPath = outputPath + ".tar.gz"
	}
	return core.CreateElvArchive(outputPath, templateHTML, styleCSS, functionJS, configYAML)
}

// WriteElvDocFromDir writes a .elv (tar.gz) archive from a directory containing elvdoc files.
func WriteElvDocFromDir(outputPath, sourceDir string) error {
	ext := filepath.Ext(outputPath)
	if ext == ".elv" {
		outputPath = outputPath + ".tar.gz"
	}
	return core.CreateElvArchiveFromDir(outputPath, sourceDir)
}

// IsValidElvDocArchive checks if a file is a valid elvdoc archive.
func IsValidElvDocArchive(path string) bool {
	return core.IsValidElvArchive(path)
}

// Example usage:
//  template, css, js, yaml, err := elvdoc.ReadElvDoc("format/elvdoc")
//  err = elvdoc.WriteElvDoc("output.elv", template, css, js, yaml)
//  valid := elvdoc.IsValidElvDocArchive("output.elv.tar.gz")
