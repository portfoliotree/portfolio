package portfolio

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
)

// ParseDocumentFile opens a file and parses the contents into a Document
// It supports YAML files at the moment but may support other encodings in the future.
func ParseDocumentFile(specificationFilePath string) ([]Document, error) {
	if err := checkPortfolioFileName(specificationFilePath); err != nil {
		return nil, err
	}
	f, err := os.Open(specificationFilePath)
	if err != nil {
		return nil, err
	}
	defer closeAndIgnoreErrors(f)
	return portfoliosFromFile(specificationFilePath, f)
}

func checkPortfolioFileName(fileName string) error {
	switch {
	case strings.HasSuffix(fileName, "_portfolio.yml"),
		strings.HasSuffix(fileName, "_portfolio.yaml"):
		return nil
	default:
		return fmt.Errorf("expected a YAML file: it must have a _portfolio.yml file name suffix")
	}
}

func portfoliosFromFile(fileName string, file fs.File) ([]Document, error) {
	result, err := ParseDocuments(file)
	if err != nil {
		return result, err
	}
	// p := filepath.ToSlash(fileName)
	return result, nil
}
