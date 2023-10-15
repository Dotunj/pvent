package util

import (
	"io"
	"os"
)

func ReadJSONTarget(path string) ([]byte, error) {
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
