package core

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// IsValidElvArchive checks if a .tar.gz archive contains 'elvdoc/config.yaml' with required fields.
func IsValidElvArchive(path string) bool {
	if !strings.HasSuffix(strings.ToLower(path), ".tar.gz") {
		return false
	}
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return false
	}
	defer gz.Close()
	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return false
		}
		if hdr.Name == "elvdoc/config.yaml" {
			var config struct {
				Version string      `yaml:"version"`
				Elvdoc  interface{} `yaml:"elvdoc"`
			}
			dec := yaml.NewDecoder(tr)
			if dec.Decode(&config) == nil && config.Version != "" && config.Elvdoc != nil {
				return true
			}
			return false
		}
	}
	return false
}

// CreateElvArchive creates a .tar.gz archive containing an elvdoc folder with the provided files.
// The files parameter should contain the content of template.html, style.css, function.js, and config.yaml respectively.
func CreateElvArchive(outputPath string, templateHTML, styleCSS, functionJS, configYAML string) error {
	// Create the output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Create gzip writer
	gzWriter := gzip.NewWriter(outFile)
	defer gzWriter.Close()

	// Create tar writer
	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	// Define files to add to the archive
	files := map[string][]byte{
		"elvdoc/template.html": []byte(templateHTML),
		"elvdoc/style.css":     []byte(styleCSS),
		"elvdoc/function.js":   []byte(functionJS),
		"elvdoc/config.yaml":   []byte(configYAML),
	}

	// Add each file to the tar archive
	for filename, content := range files {
		header := &tar.Header{
			Name:    filename,
			Size:    int64(len(content)),
			Mode:    0644,
			ModTime: time.Now(),
		}

		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		if _, err := tarWriter.Write(content); err != nil {
			return err
		}
	}

	return nil
}

// CreateElvArchiveFromFiles creates a .tar.gz archive from file paths.
// This function reads the files from disk and creates the archive.
func CreateElvArchiveFromFiles(outputPath, templatePath, stylePath, functionPath, configPath string) error {
	// Read all files
	templateHTML, err := os.ReadFile(templatePath)
	if err != nil {
		return err
	}

	styleCSS, err := os.ReadFile(stylePath)
	if err != nil {
		return err
	}

	functionJS, err := os.ReadFile(functionPath)
	if err != nil {
		return err
	}

	configYAML, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	// Create the archive
	return CreateElvArchive(outputPath, string(templateHTML), string(styleCSS), string(functionJS), string(configYAML))
}

// CreateElvArchiveFromDir creates a .tar.gz archive from a directory containing elvdoc files.
// The directory should contain template.html, style.css, function.js, and config.yaml files.
func CreateElvArchiveFromDir(outputPath, sourceDir string) error {
	templatePath := filepath.Join(sourceDir, "template.html")
	stylePath := filepath.Join(sourceDir, "style.css")
	functionPath := filepath.Join(sourceDir, "function.js")
	configPath := filepath.Join(sourceDir, "config.yaml")

	return CreateElvArchiveFromFiles(outputPath, templatePath, stylePath, functionPath, configPath)
}

// ReadFileAsString reads a file and returns its content as a string
func ReadFileAsString(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// ReadElvDocFiles reads all elvdoc files from a directory and returns their content as strings
func ReadElvDocFiles(sourceDir string) (templateHTML, styleCSS, functionJS, configYAML string, err error) {
	templatePath := filepath.Join(sourceDir, "template.html")
	stylePath := filepath.Join(sourceDir, "style.css")
	functionPath := filepath.Join(sourceDir, "function.js")
	configPath := filepath.Join(sourceDir, "config.yaml")

	templateHTML, err = ReadFileAsString(templatePath)
	if err != nil {
		return "", "", "", "", err
	}

	styleCSS, err = ReadFileAsString(stylePath)
	if err != nil {
		return "", "", "", "", err
	}

	functionJS, err = ReadFileAsString(functionPath)
	if err != nil {
		return "", "", "", "", err
	}

	configYAML, err = ReadFileAsString(configPath)
	if err != nil {
		return "", "", "", "", err
	}

	return templateHTML, styleCSS, functionJS, configYAML, nil
}

// TestCreateArchiveFromFormatDir is a test function that reads files from format/elvdoc
// and creates a test archive to verify the functionality
func TestCreateArchiveFromFormatDir(outputPath string) error {
	// Read the files from format/elvdoc directory
	templateHTML, styleCSS, functionJS, configYAML, err := ReadElvDocFiles("format/elvdoc")
	if err != nil {
		return err
	}

	// Create the archive using the string content (simulating database data)
	err = CreateElvArchive(outputPath, templateHTML, styleCSS, functionJS, configYAML)
	if err != nil {
		return err
	}

	// Verify the created archive is valid
	if !IsValidElvArchive(outputPath) {
		return fmt.Errorf("created archive failed validation")
	}

	return nil
}
