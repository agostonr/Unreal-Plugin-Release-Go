package executor

import (
	"fmt"
	"os"
	"os/exec"
)

/*
Uses Powershell for both calling the builder bat file and the compress to zip command.
*/
type WindowsExecutor struct {
}

/*
Creates the command that calls the RunUAT.bat file.
*/
func (e WindowsExecutor) CreateBuilderCommand(buildScriptPath string, pluginLocation string, outputDir string) *exec.Cmd {
	args := []string{
		"BuildPlugin",
		"-Plugin=" + pluginLocation,
		"-Package=" + outputDir,
		"-Rocket",
	}

	buildCmd := exec.Command("cmd", append([]string{"/C", buildScriptPath}, args...)...)
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	return buildCmd
}

/*
Creates the command that calls the zip command.
*/
func (e WindowsExecutor) CreateZipCommand(sourceDir string) *exec.Cmd {
	// Example PowerShell command:
	// Compress-Archive -Path "C:\MyFolder\*" -DestinationPath "C:\MyZip.zip" -Force
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf(`Compress-Archive -Path "%s\*" -DestinationPath "%s" -Force`, sourceDir, sourceDir+".zip"))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}
