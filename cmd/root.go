/*
Copyright © 2025 AMSIAMUN
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"

	"unreal-plugin-release/app"
	"unreal-plugin-release/executor"
	"unreal-plugin-release/model"
)

var cmdInput = model.CmdInput{}

func init() {
	rootCmd.Flags().StringVar(&cmdInput.EngineVersions, "engine-versions", "", "Comma-separated list of Unreal engine versions")
	rootCmd.Flags().BoolVar(&cmdInput.SkipDocs, "skip-docs", false, "Omit copying documentation")
}

var rootCmd = &cobra.Command{
	Use:   "unreal-plugin-release",
	Short: "Build Unreal plugins across multiple engine versions.",
	Long: `Build Unreal plugins in batch to various Unreal Engine versions.

REQUIRES a config.json file next to the executable, which must contain:
  - engineBaseDirectory: the folder that contains the UE_5.1, UE_5.2 etc folders
  - buildScriptPath: the path to the RunUAT file within the engine dir
  - outputBaseDirectory: the path to the folder that will contain the built content
  - pluginPath: the path to the .uplugin file to be built
  - docsPath: (optional) the path to the pdf documentation

If documentation is enabled, a FilterPlugin.ini file must also exist next to the executable.
It should contain the expected internal documentation path like so:

  [FilterPlugin]
  /Documentation/My_Documentation.pdf`,
	Run: runRootCommand,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runRootCommand(cmd *cobra.Command, args []string) {
	if !isEngineVersionsValid() {
		os.Exit(1)
	}

	execPath, err := os.Executable()
	if err != nil {
		panic("Executable file not found")
	}

	config, err := createAndValidateConfig(app.GetFullPathForFileInExecDir(execPath, model.ConfigFile))
	if err != nil {
		os.Exit(1)
	}

	app.NewPluginBuilder(config, executor.NewExecutor()).BuildPluginsForSelectedVersions(cmdInput, execPath)
	fmt.Println("✅ All builds completed successfully.")
}

func isEngineVersionsValid() bool {
	if cmdInput.EngineVersions == "" {
		fmt.Println("Missing required flag: --engine-versions is required.")
		return false
	}

	validatorExpression := regexp.MustCompile(`^\d+\.\d+(,\d+\.\d+)*$`)
	if !validatorExpression.MatchString(cmdInput.EngineVersions) {
		fmt.Println("Spelling error in unreal engine versions. Must be MAJOR.MINOR e.g. 5.6, separated by commas.")
		return false
	}

	return true
}

func createAndValidateConfig(configPath string) (*model.Config, error) {
	config, err := app.CreateConfig(configPath)
	if err != nil || config == nil {
		fmt.Print("Failed to load config file")
		return nil, err
	}

	if !isConfigValid(config) {
		fmt.Printf("The config file contains invalid path.")
		return nil, errors.New("invalid config")
	}

	return config, nil
}

func isConfigValid(config *model.Config) bool {
	// the bat file path is a relative path, so it cannot be validated here, only after the full path per version is assembled.
	return app.IsPathExist(config.EngineBaseDirectory) &&
		app.IsPathExist(config.OutputBaseDirectory) &&
		!app.IsPathEqual(config.EngineBaseDirectory, config.OutputBaseDirectory) &&
		isPluginLocationValid(config.PluginPath) &&
		config.BuildScriptPath != ""
}

func isPluginLocationValid(pluginPath string) bool {
	if pluginPath == "" {
		fmt.Println("--plugin-location is mandatory.")
		return false
	}

	if !app.IsPathExist(pluginPath) {
		fmt.Println("Plugin file not found:", pluginPath)
		return false
	}

	if !app.IsFile(pluginPath) {
		fmt.Println("Plugin location needs to point to a file, not a directory.")
		return false
	}

	return true
}
