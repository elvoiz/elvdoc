package utils

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"strings"

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
