// responsible for file operations like opening, reading, copying or deleting files, or building paths.
package app

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"unreal-plugin-release/model"
)

/*
Create the configuration dto from the config file, by location.
*/
func CreateConfig(path string) (*model.Config, error) {
	file := openConfigFile(path)
	decoder := json.NewDecoder(file)
	defer file.Close()
	config := model.Config{}
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

/*
Get the full path for the given file that's next to the executable.
*/
func GetFullPathForFileInExecDir(exePath, filename string) string {
	return filepath.Join(filepath.Dir(exePath), filename)
}

/*
Determines if the path exists and is valid.
*/
func IsPathExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

/*
Determines if two paths are the same.
*/
func IsPathEqual(p1 string, p2 string) bool {
	abs1, err1 := filepath.Abs(filepath.Clean(p1))
	abs2, err2 := filepath.Abs(filepath.Clean(p2))

	if err1 != nil || err2 != nil {
		return false
	}

	return abs1 == abs2
}

/*
Determines if a path is a file.
*/
func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func openConfigFile(path string) *os.File {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening config file:", err)
		os.Exit(1)
	}
	return f
}

func createPluginName(pluginLocation string) string {
	return strings.TrimSuffix(filepath.Base(pluginLocation), filepath.Ext(pluginLocation))
}

func combineOutputDir(pluginName string, version string, outputDir string) string {
	return filepath.Join(outputDir, pluginName+"_"+version)
}

func createBatFilePath(engineBaseDir string, version string, buildScriptPath string) string {
	return filepath.Join(engineBaseDir, "UE_"+version, buildScriptPath)
}

func isFilePathValid(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func removeDirectory(path string) {
	if !isDangerousPath(path) {
		if removeErr := os.RemoveAll(path); removeErr != nil {
			fmt.Println("Failed to remove output directory", removeErr)
		}
	} else {
		fmt.Println("Output directory is dangerous to delete.")
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

func deleteSubFolder(baseDir string, dirToDelete string) {
	fullPath := filepath.Join(baseDir, dirToDelete)
	if err := os.RemoveAll(fullPath); err != nil {
		fmt.Println("⚠️ Failed to delete:", fullPath, "->", err)
	}
}

func createConfigFolderWithIni(releaseDir string, sourceIniPath string) error {
	configDir := filepath.Join(releaseDir, model.ConfigDirectoryName)
	if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
		fmt.Println("⚠️ Failed to create Config dir:", err)
		return err
	}

	destIni := filepath.Join(configDir, model.PluginConfigurationIniFileName)
	if err := copyFile(sourceIniPath, destIni); err != nil {
		fmt.Println("⚠️ Failed to copy INI file:", err)
		return err
	}

	return nil
}

func copyPdfIntoDocsFolderAndRename(releaseDir string, docsPath string, filterPluginFilePath string) error {
	// 1. Read FilterPlugin.ini next to the .exe
	filterData, err := os.ReadFile(filterPluginFilePath)
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
	srcFile, err := os.Open(docsPath)
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
