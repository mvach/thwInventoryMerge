package app

import (
	"fmt"
	"os"
	"strings"
	"thwInventoryMerge/utils"

	"github.com/gogs/chardet"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
)

type encodingProvider interface {
	GetFileEncoding(filePath string) (encoding.Encoding, error)
}

type encodingProviderImpl struct {
	logger utils.Logger
}

func NewEncodingProvider(logger utils.Logger) encodingProvider {
	return &encodingProviderImpl{
		logger: logger,
	}
}

func (e *encodingProviderImpl) GetFileEncoding(filePath string) (encoding.Encoding, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file '%s': %w", filePath, err)
	}

	result, err := chardet.NewTextDetector().DetectBest(data)
	if err != nil {
		return nil, fmt.Errorf("failed to detect encoding of file '%s': %w", filePath, err)
	}

	enc, err := e.getEncodingByName(result.Charset)
	if err != nil {
		return nil, fmt.Errorf("failed to get encoding: %w", err)
	}

	e.logger.Info(fmt.Sprintf("File %s has encoding: %s", filePath, result.Charset))

	return enc, nil
}

func (e *encodingProviderImpl) getEncodingByName(name string) (encoding.Encoding, error) {
	switch strings.ToLower(name) {
	case "utf-8":
		return unicode.UTF8, nil
	case "iso-8859-1":
		return charmap.ISO8859_1, nil
	default:
		return nil, fmt.Errorf("unsupported encoding: %s", name)
	}
}
