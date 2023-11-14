package util

import (
	"errors"
	"io"
	"os"
)

func ReadJSONTarget(path string) ([]byte, error) {
	if IsStringEmpty(path) {
		return nil, errors.New("target path cannot be empty")
	}

	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	payload, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	return payload, nil
}
