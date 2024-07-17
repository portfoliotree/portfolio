package portfolio

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
)

// ParseSpecificationFile opens a file and parses the contents into a Specification
// It supports YAML files at the moment but may support other encodings in the future.
func ParseSpecificationFile(specificationFilePath string) ([]Document, error) {
	if err := checkPortfolioFileName(specificationFilePath); err != nil {
		return nil, err
	}
	f, err := os.Open(specificationFilePath)
	if err != nil {
		return nil, err
	}
	defer closeAndIgnoreError(f)
	return portfoliosFromFile(f)
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

func portfoliosFromFile(file fs.File) ([]Document, error) {
	result, err := ParseDocuments(file)
	if err != nil {
		return result, err
	}
	return result, nil
}

func WalkDirectoryAndParseSpecificationFiles(dir fs.FS) ([]Document, error) {
	var result []Document
	return result, fs.WalkDir(dir, ".", func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if filePath != "." && strings.HasPrefix(path.Base(filePath), ".") {
				return fs.SkipDir
			}
			return nil
		}
		if err := checkPortfolioFileName(filePath); err != nil {
			return nil
		}
		f, err := dir.Open(filePath)
		if err != nil {
			return err
		}
		defer closeAndIgnoreError(f)
		specs, err := portfoliosFromFile(f)
		if err != nil {
			return err
		}
		result = append(result, specs...)
		return nil
	})
}
