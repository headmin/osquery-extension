package intellij_plugins

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/osquery/osquery-go/plugin/table"
)

type IntelliJPlugin struct {
	AppName       string
	PluginDirName string
	PluginName    string
	BinaryVersion string
	JarName       string
	FullPath      string
}

func IntelliJPluginsColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("app_name"),
		table.TextColumn("plugin"),
		table.TextColumn("component_name"),
		table.TextColumn("component_version"),
		table.TextColumn("path"),
	}
}

func IntelliJPluginsGenerate(ctx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	userDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	jetbrainsDir := filepath.Join(userDir, "Library/Application Support/JetBrains")

	plugins := []IntelliJPlugin{}

	re := regexp.MustCompile(`^(.*?)-([\d\.]+)\.jar$`)

	err = filepath.Walk(jetbrainsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".jar") {
			pluginDir := filepath.Dir(filepath.Dir(path))
			if filepath.Base(filepath.Dir(pluginDir)) == "plugins" {
				appName := filepath.Base(filepath.Dir(filepath.Dir(pluginDir)))
				matches := re.FindStringSubmatch(filepath.Base(path))
				if len(matches) == 3 {
					plugins = append(plugins, IntelliJPlugin{
						AppName:       appName,
						PluginDirName: filepath.Base(pluginDir),
						PluginName:    filepath.Base(pluginDir),
						BinaryVersion: matches[2],
						JarName:       matches[1],
						FullPath:      path,
					})
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	results := make([]map[string]string, 0, len(plugins))

	for _, plugin := range plugins {
		result := map[string]string{
			"app_name":          plugin.AppName,
			"plugin":            plugin.PluginName,
			"component_version": plugin.BinaryVersion,
			"component_name":    plugin.JarName,
			"path":              plugin.FullPath,
		}
		results = append(results, result)
	}

	return results, nil
}
