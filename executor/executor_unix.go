package executor

import (
	"os/exec"
)

type UnixExecutor struct {
}

/*
Implement this method to call the unix version of Epic's RunUAT.bat, which builds the plugin for the release.
To see the logs of the build script, don't forget to route os logs to that of this executable, e.g.

	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
*/
func (e UnixExecutor) CreateBuilderCommand(buildScriptPath string, pluginLocation string, outputDir string) *exec.Cmd {
	panic("Unix executor for build script command is not implemented.")
}

/*
Implement this method to zip the created plugin folder.
To see the logs of the zip command, don't forget to route os logs to that of this executable, e.g.

	zipCmd.Stdout = os.Stdout
	zipCmd.Stderr = os.Stderr
*/
func (e UnixExecutor) CreateZipCommand(sourceDir string) *exec.Cmd {
	panic("Unix executor for zip command is not implemented.")
}
