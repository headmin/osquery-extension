package main

import (
	"flag"
	"log"
	"runtime"
	"time"

	asdf_installs "github.com/macadmins/osquery-extension/tables/asdf_installs"
	"github.com/macadmins/osquery-extension/tables/chromeuserprofiles"
	"github.com/macadmins/osquery-extension/tables/fileline"
	"github.com/macadmins/osquery-extension/tables/filevaultusers"

	intellij_plugins "github.com/macadmins/osquery-extension/tables/intellij_plugins"
	macosprofiles "github.com/macadmins/osquery-extension/tables/macos_profiles"
	"github.com/macadmins/osquery-extension/tables/mdm"
	"github.com/macadmins/osquery-extension/tables/munki"
	"github.com/macadmins/osquery-extension/tables/networkquality"
	"github.com/macadmins/osquery-extension/tables/pendingappleupdates"
	"github.com/macadmins/osquery-extension/tables/puppet"
	"github.com/macadmins/osquery-extension/tables/unifiedlog"
	vscode_plugins "github.com/macadmins/osquery-extension/tables/vscode_plugins"
	osquery "github.com/osquery/osquery-go"
	"github.com/osquery/osquery-go/plugin/table"
)

func main() {
	var (
		flSocketPath = flag.String("socket", "", "")
		flTimeout    = flag.Int("timeout", 0, "")
		_            = flag.Int("interval", 0, "")
		_            = flag.Bool("verbose", false, "")
	)
	flag.Parse()

	// allow for osqueryd to create the socket path otherwise it will error
	time.Sleep(3 * time.Second)

	server, err := osquery.NewExtensionManagerServer(
		"macadmins_extension",
		*flSocketPath,
		osquery.ServerTimeout(time.Duration(*flTimeout)*time.Second),
	)
	if err != nil {
		log.Fatalf("Error creating extension: %s\n", err)
	}

	// Create and register a new table plugin with the server.
	// Adding a new table? Add it to the list and the loop below will handle
	// the registration for you.
	plugins := []osquery.OsqueryPlugin{
		table.NewPlugin("puppet_info", puppet.PuppetInfoColumns(), puppet.PuppetInfoGenerate),
		table.NewPlugin("puppet_logs", puppet.PuppetLogsColumns(), puppet.PuppetLogsGenerate),
		table.NewPlugin("puppet_state", puppet.PuppetStateColumns(), puppet.PuppetStateGenerate),
		table.NewPlugin("puppet_facts", puppet.PuppetFactsColumns(), puppet.PuppetFactsGenerate),
		table.NewPlugin("google_chrome_profiles", chromeuserprofiles.GoogleChromeProfilesColumns(), chromeuserprofiles.GoogleChromeProfilesGenerate),
		table.NewPlugin("file_lines", fileline.FileLineColumns(), fileline.FileLineGenerate),
	}

	// Platform specific tables
	// if runtime.GOOS == "windows" {
	// If there were windows only tables, they would go here
	// }

	if runtime.GOOS == "darwin" {
		plugins = append(plugins, table.NewPlugin("filevault_users", filevaultusers.FileVaultUsersColumns(), filevaultusers.FileVaultUsersGenerate))
		plugins = append(plugins, table.NewPlugin("macos_profiles", macosprofiles.MacOSProfilesColumns(), macosprofiles.MacOSProfilesGenerate))
		plugins = append(plugins, table.NewPlugin("mdm", mdm.MDMInfoColumns(), mdm.MDMInfoGenerate))
		plugins = append(plugins, table.NewPlugin("munki_info", munki.MunkiInfoColumns(), munki.MunkiInfoGenerate))
		plugins = append(plugins, table.NewPlugin("munki_installs", munki.MunkiInstallsColumns(), munki.MunkiInstallsGenerate))
		plugins = append(plugins, table.NewPlugin("network_quality", networkquality.NetworkQualityColumns(), networkquality.NetworkQualityGenerate))
		plugins = append(plugins, table.NewPlugin("pending_apple_updates", pendingappleupdates.PendingAppleUpdatesColumns(), pendingappleupdates.PendingAppleUpdatesGenerate))
		plugins = append(plugins, table.NewPlugin("macadmins_unified_log", unifiedlog.UnifiedLogColumns(), unifiedlog.UnifiedLogGenerate))
		plugins = append(plugins, table.NewPlugin("vscode_plugins", vscode_plugins.ExtensionsColumns(), vscode_plugins.ExtensionsGenerate))
		plugins = append(plugins, table.NewPlugin("asdf_installs", asdf_installs.AsdfInstallsColumns(), asdf_installs.AsdfInstallsGenerate))
		plugins = append(plugins, table.NewPlugin("intellij_plugins", intellij_plugins.IntelliJPluginsColumns(), intellij_plugins.IntelliJPluginsGenerate))

	}

	for _, p := range plugins {
		server.RegisterPlugin(p)
	}

	// Start the server. It will run forever unless an error bubbles up.
	if err := server.Run(); err != nil {
		log.Fatalln(err)
	}
}
