package cli

import (
	"embed"
	_ "embed"
)

//go:embed spec/*
var specFS embed.FS

// GetSpecFile reads a file from the embedded spec filesystem
func GetSpecFile(path string) ([]byte, error) {
	return specFS.ReadFile("spec/" + path)
}
