package core

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTestTarGz(t *testing.T, files map[string]string) string {
	t.Helper()
	tempFile, err := os.CreateTemp("", "test-*.tar.gz")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer tempFile.Close()
	gz := gzip.NewWriter(tempFile)
	tr := tar.NewWriter(gz)
	for name, content := range files {
		hdr := &tar.Header{
			Name: name,
			Mode: 0600,
			Size: int64(len(content)),
		}
		if err := tr.WriteHeader(hdr); err != nil {
			t.Fatalf("failed to write tar header: %v", err)
		}
		if _, err := tr.Write([]byte(content)); err != nil {
			t.Fatalf("failed to write tar content: %v", err)
		}
	}
	tr.Close()
	gz.Close()
	return tempFile.Name()
}

func TestCreateElvArchive(t *testing.T) {
	// Test data simulating database content
	templateHTML := `<!DOCTYPE html>
<html>
<head><title>Test Template</title></head>
<body><h1>Hello World</h1></body>
</html>`

	styleCSS := `body { font-family: Arial, sans-serif; color: #333; }`

	functionJS := `function hello() { console.log("Hello from JS!"); }`

	configYAML := `version: "1.0.0"
elvdoc:
  name: "test-package"
  description: "Test package for unit testing"`

	// Create temporary file for output
	tempFile, err := os.CreateTemp("", "test-archive-*.tar.gz")
	assert.NoError(t, err)
	tempFile.Close()
	defer os.Remove(tempFile.Name())

	// Test CreateElvArchive with string content
	err = CreateElvArchive(tempFile.Name(), templateHTML, styleCSS, functionJS, configYAML)
	assert.NoError(t, err, "CreateElvArchive should succeed")

	// Verify the archive was created and is valid
	assert.True(t, IsValidElvArchive(tempFile.Name()), "created archive should be valid")

	// Verify the archive contains the correct files with correct content
	verifyArchiveContent(t, tempFile.Name(), templateHTML, styleCSS, functionJS, configYAML)
}

func TestReadFileAsString(t *testing.T) {
	// Create a temporary file with test content
	tempFile, err := os.CreateTemp("", "test-*.txt")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	testContent := "Hello, this is test content!"
	_, err = tempFile.WriteString(testContent)
	assert.NoError(t, err)
	tempFile.Close()

	// Test reading the file as string
	content, err := ReadFileAsString(tempFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, testContent, content)

	// Test reading non-existent file
	_, err = ReadFileAsString("non-existent-file.txt")
	assert.Error(t, err)
}

func TestReadElvDocFiles(t *testing.T) {
	// Create temporary directory structure
	tempDir, err := os.MkdirTemp("", "elvdoc-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Test file contents
	templateContent := "<html><body>Test Template</body></html>"
	styleContent := "body { margin: 0; }"
	functionContent := "console.log('test');"
	configContent := "version: 1.0.0\nelvdoc: {name: test}"

	// Create test files
	files := map[string]string{
		"template.html": templateContent,
		"style.css":     styleContent,
		"function.js":   functionContent,
		"config.yaml":   configContent,
	}

	for filename, content := range files {
		filePath := tempDir + string(os.PathSeparator) + filename
		err := ioutil.WriteFile(filePath, []byte(content), 0644)
		assert.NoError(t, err)
	}

	// Test reading all files
	template, style, function, config, err := ReadElvDocFiles(tempDir)
	assert.NoError(t, err)
	assert.Equal(t, templateContent, template)
	assert.Equal(t, styleContent, style)
	assert.Equal(t, functionContent, function)
	assert.Equal(t, configContent, config)

	// Test with missing file
	os.Remove(tempDir + string(os.PathSeparator) + "template.html")
	_, _, _, _, err = ReadElvDocFiles(tempDir)
	assert.Error(t, err)
}

func TestCreateElvArchiveFromFiles(t *testing.T) {
	// Create temporary directory and files
	tempDir, err := os.MkdirTemp("", "elvdoc-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Test file contents
	templateContent := "<html><body>Test Template</body></html>"
	styleContent := "body { margin: 0; }"
	functionContent := "console.log('test');"
	configContent := "version: 1.0.0\nelvdoc: {name: test}"

	// Create test files
	files := map[string]string{
		"template.html": templateContent,
		"style.css":     styleContent,
		"function.js":   functionContent,
		"config.yaml":   configContent,
	}

	filePaths := make(map[string]string)
	for filename, content := range files {
		filePath := tempDir + string(os.PathSeparator) + filename
		err := ioutil.WriteFile(filePath, []byte(content), 0644)
		assert.NoError(t, err)
		filePaths[filename] = filePath
	}

	// Create output archive
	tempFile, err := os.CreateTemp("", "test-archive-*.tar.gz")
	assert.NoError(t, err)
	tempFile.Close()
	defer os.Remove(tempFile.Name())

	// Test CreateElvArchiveFromFiles
	err = CreateElvArchiveFromFiles(
		tempFile.Name(),
		filePaths["template.html"],
		filePaths["style.css"],
		filePaths["function.js"],
		filePaths["config.yaml"],
	)
	assert.NoError(t, err)

	// Verify the archive
	assert.True(t, IsValidElvArchive(tempFile.Name()))
	verifyArchiveContent(t, tempFile.Name(), templateContent, styleContent, functionContent, configContent)
}

func TestCreateElvArchiveFromDir(t *testing.T) {
	// Create temporary directory with elvdoc files
	tempDir, err := os.MkdirTemp("", "elvdoc-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Test file contents
	templateContent := "<html><body>Test Template</body></html>"
	styleContent := "body { margin: 0; }"
	functionContent := "console.log('test');"
	configContent := "version: 1.0.0\nelvdoc: {name: test}"

	// Create test files
	files := map[string]string{
		"template.html": templateContent,
		"style.css":     styleContent,
		"function.js":   functionContent,
		"config.yaml":   configContent,
	}

	for filename, content := range files {
		filePath := tempDir + string(os.PathSeparator) + filename
		err := ioutil.WriteFile(filePath, []byte(content), 0644)
		assert.NoError(t, err)
	}

	// Create output archive
	tempFile, err := os.CreateTemp("", "test-archive-*.tar.gz")
	assert.NoError(t, err)
	tempFile.Close()
	defer os.Remove(tempFile.Name())

	// Test CreateElvArchiveFromDir
	err = CreateElvArchiveFromDir(tempFile.Name(), tempDir)
	assert.NoError(t, err)

	// Verify the archive
	assert.True(t, IsValidElvArchive(tempFile.Name()))
	verifyArchiveContent(t, tempFile.Name(), templateContent, styleContent, functionContent, configContent)
}

// Helper function to verify archive content
func verifyArchiveContent(t *testing.T, archivePath, expectedTemplate, expectedStyle, expectedFunction, expectedConfig string) {
	t.Helper()

	// Open and read the archive
	f, err := os.Open(archivePath)
	assert.NoError(t, err)
	defer f.Close()

	gz, err := gzip.NewReader(f)
	assert.NoError(t, err)
	defer gz.Close()

	tr := tar.NewReader(gz)

	foundFiles := make(map[string]string)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		assert.NoError(t, err)

		content, err := ioutil.ReadAll(tr)
		assert.NoError(t, err)

		foundFiles[hdr.Name] = string(content)
	}

	// Verify all expected files are present with correct content
	assert.Equal(t, expectedTemplate, foundFiles["elvdoc/template.html"], "template.html content should match")
	assert.Equal(t, expectedStyle, foundFiles["elvdoc/style.css"], "style.css content should match")
	assert.Equal(t, expectedFunction, foundFiles["elvdoc/function.js"], "function.js content should match")
	assert.Equal(t, expectedConfig, foundFiles["elvdoc/config.yaml"], "config.yaml content should match")

	// Verify we have exactly 4 files
	assert.Len(t, foundFiles, 4, "archive should contain exactly 4 files")
}

// TestCreateElvArchiveToDisk creates a non-temporary archive file for manual inspection
func TestCreateElvArchiveToDisk(t *testing.T) {
	outputPath := "test-elvdoc-archive.elv"
	templateHTML := `<!DOCTYPE html>\n<html>\n<head><title>Manual Test</title></head>\n<body><h1>Manual Archive</h1></body>\n</html>`
	styleCSS := `body { background: #fafafa; color: #222; }`
	functionJS := `function greet() { alert(\"Hello Manual!\"); }`
	configYAML := `version: \"1.0.0\"\nelvdoc:\n  name: \"manual-test\"\n  description: \"Archive for manual inspection\"`

	err := CreateElvArchive(outputPath, templateHTML, styleCSS, functionJS, configYAML)
	assert.NoError(t, err, "CreateElvArchive should succeed for manual file")
	// Do not remove the file, user will check it manually
}
