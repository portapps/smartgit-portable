//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=res/papp.ico -manifest=res/papp.manifest
package main

import (
	"os"
	"strings"

	. "github.com/portapps/portapps"
	"github.com/portapps/portapps/pkg/utl"
)

var (
	app *App
)

func init() {
	var err error

	// Init app
	if app, err = New("smartgit-portable", "SmartGit"); err != nil {
		Log.Fatal().Err(err).Msg("Cannot initialize application. See log file for more info.")
	}
}

func main() {
	utl.CreateFolder(app.DataPath)
	app.Process = utl.PathJoin(app.AppPath, "bin", "smartgit.exe")
	app.WorkingDir = utl.PathJoin(app.AppPath, "bin")

	// create err folder
	utl.CreateFolder(app.DataPath, "err")

	// create default smartgit.vmoptions if not found
	customSmartgitOptionsPath := utl.PathJoin(app.DataPath, "smartgit.vmoptions")
	if !utl.Exists(customSmartgitOptionsPath) {
		if err := utl.CreateFile(customSmartgitOptionsPath, `-Xmx1024m
-Dsmartgit.disableBugReporting=true
`); err != nil {
			Log.Fatal().Err(err).Msg("Cannot write default smartgit.vmoptions")
		}
	}

	// override system smartgit.vmoptions
	smartgitOptionsPath := utl.PathJoin(app.AppPath, "bin", "smartgit.vmoptions")
	if err := utl.CreateFile(smartgitOptionsPath, strings.Replace(`-Dsmartboot.sourceDirectory={{ DATA_PATH }}\.updates
-Dsmartgit.settings={{ DATA_PATH }}\.settings
-Dsmartgit.updateCheck.enabled=false
-Dsmartgit.updateCheck.automatic=false
-Dsmartgit.updateCheck.checkForLatestBuildVisible=false
-Dsmartgit.disableBugReporting=true
-XX:ErrorFile={{ DATA_PATH }}\err\hs_err_pid%p.log
-include-options {{ DATA_PATH }}\smartgit.vmoptions
`, "{{ DATA_PATH }}", utl.FormatWindowsPath(app.DataPath), -1)); err != nil {
		Log.Fatal().Err(err).Msg("Cannot write system smartgit.vmoptions")
	}

	// set JAVA_HOME
	utl.OverrideEnv("SMARTGIT_JAVA_HOME", utl.PathJoin(app.AppPath, "jre"))

	app.Launch(os.Args[1:])
}
