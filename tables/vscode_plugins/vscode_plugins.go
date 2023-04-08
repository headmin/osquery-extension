package vscode_plugins

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/osquery/osquery-go/plugin/table"
)

type Extension struct {
	Name string
	Path string
}

func ExtensionsColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("name"),
		table.TextColumn("path"),
	}
}

var ErrUserHomeDir = errors.New("failed to get user home directory")

func ExtensionsGenerate(ctx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, ErrUserHomeDir
	}

	extensions, err := findExtensions(filepath.Join(homeDir, ".vscode/extensions"))
	if err != nil {
		return nil, err
	}

	results := make([]map[string]string, 0, len(extensions))

	for _, extension := range extensions {
		result := map[string]string{
			"name": extension.Name,
			"path": extension.Path,
		}
		results = append(results, result)
	}

	return results, nil
}

func findExtensions(rootDir string) ([]Extension, error) {
	var extensions []Extension

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		dirName := filepath.Base(path)

		for _, prefix := range []string{"ms-", "vscode-"} {
			if len(dirName) >= len(prefix) && (dirName == prefix || dirName[:len(prefix)] == prefix) {
				extensions = append(extensions, Extension{Name: dirName, Path: filepath.Clean(path)})
				return filepath.SkipDir
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return extensions, nil
}
