// contains the business logic of the application
package app

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"unreal-plugin-release/model"
)

// the application that encapsulates the core business logic with a configuration as input
type PluginBuilder struct {
	config *model.Config
}

/*
Constructor for the plugin builder.
*/
func NewPluginBuilder(config *model.Config) *PluginBuilder {
	return &PluginBuilder{config}
}

/*
Builds the plugins for all selected versions.
*/
func (pb *PluginBuilder) BuildPluginsForSelectedVersions(cmdInput model.CmdInput) {
	versions := pb.collectVersions(cmdInput.EngineVersions)
	pluginName := createPluginName(cmdInput.PluginLocation)

	for _, version := range versions {
		version = strings.TrimSpace(version)
		if version == "" {
			continue
		}

		outputDir := combineOutputDir(pluginName, version, pb.config.OutputBaseDirectory)
		pb.runBuildForEngineVersion(version, outputDir, cmdInput)
		pb.postProcessRelease(outputDir, cmdInput.DocsPath)
	}
}

func (pb *PluginBuilder) postProcessRelease(outputDir string, docsPath string) {
	pb.removeUnneededFolders(outputDir)

	if docsPath != "" {
		if err := pb.handleDocumentation(outputDir, docsPath); err != nil {
			return
		}
	}

	if err := pb.zipWithPowerShell(outputDir); err != nil {
		fmt.Println("⚠️ Failed to zip using PowerShell:", err)
	}
}

func (pb *PluginBuilder) zipWithPowerShell(sourceDir string) error {
	// Example PowerShell command:
	// Compress-Archive -Path "C:\MyFolder\*" -DestinationPath "C:\MyZip.zip" -Force
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf(`Compress-Archive -Path "%s\*" -DestinationPath "%s" -Force`, sourceDir, sourceDir+".zip"))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (pb *PluginBuilder) handleDocumentation(releaseDir string, docsPath string) error {
	if err := createConfigFolderWithIni(releaseDir); err != nil {
		return err
	}

	if err := copyPdfIntoDocsFolderAndRename(releaseDir, docsPath); err != nil {
		return err
	}

	return nil
}

func (pb *PluginBuilder) removeUnneededFolders(releaseDir string) {
	for _, folder := range pb.getUnneededFolders() {
		deleteSubFolder(releaseDir, folder)
	}
}

func (pb *PluginBuilder) getUnneededFolders() []string {
	return []string{"Binaries", "Build", "Intermediate", "Saved"}
}

func (pb *PluginBuilder) runBuildForEngineVersion(version string, outputDir string, cmdInput model.CmdInput) {
	batPath := pb.makeBatFilePath(version)
	cmdArgs := pb.createBuilderCommandArgs(cmdInput.PluginLocation, outputDir)

	fmt.Println("======================================")
	fmt.Println("Building for UE version", version)
	fmt.Println("Output to:", outputDir)
	fmt.Println("======================================")

	err := pb.createBuilderCommand(batPath, cmdArgs).Run()
	if err != nil {
		fmt.Println("Build failed for", version, ":", err)
		removeDirectory(pb.config.OutputBaseDirectory)
		os.Exit(1)
	}
}

func (pb *PluginBuilder) createBuilderCommand(batPath string, args []string) *exec.Cmd {
	buildCmd := exec.Command("cmd", append([]string{"/C", batPath}, args...)...)
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	return buildCmd
}

func (pb *PluginBuilder) createBuilderCommandArgs(pluginLocation string, outputDir string) []string {
	return []string{
		"BuildPlugin",
		"-Plugin=" + pluginLocation,
		"-Package=" + outputDir,
		"-Rocket",
	}
}

func (pb *PluginBuilder) collectVersions(input string) []string {
	return strings.Split(input, ",")
}

func (pb *PluginBuilder) makeBatFilePath(version string) string {
	batPath := createBatFilePath(pb.config.EngineBaseDirectory, version, pb.config.BuildScriptPath)
	if err := validateFilePath(batPath); err != nil {
		fmt.Println("Build script not found for engine version", version, ":", batPath)
		os.Exit(1)
	}
	return batPath
}
