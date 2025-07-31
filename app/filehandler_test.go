package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"unreal-plugin-release/model"
)

type isDangerousPathTestData struct {
	path     string
	expected bool
}

func TestCreateValidConfigShouldReturnConfigWithNoError(t *testing.T) {
	// given
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")
	engineDir := filepath.Join(tempDir, "Engine")
	if err := os.MkdirAll(engineDir, 0755); err != nil {
		t.Fatal("Failed to create engine dir")
	}
	outputDir := filepath.Join(tempDir, "Output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatal("Failed to create output dir")
	}

	// create expected config
	expected := model.Config{
		EngineBaseDirectory: engineDir,
		BuildScriptPath:     "RunUAT.bat",
		OutputBaseDirectory: outputDir,
	}
	data, marshallErr := json.MarshalIndent(expected, "", " ")
	if marshallErr != nil {
		t.Fatal("Failed to marshall config file")
	}
	if err := os.WriteFile(configPath, data, 0755); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// when
	config, err := CreateConfig(configPath)

	// then
	if err != nil || config == nil {
		t.Error("Configuration creation failed")
	}
}

func TestCreateMissingConfigShouldFail(t *testing.T) {
	// given
	tempDir := t.TempDir()

	// when
	config, err := CreateConfig(tempDir)

	// then
	if config != nil || err == nil {
		t.Error("A non existing config file should return nil and an error.")
	}
}

func TestGetFullPathForFileInExecDirShouldEndWithTheFileName(t *testing.T) {
	// given
	filename := "FilterPlugin.ini"
	base := t.TempDir()
	exePath := filepath.Join(base, "script.exe")
	expected := filepath.Join(base, "FilterPlugin.ini")

	// when
	actual := GetFullPathForFileInExecDir(exePath, filename)

	// then
	if actual != expected {
		t.Errorf("Expected path to end with %s, got %s", filename, actual)
	}
}

func TestIsPathExistWithExistingPath(t *testing.T) {
	// given
	tempDir := t.TempDir()

	// when
	actual := IsPathExist(tempDir)

	// then
	if !actual {
		t.Error("An existing folder should not be deemed non existent.")
	}
}

func TestIsPathExistWithNonExistingPath(t *testing.T) {
	// given
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "NonExistentDir")

	// when
	actual := IsPathExist(path)

	// then
	if actual {
		t.Error("A non existent folder should not be deemed existing.")
	}
}

func TestIsPathEqualShouldReturnTrueWithSamePath(t *testing.T) {
	// given
	tempDir := t.TempDir()

	// when
	actual := IsPathEqual(tempDir, tempDir)

	// then
	if !actual {
		t.Error("The same path should be deemed equal to itself.")
	}
}

func TestIsPathEqualShouldReturnFalseWithDifferentPaths(t *testing.T) {
	// given
	tempDir := t.TempDir()
	p1 := filepath.Join(tempDir, "One")
	p2 := filepath.Join(tempDir, "Two")

	// when
	actual := IsPathEqual(p1, p2)

	// then
	if actual {
		t.Error("Two different filepaths should be deemed unequal.")
	}
}

func TestIsFileShouldReturnTrueForAFile(t *testing.T) {
	// given
	tempDir := t.TempDir()
	file := filepath.Join(tempDir, "file.txt")
	f, err := os.Create(file)
	if err != nil {
		t.Fatal("Failed to create file in temp dir")
	}
	f.Close()

	// when
	actual := IsFile(file)

	// then
	if !actual {
		t.Error("A file should not be deemed a directory.")
	}
}

func TestIsFileShouldReturnFalseForADirectory(t *testing.T) {
	// given
	tempDir := t.TempDir()

	// when
	actual := IsFile(tempDir)

	// then
	if actual {
		t.Error("A directory should not be considered a file.")
	}
}

func TestCreatePluginNameShouldReturnThePluginNameWithoutExtension(t *testing.T) {
	// given
	tempDir := t.TempDir()
	name := "MyPlugin"
	plugin := filepath.Join(tempDir, name+".upligin")

	// when
	actual := createPluginName(plugin)

	// then
	if name != actual {
		t.Errorf("Expected: %q, Actual: %q", name, actual)
	}
}

func TestCombineOutputDirShouldPutPluginNameAndVersionTogether(t *testing.T) {
	// given
	tempDir := t.TempDir()
	expected := tempDir + "\\MyPlugin_5.5"
	fmt.Printf("Expected directory: %q", expected)

	// when
	actual := combineOutputDir("MyPlugin", "5.5", tempDir)

	// then
	if expected != actual {
		t.Errorf("Expected: %q, Actual: %q", expected, actual)
	}
}

func TestCreateBatFilePathShouldCombinePathWithEngineBaseDir(t *testing.T) {
	// given
	engineBase := t.TempDir()
	buildScriptPath := "Tools\\RunUAT.bat"
	version := "5.5"

	// when
	actual := createBatFilePath(engineBase, version, buildScriptPath)

	// then
	if engineBase+"\\UE_5.5\\Tools\\RunUAT.bat" != actual {
		t.Error("Path combination to bat file was incorrect.")
	}
}

func TestExistingFilePathShouldBeValid(t *testing.T) {
	// given
	tempDir := t.TempDir()
	file := filepath.Join(tempDir, "file.txt")
	f, err := os.Create(file)
	if err != nil {
		t.Fatal("Failed to create file in temp dir")
	}
	f.Close()

	// when
	actual := isFilePathValid(file)

	// then
	if !actual {
		t.Error("An existing file's path should be marked as valid.")
	}
}

func TestNonExistingFilePathShouldBeInvalid(t *testing.T) {
	// given
	tempDir := t.TempDir()
	file := filepath.Join(tempDir, "file.txt")

	// when
	actual := isFilePathValid(file)

	// then
	if actual {
		t.Error("A non existing file's path should be marked as invalid.")
	}
}

func TestRemoveDirectoryShouldRemoveNonDangerousPath(t *testing.T) {
	// given
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "Directory")
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatal("Failed to create dir")
	}

	// when
	removeDirectory(path)

	// then
	if _, err := os.Stat(path); err == nil {
		t.Error("The directory should have been removed.")
	}
}

func TestRemoveDirectoryShouldNotRemoveDangerousPath(t *testing.T) {
	// given
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "windows\\directory")
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatal("Failed to create dir")
	}

	// when
	removeDirectory(path)

	// then
	if _, err := os.Stat(path); err != nil {
		t.Error("The directory should not have been removed.")
	}
}

func TestIsDangerousPath(t *testing.T) {
	for i, tt := range createIsDangerousPathTestData() {
		t.Run("EngineVersionTest #"+strconv.Itoa(i), func(t *testing.T) {
			// when
			actual := isDangerousPath(tt.path)

			// then
			if tt.expected != actual {
				t.Errorf("%d case: expected: %v, actual: %v", i, tt.expected, actual)
			}
		})
	}
}

func TestDeleteSubFolder(t *testing.T) {
	// given
	base := t.TempDir()
	subfolderName := "Subfolder"
	path := filepath.Join(base, subfolderName)
	if err := os.Mkdir(path, 0755); err != nil {
		t.Fatal("Failed to create dir")
	}

	// when
	deleteSubFolder(base, subfolderName)

	// then
	_, subfolderErr := os.Stat(path)
	baseInfo, baseErr := os.Stat(base)

	if subfolderErr == nil {
		t.Error("The subfolder should have been deleted.")
	}

	if baseErr != nil || !baseInfo.IsDir() {
		t.Error("The directory should have been left intact.")
	}
}

func TestCreateConfigFolderWithIni(t *testing.T) {
	// given
	releaseDir := t.TempDir()
	filePath := filepath.Join(t.TempDir(), "FilterPlugin.ini")
	content := []byte("testing")
	os.WriteFile(filePath, content, 0755)

	// when
	createConfigFolderWithIni(releaseDir, filePath)

	// then
	copiedData, err := os.ReadFile(filepath.Join(releaseDir, "Config", "FilterPlugin.ini"))
	if err != nil {
		t.Error("The file should have been copied to the release directory's appropriate folder.")
	}

	if string(copiedData) != "testing" {
		t.Error("The contents of the copied file does not match the contents of the original one.")
	}
}

func TestCopyPdfIntoDocsFolderAndRename(t *testing.T) {
	// given
	releaseDir := t.TempDir()
	testDataFolder := "testdata"
	docsPath := filepath.Join(testDataFolder, "testing.pdf")
	filterPluginPath := filepath.Join(testDataFolder, "FilterPluginTest.ini")

	// when
	copyPdfIntoDocsFolderAndRename(releaseDir, docsPath, filterPluginPath)

	// then
	data, err := os.ReadFile(filepath.Join(releaseDir, "Docs", "My_Docs.pdf"))
	if err != nil {
		t.Error("Moving documentation or renaming it was not correct.")
	}

	if len(data) == 0 {
		t.Error("The contents of the copied docs file are empty.")
	}
}

// test data
func createIsDangerousPathTestData() []isDangerousPathTestData {
	return []isDangerousPathTestData{
		{
			"C:\\",
			true,
		},
		{
			"D:\\",
			true,
		},
		{
			"E:\\",
			true,
		},
		{
			"C:\\Windows\\AnyDirectory",
			true,
		},
		{
			"D:\\Windows\\AnyDirectory",
			true,
		},
		{
			"D:\\Games\\MyPlugin\\Release",
			false,
		},
		{
			"C:\\Users\\JohnDoe\\Documents\\Plugins\\MyPlugin",
			false,
		},
	}
}
