package pvent

import (
	"embed"
	"strings"
)

//go:embed VERSION
var version embed.FS

func GetVersion() string {
	v := "0.1.0"

	data, err := version.ReadFile("VERSION")
	if err != nil {
		return v
	}

	v = strings.TrimSpace(string(data))
	return v
}
