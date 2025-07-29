/*
Copyright © 2025 AMSIAMUN
*/
package cmd

import (
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"

	"unreal-plugin-release/app"
	"unreal-plugin-release/model"
)

var cmdInput = model.CmdInput{}

func init() {
	rootCmd.Flags().StringVar(&cmdInput.EngineVersions, "engine-versions", "", "Comma-separated list of Unreal engine versions")
	rootCmd.Flags().StringVar(&cmdInput.PluginLocation, "plugin-location", "", "Full path to the .uplugin file")
	rootCmd.Flags().StringVar(&cmdInput.DocsPath, "docs-path", "", "(Optional) The full path to the documentation pdf, if there is one for the plugin. If present, it will be copied to the path in the plugin where the FilterPlugin.ini file directs.")
}

var rootCmd = &cobra.Command{
	Use:   "unreal-plugin-release",
	Short: "Build Unreal plugins across multiple engine versions.",
	Long: `Build Unreal plugins in batch to various Unreal Engine versions.

REQUIRES a config.json file next to the executable, which must contain:
  - engineBaseDirectory: the folder that contains the UE_5.1, UE_5.2 etc folders
  - buildScriptPath: the path to the RunUAT file within the engine dir
  - outputBaseDirectory: the path to the folder that will contain the built content

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
	if !isInputValid() {
		os.Exit(1)
	}

	config := app.CreateConfig(app.GetFullPathForFileInExecDir(model.ConfigFile))
	app.NewPluginBuilder(config).BuildPluginsForSelectedVersions(cmdInput)

	fmt.Println("✅ All builds completed successfully.")
}

func isInputValid() bool {
	return isEngineVersionsValid() && isPluginLocationValid()
}

func isPluginLocationValid() bool {
	if cmdInput.PluginLocation == "" {
		fmt.Println("--plugin-location is mandatory.")
		return false
	}

	if !app.IsPathExist(cmdInput.PluginLocation) {
		fmt.Println("Plugin file not found:", cmdInput.PluginLocation)
		return false
	}

	return true
}

func isEngineVersionsValid() bool {
	if cmdInput.EngineVersions == "" {
		fmt.Println("Missing required flags: --engine-versions and --plugin-location are both required.")
		return false
	}

	validatorExpression := regexp.MustCompile(`^\d+\.\d+(,\d+\.\d+)*$`)
	if !validatorExpression.MatchString(cmdInput.EngineVersions) {
		fmt.Println("Spelling error in unreal engine versions. Must be MAJOR.MINOR e.g. 5.6, separated by commas.")
		return false
	}

	return true
}
