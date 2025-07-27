/*
Copyright © 2025 AMSIAMUN
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type Config struct {
	EngineBaseDirectory string `json:"engineBaseDirectory"`
	BuildScriptPath     string `json:"buildScriptPath"`
	OutputBaseDirectory string `json:"outputBaseDirectory"`
	DocumentationPath   string `json:"documentationPath"`
}

const configFile = "config.json"
const configDirectoryName = "Config"
const pluginConfigurationIniFileName = "FilterPlugin.ini"

var engineVersions string
var pluginLocation string
var skipDocs bool

func init() {
	rootCmd.Flags().StringVar(&engineVersions, "engine-versions", "", "Comma-separated list of Unreal engine versions")
	rootCmd.Flags().StringVar(&pluginLocation, "plugin-location", "", "Full path to the .uplugin file")
	rootCmd.Flags().BoolVar(&skipDocs, "skip-docs", false, "(Optional) Skip adding documentation (FilterPlugin.ini and pdf doc location is needed otherwise)")
}

var rootCmd = &cobra.Command{
	Use:   "unreal-plugin-release",
	Short: "Build Unreal plugins across multiple engine versions.",
	Long: `Build Unreal plugins in batch to various Unreal Engine versions.

REQUIRES a config.json file next to the executable, which must contain:
  - engineBaseDirectory: the folder that contains the UE_5.1, UE_5.2 etc folders
  - buildScriptPath: the path to the RunUAT file within the engine dir
  - outputBaseDirectory: the path to the folder that will contain the built content
  - documentationPath: (optional) full path to a PDF file to include as plugin documentation

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

	config := loadConfigFile(configFile)
	versions := collectVersions(engineVersions)
	pluginName := createPluginName(pluginLocation)

	for _, version := range versions {
		version = strings.TrimSpace(version)
		if version == "" {
			continue
		}

		outputDir := combineOutputDir(pluginName, version, config)
		runBuildForEngineVersion(version, outputDir, config)
		postProcessRelease(outputDir, config)
	}

	fmt.Println("✅ All builds completed successfully.")
}

func collectVersions(input string) []string {
	return strings.Split(input, ",")
}

func createPluginName(pluginLocation string) string {
	return strings.TrimSuffix(filepath.Base(pluginLocation), filepath.Ext(pluginLocation))
}

func loadConfigFile(configFile string) *Config {
	file := openConfigFile(configFile)
	defer file.Close()
	return createConfig(file)
}

func combineOutputDir(pluginName string, version string, config *Config) string {
	return filepath.Join(config.OutputBaseDirectory, pluginName+"_"+version)
}

func postProcessRelease(releaseDir string, config *Config) {
	removeUnneededFolders(releaseDir)

	if !skipDocs {
		if err := handleDocumentation(releaseDir, config); err != nil {
			return
		}
	}

	zipPath := releaseDir + ".zip"
	if err := zipWithPowerShell(releaseDir, zipPath); err != nil {
		fmt.Println("⚠️ Failed to zip using PowerShell:", err)
	}
}

func zipWithPowerShell(sourceDir, zipPath string) error {
	// Example PowerShell command:
	// Compress-Archive -Path "C:\MyFolder\*" -DestinationPath "C:\MyZip.zip" -Force
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf(`Compress-Archive -Path "%s\*" -DestinationPath "%s" -Force`, sourceDir, zipPath))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func handleDocumentation(releaseDir string, config *Config) error {
	configErr := createConfigFolderWithIni(releaseDir)
	if configErr != nil {
		return configErr
	}

	if err := copyPdfIntoDocsFolderAndRename(releaseDir, config); err != nil {
		return err
	}

	return nil
}

func copyPdfIntoDocsFolderAndRename(releaseDir string, config *Config) error {
	// 1. Read FilterPlugin.ini next to the .exe
	filterPluginPath := filepath.Join(filepath.Dir(os.Args[0]), pluginConfigurationIniFileName)
	filterData, err := os.ReadFile(filterPluginPath)
	if err != nil {
		return fmt.Errorf("failed to read FilterPlugin.ini: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(filterData)), "\n")
	if len(lines) < 2 {
		return fmt.Errorf("FilterPlugin.ini must contain a second line for the doc path")
	}

	// 2. Get the relative path from the second line (e.g. /Docs/DocName.pdf)
	relativeDocPath := strings.TrimSpace(lines[1])
	relativeDocPath = strings.TrimPrefix(relativeDocPath, "/")
	targetDocPath := filepath.Join(releaseDir, relativeDocPath)

	// 3. Ensure target directories exist
	targetDir := filepath.Dir(targetDocPath)
	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create doc target folder: %w", err)
	}

	// 4. Copy the documentation PDF to the correct path
	srcFile, err := os.Open(config.DocumentationPath)
	if err != nil {
		return fmt.Errorf("failed to open documentation file: %w", err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(targetDocPath)
	if err != nil {
		return fmt.Errorf("failed to create destination doc file: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy doc file: %w", err)
	}

	return nil
}

func createConfigFolderWithIni(releaseDir string) error {
	configDir := filepath.Join(releaseDir, configDirectoryName)
	if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
		fmt.Println("⚠️ Failed to create Config dir:", err)
		return err
	}

	destIni := filepath.Join(configDir, pluginConfigurationIniFileName)
	if err := copyFile(pluginConfigurationIniFileName, destIni); err != nil {
		fmt.Println("⚠️ Failed to copy INI file:", err)
		return err
	}

	return nil
}

func copyFile(src, dest string) error {
	input, err := os.Open(src)
	if err != nil {
		return err
	}
	defer input.Close()

	output, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer output.Close()

	_, err = io.Copy(output, input)
	return err
}

func removeUnneededFolders(releaseDir string) {
	unneededFolders := []string{"Binaries", "Build", "Intermediate", "Saved"}
	for _, folder := range unneededFolders {
		fullPath := filepath.Join(releaseDir, folder)
		if err := os.RemoveAll(fullPath); err != nil {
			fmt.Println("⚠️ Failed to delete:", fullPath, "->", err)
		}
	}
}

func runBuildForEngineVersion(version string, outputDir string, config *Config) {
	batPath := filepath.Join(config.EngineBaseDirectory, "UE_"+version, config.BuildScriptPath)
	validateBatFilePath(batPath, version)

	cmdArgs := []string{
		"BuildPlugin",
		"-Plugin=" + pluginLocation,
		"-Package=" + outputDir,
		"-Rocket",
	}

	fmt.Println("======================================")
	fmt.Println("Building for UE version", version)
	fmt.Println("Command:", batPath, strings.Join(cmdArgs, " "))
	fmt.Println("Output to:", outputDir)
	fmt.Println("======================================")

	buildCmd := exec.Command("cmd", append([]string{"/C", batPath}, cmdArgs...)...)
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	if err := buildCmd.Run(); err != nil {
		fmt.Println("Build failed for", version, ":", err)
		fmt.Println("Deleting output directory:", config.OutputBaseDirectory)

		if !isDangerousPath(config.OutputBaseDirectory) {
			if removeErr := os.RemoveAll(config.OutputBaseDirectory); removeErr != nil {
				fmt.Println("Failed to remove output directory", removeErr)
			}
		} else {
			fmt.Println("Output directory is dangerous to delete.")
		}

		os.Exit(1)
	}
}

func isDangerousPath(path string) bool {
	lower := strings.ToLower(filepath.Clean(path))

	// Check for root drive (C:\, D:\, etc.)
	vol := filepath.VolumeName(lower)
	root := filepath.Join(vol, "\\")

	if lower == root {
		return true
	}

	// Check if it contains "windows" directory
	if strings.Contains(lower, `\windows`) {
		return true
	}

	return false
}

func validateBatFilePath(path string, version string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("Build script not found for engine version", version, ":", path)
		os.Exit(1)
	}
}

func createConfig(file *os.File) *Config {
	decoder := json.NewDecoder(file)
	config := Config{}
	if err := decoder.Decode(&config); err != nil {
		fmt.Println("Error parsing config file:", err)
		os.Exit(1)
	}

	return &config
}

func openConfigFile(path string) *os.File {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening config file:", err)
		os.Exit(1)
	}
	return f
}

func isInputValid() bool {
	if engineVersions == "" || pluginLocation == "" {
		fmt.Println("Missing required flags: --engine-versions and --plugin-location are both required.")
		return false
	}

	if _, err := os.Stat(pluginLocation); os.IsNotExist(err) {
		fmt.Println("Plugin file not found:", pluginLocation)
		return false
	}

	return true
}
