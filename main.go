//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/portapps/portapps/v3"
	"github.com/portapps/portapps/v3/pkg/log"
)

var (
	app *portapps.App
)

func init() {
	var err error

	// Init app
	if app, err = portapps.New("smartgit-portable", "SmartGit"); err != nil {
		log.Fatal().Err(err).Msg("Cannot initialize application. See log file for more info.")
	}
}

func main() {
	if err := os.MkdirAll(app.DataPath, os.ModePerm); err != nil {
		log.Fatal().Err(err).Msg("Cannot create data directory.")
	}
	app.Process = filepath.Join(app.AppPath, "bin", "smartgit.exe")
	app.WorkingDir = filepath.Join(app.AppPath, "bin")

	// create err folder
	if err := os.MkdirAll(filepath.Join(app.DataPath, "err"), os.ModePerm); err != nil {
		log.Fatal().Err(err).Msg("Cannot create error directory.")
	}

	// create default smartgit.vmoptions if not found
	customSmartgitOptionsPath := filepath.Join(app.DataPath, "smartgit.vmoptions")
	if _, err := os.Stat(customSmartgitOptionsPath); os.IsNotExist(err) {
		if err := os.WriteFile(customSmartgitOptionsPath, []byte(`-Xmx1024m
-Dsmartgit.disableBugReporting=true
`), 0644); err != nil {
			log.Fatal().Err(err).Msg("Cannot write default smartgit.vmoptions")
		}
	}

	// override system smartgit.vmoptions
	smartgitOptionsPath := filepath.Join(app.AppPath, "bin", "smartgit.vmoptions")
	if err := os.WriteFile(smartgitOptionsPath, []byte(strings.Replace(`-Dsmartboot.sourceDirectory={{ DATA_PATH }}\.updates
-Dsmartgit.settings={{ DATA_PATH }}\.settings
-Dsmartgit.updateCheck.enabled=false
-Dsmartgit.updateCheck.automatic=false
-XX:ErrorFile={{ DATA_PATH }}\err\hs_err_pid%p.log
-include-options {{ DATA_PATH }}\smartgit.vmoptions
`, "{{ DATA_PATH }}", filepath.FromSlash(app.DataPath), -1)), 0644); err != nil {
		log.Fatal().Err(err).Msg("Cannot write system smartgit.vmoptions")
	}

	// set JAVA_HOME
	os.Setenv("SMARTGIT_JAVA_HOME", filepath.Join(app.AppPath, "jre"))

	defer app.Close()
	app.Launch(os.Args[1:])
}
