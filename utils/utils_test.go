package utils

import (
	"archive/tar"
	"compress/gzip"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTestTarGz(t *testing.T, files map[string]string) string {
	t.Helper()
	tempFile, err := ioutil.TempFile("", "test-*.tar.gz")
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

func TestIsValidElvArchive(t *testing.T) {
	// Valid config.yaml
	config := `version: 1.0.0\nelvdoc: {name: test}`
	files := map[string]string{"elvdoc/config.yaml": config}
	archive := createTestTarGz(t, files)
	defer os.Remove(archive)
	assert.True(t, IsValidElvArchive(archive), "should detect valid elvdoc archive")

	// Invalid config.yaml (missing fields)
	badConfig := `foo: bar`
	files = map[string]string{"elvdoc/config.yaml": badConfig}
	archive2 := createTestTarGz(t, files)
	defer os.Remove(archive2)
	assert.False(t, IsValidElvArchive(archive2), "should not detect archive with invalid config.yaml")

	// No config.yaml
	files = map[string]string{"other.txt": "data"}
	archive3 := createTestTarGz(t, files)
	defer os.Remove(archive3)
	assert.False(t, IsValidElvArchive(archive3), "should not detect archive without config.yaml")
}
