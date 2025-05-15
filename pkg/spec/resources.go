package spec

import (
	"embed"
)

//go:embed instacli/instacli-spec/**
var specFS embed.FS

// GetSpecFile reads a file from the embedded Instacli spec filesystem
func GetSpecFile(path string) ([]byte, error) {
	return specFS.ReadFile("instacli/instacli-spec/" + path)
}
