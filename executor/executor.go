/*
This file contains the definition of running any executable from within the application, and the
Windows implementation for it.

Developers using this, who use Mac or Linux need to make their own implementations for the same tasks,
based on the required operation detailed.
*/
package executor

import (
	"fmt"
	"os/exec"
	"runtime"
)

/*
Creates the commands that are ran as subprocesses in the application.
Different platforms like Mac or Linux need their own implementation.
*/
type SubprocessExecutor interface {
	CreateZipCommand(sourceDir string) *exec.Cmd
	CreateBuilderCommand(buildScriptPath string, pluginLocation string, outputDir string) *exec.Cmd
}

/*
Automatically pick which OS we are on, and return the appropriate executor for it.
*/
func NewExecutor() SubprocessExecutor {
	switch runtime.GOOS {
	case "windows":
		return WindowsExecutor{}
	case "darwin", "linux":
		return UnixExecutor{}
	default:
		panic(fmt.Sprintf("unsupported OS: %s", runtime.GOOS))
	}
}
