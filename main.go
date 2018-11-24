//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=res/papp.ico
package main

import (
	"os"
	"runtime"
	"strings"

	. "github.com/portapps/portapps"
)

func init() {
	Papp.ID = "smartgit-portable"
	Papp.Name = "SmartGit"
	Init()
}

func main() {
	smartgitExe := "smartgit32.exe"
	if runtime.GOARCH == "amd64" {
		smartgitExe = "smartgit.exe"
	}

	Papp.AppPath = AppPathJoin("app")
	Papp.DataPath = CreateFolder(AppPathJoin("data"))
	Papp.Process = PathJoin(Papp.AppPath, "bin", smartgitExe)
	Papp.Args = nil
	Papp.WorkingDir = PathJoin(Papp.AppPath, "bin")

	CreateFolder(PathJoin(Papp.DataPath, "err"))

	// create default smartgit.vmoptions if not found
	customSmartgitOptionsPath := PathJoin(Papp.DataPath, "smartgit.vmoptions")
	if !Exists(customSmartgitOptionsPath) {
		if err := CreateFile(customSmartgitOptionsPath, `-Xmx1024m
-Dsmartgit.disableBugReporting=true
`); err != nil {
			Log.Errorf("Cannot write default smartgit.vmoptions: %s", err)
		}
	}

	// override system smartgit.vmoptions
	smartgitOptionsPath := PathJoin(Papp.AppPath, "bin", "smartgit.vmoptions")
	if err := CreateFile(smartgitOptionsPath, strings.Replace(`-Dsmartboot.sourceDirectory={{ DATA_PATH }}\.updates
-Dsmartgit.settings={{ DATA_PATH }}\.settings
-Dsmartgit.updateCheck.enabled=false
-Dsmartgit.updateCheck.automatic=false
-Dsmartgit.updateCheck.checkForLatestBuildVisible=false
-Dsmartgit.disableBugReporting=true
-XX:ErrorFile={{ DATA_PATH }}\err\hs_err_pid%p.log
-include-options {{ DATA_PATH }}\smartgit.vmoptions
`, "{{ DATA_PATH }}", FormatWindowsPath(Papp.DataPath), -1)); err != nil {
		Log.Errorf("Cannot write system smartgit.vmoptions: %s", err)
	}

	Launch(os.Args[1:])
}
