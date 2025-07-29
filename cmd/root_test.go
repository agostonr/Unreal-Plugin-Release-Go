package cmd

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

type engineVersionTestData struct {
	engineVersions string
	expected       bool
}

type validInputTestData struct {
	engineVersions string
	pluginLocation string
	expected       bool
}

func TestIsInputValid(t *testing.T) {
	for i, tt := range createIsInputValidTestData() {
		t.Run("IsInputValid #"+strconv.Itoa(i), func(t *testing.T) {
			// given
			if tt.pluginLocation != "" {
				cmdInput.PluginLocation = createPluginAtTempDir(tt.pluginLocation, t)
			}
			cmdInput.EngineVersions = tt.engineVersions

			// when
			actual := isInputValid()

			// then
			if actual != tt.expected {
				t.Errorf("Expected isValid: %t, actual isValid: %t", tt.expected, actual)
			}
		})
	}
}

func TestValidPluginLocation(t *testing.T) {
	// given
	cmdInput.PluginLocation = createPluginAtTempDir("MyPlugin.uplugin", t)

	// when
	actual := isPluginLocationValid()

	// then
	if !actual {
		t.Error("A valid plugin location must not be marked invalid.")
	}
}

func TestInvalidPluginLocation(t *testing.T) {
	// given
	base := t.TempDir()
	missingPath := filepath.Join(base, "MissingDirectory", "NoPlugin.uplugin")
	cmdInput.PluginLocation = missingPath

	// when
	actual := isPluginLocationValid()

	// then
	if actual {
		t.Error("A missing plugin location must be invalid.")
	}
}

func TestMissingPluginLocation(t *testing.T) {
	// given
	cmdInput.PluginLocation = ""

	// when
	actual := isPluginLocationValid()

	// then
	if actual {
		t.Error("An empty plugin location must be invalid.")
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

// create data for tests
func createPluginAtTempDir(pluginName string, t *testing.T) string {
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
			pluginLocation: "",
			expected:       false,
		},
		{
			engineVersions: "",
			pluginLocation: "MyPlugin.uplugin",
			expected:       false,
		},
		{
			engineVersions: "5.4",
			pluginLocation: "MyPlugin.uplugin",
			expected:       true,
		},
	}
}
