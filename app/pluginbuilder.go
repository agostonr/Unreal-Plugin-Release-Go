// contains the business logic of the application
package app

import (
	"fmt"
	"os"
	"strings"

	"unreal-plugin-release/executor"
	"unreal-plugin-release/model"
)

// the application that encapsulates the core business logic with a configuration as input
type PluginBuilder struct {
	config *model.Config
	runner executor.SubprocessExecutor
}

/*
Constructor for the plugin builder.
*/
func NewPluginBuilder(config *model.Config, runner executor.SubprocessExecutor) *PluginBuilder {
	return &PluginBuilder{config, runner}
}

/*
Builds the plugins for all selected versions.
*/
func (pb *PluginBuilder) BuildPluginsForSelectedVersions(cmdInput model.CmdInput, execPath string) {
	versions := pb.collectVersions(cmdInput.EngineVersions)
	pluginName := createPluginName(pb.config.PluginPath)

	for _, version := range versions {
		version = strings.TrimSpace(version)
		if version == "" {
			continue
		}

		outputDir := combineOutputDir(pluginName, version, pb.config.OutputBaseDirectory)
		pb.runBuildForEngineVersion(version, outputDir, pb.config.PluginPath)
		pb.postProcessRelease(outputDir, execPath, pb.config.DocsPath, cmdInput)
	}
}

func (pb *PluginBuilder) postProcessRelease(outputDir, execPath, docsPath string, cmdInput model.CmdInput) {
	pb.removeUnneededFolders(outputDir)

	if docsPath != "" && !cmdInput.SkipDocs {
		if err := pb.handleDocumentation(outputDir, execPath, docsPath); err != nil {
			return
		}
	}

	if err := pb.runner.CreateZipCommand(outputDir).Run(); err != nil {
		fmt.Println("⚠️ Failed to zip using PowerShell:", err)
	}
}

func (pb *PluginBuilder) handleDocumentation(releaseDir, execPath, docsPath string) error {
	sourceIni := GetFullPathForFileInExecDir(execPath, model.PluginConfigurationIniFileName)
	if err := createConfigFolderWithIni(releaseDir, sourceIni); err != nil {
		return err
	}

	if err := copyPdfIntoDocsFolderAndRename(releaseDir, docsPath, GetFullPathForFileInExecDir(execPath, model.PluginConfigurationIniFileName)); err != nil {
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

func (pb *PluginBuilder) runBuildForEngineVersion(version, outputDir, pluginPath string) {
	buildScriptPath := pb.makeBuildScriptFilePath(version)

	fmt.Println("======================================")
	fmt.Println("Building for UE version", version)
	fmt.Println("Output to:", outputDir)
	fmt.Println("======================================")

	err := pb.runner.CreateBuilderCommand(buildScriptPath, pluginPath, outputDir).Run()
	if err != nil {
		fmt.Println("Build failed for", version, ":", err)
		removeDirectory(pb.config.OutputBaseDirectory)
		os.Exit(1)
	}
}

func (pb *PluginBuilder) collectVersions(input string) []string {
	return strings.Split(input, ",")
}

func (pb *PluginBuilder) makeBuildScriptFilePath(version string) string {
	batPath := createBatFilePath(pb.config.EngineBaseDirectory, version, pb.config.BuildScriptPath)
	if !isFilePathValid(batPath) {
		fmt.Println("Build script not found for engine version", version, ":", batPath)
		os.Exit(1)
	}
	return batPath
}
