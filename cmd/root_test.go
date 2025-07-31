package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"unreal-plugin-release/model"
)

type engineVersionTestData struct {
	engineVersions string
	expected       bool
}

type validInputTestData struct {
	engineVersions string
	expected       bool
}

func TestIsInputValid(t *testing.T) {
	for i, tt := range createIsInputValidTestData() {
		t.Run("IsInputValid #"+strconv.Itoa(i), func(t *testing.T) {
			// given
			cmdInput.EngineVersions = tt.engineVersions

			// when
			actual := isEngineVersionsValid()

			// then
			if actual != tt.expected {
				t.Errorf("Case %d: Expected isValid: %t, actual isValid: %t", i, tt.expected, actual)
			}
		})
	}
}

func TestValidPluginLocation(t *testing.T) {
	// when
	actual := isPluginLocationValid(createPluginAtTempDir("MyPlugin.uplugin", t))

	// then
	if !actual {
		t.Error("A valid plugin location must not be marked invalid.")
	}
}

func TestInvalidPluginLocation(t *testing.T) {
	// given
	base := t.TempDir()
	missingPath := filepath.Join(base, "MissingDirectory", "NoPlugin.uplugin")

	// when
	actual := isPluginLocationValid(missingPath)

	// then
	if actual {
		t.Error("A missing plugin location must be invalid.")
	}
}

func TestMissingPluginLocation(t *testing.T) {
	// when
	actual := isPluginLocationValid("")

	// then
	if actual {
		t.Error("An empty plugin location must be invalid.")
	}
}

func TestNotFilePluginLocation(t *testing.T) {
	// given
	dir := t.TempDir()

	// when
	actual := isPluginLocationValid(dir)

	// then
	if actual {
		t.Error("A directory as plugin location must be invalid.")
	}
}

func TestIsEngineVersionsValid(t *testing.T) {
	for i, tt := range createEngineVersionsTestData() {
		t.Run("EngineVersionTest #"+strconv.Itoa(i), func(t *testing.T) {
			// given
			cmdInput.EngineVersions = tt.engineVersions

			// when
			actual := isEngineVersionsValid()

			// then
			if actual != tt.expected {
				t.Errorf("expected %t, got %t", tt.expected, actual)
			}
		})
	}
}

func TestValidPathsShouldMakeValidConfig(t *testing.T) {
	// given
	base := t.TempDir()
	enginePath := filepath.Join(base, "Engine")
	outputPath := filepath.Join(base, "Output")
	if err := os.MkdirAll(enginePath, 0755); err != nil {
		t.Fatal("Failed to create temp dir for Engine.")
	}
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		t.Fatal("Failed to create temp dir for Output.")
	}
	upluginPath := filepath.Join(enginePath, "MyPlugin.uplugin")
	f, err := os.Create(upluginPath)
	if err != nil {
		t.Fatal("Failed to create uplugin.")
	}
	defer f.Close()
	var config = model.Config{
		EngineBaseDirectory: enginePath,
		OutputBaseDirectory: outputPath,
		BuildScriptPath:     "not empty",
		PluginPath:          upluginPath,
	}

	// when
	actual := isConfigValid(&config)

	// then
	if !actual {
		t.Error("A valid configuration should not fail the validation.")
	}
}

func TestMissingEngineBaseDirectoryShouldFailValidation(t *testing.T) {
	// given
	base := t.TempDir()
	enginePath := filepath.Join(base, "Engine")
	outputPath := filepath.Join(base, "Output")
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		t.Fatal("Failed to create temp dir for Output.")
	}
	var config = model.Config{
		EngineBaseDirectory: enginePath,
		OutputBaseDirectory: outputPath,
		BuildScriptPath:     "not empty",
	}

	// when
	actual := isConfigValid(&config)

	// then
	if actual {
		t.Error("A missing engine directory should fail the validation.")
	}
}

func TestMissingOutputDirectoryShouldFailValidation(t *testing.T) {
	// given
	base := t.TempDir()
	enginePath := filepath.Join(base, "Engine")
	outputPath := filepath.Join(base, "Output")
	if err := os.MkdirAll(enginePath, 0755); err != nil {
		t.Fatal("Failed to create temp dir for Engine.")
	}
	var config = model.Config{
		EngineBaseDirectory: enginePath,
		OutputBaseDirectory: outputPath,
		BuildScriptPath:     "not empty",
	}

	// when
	actual := isConfigValid(&config)

	// then
	if actual {
		t.Error("A missing output directory should fail the validation.")
	}
}

func TestMissingBuildScriptPathShouldFailValidation(t *testing.T) {
	// given
	base := t.TempDir()
	enginePath := filepath.Join(base, "Engine")
	outputPath := filepath.Join(base, "Output")
	if err := os.MkdirAll(enginePath, 0755); err != nil {
		t.Fatal("Failed to create temp dir for Engine.")
	}
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		t.Fatal("Failed to create temp dir for Output.")
	}
	var config = model.Config{
		EngineBaseDirectory: enginePath,
		OutputBaseDirectory: outputPath,
		BuildScriptPath:     "",
	}

	// when
	actual := isConfigValid(&config)

	// then
	if actual {
		t.Error("A missing build script path should fail the validation.")
	}
}

func TestCreateAndValidateConfig(t *testing.T) {
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
	upluginPath := filepath.Join(engineDir, "MyPlugin.uplugin")
	f, err := os.Create(upluginPath)
	if err != nil {
		t.Fatal("Failed to create uplugin.")
	}
	defer f.Close()

	// create expected config
	expected := model.Config{
		EngineBaseDirectory: engineDir,
		BuildScriptPath:     "RunUAT.bat",
		OutputBaseDirectory: outputDir,
		PluginPath:          upluginPath,
	}
	data, marshallErr := json.MarshalIndent(expected, "", " ")
	if marshallErr != nil {
		t.Fatal("Failed to marshall config file")
	}
	if err := os.WriteFile(configPath, data, 0755); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// when
	actual, _ := createAndValidateConfig(configPath)

	// then
	if *actual != expected {
		t.Errorf("Configs differ. Expected: %v, Actual: %v", expected, actual)
	}
}

func TestMissingConfigShouldReturnError(t *testing.T) {
	// given
	tempDir := t.TempDir()

	// when
	config, err := createAndValidateConfig(tempDir)

	// then
	if config != nil {
		t.Error("Config should be nil")
	}
	if err == nil {
		t.Error("Validation should have returned an error")
	}
}

func TestInvalidConfigShouldReturnError(t *testing.T) {
	// given
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")
	expected := model.Config{
		EngineBaseDirectory: "some dir",
		BuildScriptPath:     "RunUAT.bat",
		OutputBaseDirectory: "another dir",
	}
	data, marshallErr := json.MarshalIndent(expected, "", " ")
	if marshallErr != nil {
		t.Fatal("Failed to marshall config file")
	}
	if err := os.WriteFile(configPath, data, 0755); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// when
	_, err := createAndValidateConfig(configPath)

	// then
	if err == nil {
		t.Error("Invalid config should send an error")
	}
}

func TestEngineBaseDirAndOutputDirMustDiffer(t *testing.T) {
	// given
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")
	directory := filepath.Join(tempDir, "Directory")
	if err := os.MkdirAll(directory, 0755); err != nil {
		t.Fatal("Failed to create dir")
	}

	// create expected config
	expected := model.Config{
		EngineBaseDirectory: directory,
		BuildScriptPath:     "RunUAT.bat",
		OutputBaseDirectory: directory,
	}
	data, marshallErr := json.MarshalIndent(expected, "", " ")
	if marshallErr != nil {
		t.Fatal("Failed to marshall config file")
	}
	if err := os.WriteFile(configPath, data, 0755); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// when
	_, err := createAndValidateConfig(configPath)

	// then
	if err == nil {
		t.Error("Output cannot be the same as the engine base directory.")
	}
}

// create data for tests
func createPluginAtTempDir(pluginName string, t *testing.T) string {
	t.Helper()

	base := t.TempDir()
	pluginFile := filepath.Join(base, pluginName)
	file, err := os.Create(pluginFile)
	if err != nil {
		t.Fatalf("Failed to create temp dir %q", pluginFile)
	}
	file.Close()
	return pluginFile
}

func createEngineVersionsTestData() []engineVersionTestData {
	return []engineVersionTestData{
		{
			engineVersions: "5.4",
			expected:       true,
		},
		{
			engineVersions: "5.4,5.5",
			expected:       true,
		},
		{
			engineVersions: "4.11,5.1,5.2,5.6",
			expected:       true,
		},
		{
			engineVersions: "5",
			expected:       false,
		},
		{
			engineVersions: "",
			expected:       false,
		},
		{
			engineVersions: "invalid",
			expected:       false,
		},
		{
			engineVersions: "5,4",
			expected:       false,
		},
	}
}

func createIsInputValidTestData() []validInputTestData {
	return []validInputTestData{
		{
			engineVersions: "5.4",
			expected:       true,
		},
		{
			engineVersions: "",
			expected:       false,
		},
	}
}
