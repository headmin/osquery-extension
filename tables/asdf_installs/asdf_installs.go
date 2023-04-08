package asdf_installs

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/osquery/osquery-go/plugin/table"
)

type AsdfBinary struct {
	Name    string
	Version string
	Path    string
}

func AsdfInstallsColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("name"),
		table.TextColumn("version"),
		table.TextColumn("path"),
	}
}

func AsdfInstallsGenerate(ctx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	userDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	asdfDir := filepath.Join(userDir, ".asdf", "installs")

	binaries, err := findAsdfBinaries(asdfDir)
	if err != nil {
		return nil, err
	}

	results := make([]map[string]string, 0, len(binaries))

	for _, binary := range binaries {
		result := map[string]string{
			"name":    binary.Name,
			"version": binary.Version,
			"path":    binary.Path,
		}
		results = append(results, result)
	}

	return results, nil
}

func findAsdfBinaries(rootDir string) ([]AsdfBinary, error) {
	var binaries []AsdfBinary

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}

		if strings.HasSuffix(info.Name(), "bin") {
			versionDir := filepath.Dir(path)
			nameDir := filepath.Dir(versionDir)

			binaryName := filepath.Base(nameDir)
			binaryVersion := filepath.Base(versionDir)

			binary := AsdfBinary{
				Name:    binaryName,
				Version: binaryVersion,
				Path:    filepath.Join(path, ".."),
			}

			binaries = append(binaries, binary)
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return binaries, nil
}
