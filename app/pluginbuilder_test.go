/*
Running the actual app would require some form of implementation
*/
package app

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"testing"

	"unreal-plugin-release/model"
)

// fake executor: run the entire build, tempdir the folders in the output and check if the app moved files and removed folders that are not needed.
type FakeExecutor struct {
}

func (e FakeExecutor) CreateBuilderCommand(buildScriptPath string, pluginLocation string, outputDir string) *exec.Cmd {
	const message = "Failed to write directory"
	// create a bunch of directories as if the plugin has been built.
	if err := os.MkdirAll(filepath.Join(outputDir, "Source"), 0755); err != nil {
		panic(message)
	}

	if err := os.MkdirAll(filepath.Join(outputDir, "Intermediate"), 0755); err != nil {
		panic(message)
	}

	if err := os.MkdirAll(filepath.Join(outputDir, "Binaries"), 0755); err != nil {
		panic(message)
	}

	if err := os.MkdirAll(filepath.Join(outputDir, "Build"), 0755); err != nil {
		panic(message)
	}

	if err := os.MkdirAll(filepath.Join(outputDir, "Content"), 0755); err != nil {
		panic(message)
	}

	if err := os.MkdirAll(filepath.Join(outputDir, "Resources"), 0755); err != nil {
		panic(message)
	}

	return createEmptyCommand()
}

func (e FakeExecutor) CreateZipCommand(sourceDir string) *exec.Cmd {
	// create an empty zip file, as if zipping the project went through
	dir := filepath.Dir(sourceDir)
	name := filepath.Base(sourceDir)
	f, err := os.Create(filepath.Join(dir, name+".zip"))
	if err != nil {
		panic("Couldn't create fake zip file")
	}
	defer f.Close()
	return createEmptyCommand()
}

func createEmptyCommand() *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.Command("cmd", "/c", "rem")
	}
	return exec.Command("true")
}

func TestGetUnneededFolders(t *testing.T) {
	// given
	config := model.Config{}
	underTest := NewPluginBuilder(&config, FakeExecutor{})
	expected := []string{"Binaries", "Build", "Intermediate", "Saved"}

	// when
	actual := underTest.getUnneededFolders()

	// then
	if !arrayContainsAll(expected, actual) {
		t.Error("The folders to remove from the build are incorrect.")
	}
}

func TestCollectVersions(t *testing.T) {
	// given
	config := model.Config{}
	underTest := NewPluginBuilder(&config, FakeExecutor{})
	expected := []string{"5.3", "5.4", "5.5"}
	input := "5.3,5.4,5.5"

	// when
	actual := underTest.collectVersions(input)

	// then
	if !arrayContainsAll(expected, actual) {
		t.Error("Engine versions were incorrectly split.")
	}
}

func TestMakeBuildScriptFilePath(t *testing.T) {
	// given
	version := "5.4"
	buildScriptRelativePath := "RunUAT.bat"
	base := t.TempDir()
	engine := makeDir(base, "Engine", t)
	writeBuildScript(engine, version, buildScriptRelativePath, t)
	writeFilterPluginFile(base, t)
	executor := FakeExecutor{}

	config := model.Config{
		EngineBaseDirectory: engine,
		BuildScriptPath:     buildScriptRelativePath,
	}
	underTest := NewPluginBuilder(&config, executor)
	expected := filepath.Join(engine, "UE_"+version, buildScriptRelativePath)

	// when
	actual := underTest.makeBuildScriptFilePath(version)

	// then
	if actual != expected {
		t.Errorf("Expected: %q, Actual: %q", expected, actual)
	}
}

func TestBuildPlugin(t *testing.T) {
	// given
	fixtures := "testdata"
	version := "5.4"
	buildScriptRelativePath := "RunUAT.bat"
	base := t.TempDir()
	engine := makeDir(base, "Engine", t)
	uplugin := makeFile(base, "MyPlugin.uplugin", t)
	writeBuildScript(engine, version, buildScriptRelativePath, t)
	docs := filepath.Join(fixtures, "testing.pdf")
	writeFilterPluginFile(base, t)
	output := makeDir(base, "Output", t)
	execPath := filepath.Join(base, "script.exe")
	executor := FakeExecutor{}
	builtPluginPath := filepath.Join(output, "MyPlugin_5.4")

	config := model.Config{
		EngineBaseDirectory: engine,
		BuildScriptPath:     buildScriptRelativePath,
		OutputBaseDirectory: output,
		PluginPath:          uplugin,
		DocsPath:            docs,
	}

	cmdInput := model.CmdInput{
		EngineVersions: version,
		SkipDocs:       false,
	}

	underTest := NewPluginBuilder(&config, executor)

	// when
	underTest.BuildPluginsForSelectedVersions(cmdInput, execPath)

	// then
	if isDirectoryExist(filepath.Join(builtPluginPath, "Binaries")) {
		t.Error("Binaries should have been removed from the release.")
	}

	if isDirectoryExist(filepath.Join(builtPluginPath, "Intermediate")) {
		t.Error("Intermediate should have been removed from the release.")
	}

	if isDirectoryExist(filepath.Join(builtPluginPath, "Build")) {
		t.Error("Build should have been removed from the release.")
	}

	if !isDirectoryExist(filepath.Join(builtPluginPath, "Source")) {
		t.Error("Source missing from release.")
	}

	if !isDirectoryExist(filepath.Join(builtPluginPath, "Resources")) {
		t.Error("Resources missing from release.")
	}

	if !isDirectoryExist(filepath.Join(builtPluginPath, "Content")) {
		t.Error("Content missing from release.")
	}

	if !isFileExist(filepath.Join(builtPluginPath, "Docs", "My_Docs.pdf")) {
		t.Error("Documentation wasn't copied or is not the right name.")
	}

	if !isFileExist(filepath.Join(builtPluginPath, "Config", "FilterPlugin.ini")) {
		t.Error("FilterPlugin file is not in the Config folder.")
	}
}

// helper for tests
func arrayContainsAll(expected []string, actual []string) bool {
	if len(expected) != len(actual) {
		return false
	}

	for _, item := range expected {
		if !slices.Contains(actual, item) {
			return false
		}
	}

	return true
}

func makeDir(base string, dir string, t *testing.T) string {
	t.Helper()
	dirPath := filepath.Join(base, dir)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		t.Fatalf("Failed to create directory at %q", dirPath)
	}
	return dirPath
}

func makeFile(base string, filename string, t *testing.T) string {
	t.Helper()
	filePath := filepath.Join(base, filename)
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		t.Fatalf("Failed to make directory %q", filePath)
	}
	file, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("Failed to create file at %q", filePath)
	}
	defer file.Close()
	return filePath
}

func writeFilterPluginFile(base string, t *testing.T) string {
	t.Helper()
	filterPluginContents, err := os.ReadFile(filepath.Join("testdata", "FilterPluginTest.ini"))
	if err != nil {
		t.Fatal("Failed to read filterplugin file")
	}

	result := filepath.Join(base, "FilterPlugin.ini")
	writeErr := os.WriteFile(result, filterPluginContents, 0755)
	if writeErr != nil {
		t.Fatal("Failed to write filter plugin file to temp dir")
	}

	return result
}

func writeBuildScript(engine, version, scriptName string, t *testing.T) string {
	t.Helper()
	path := filepath.Join(engine, "UE_"+version)
	return makeFile(path, scriptName, t)
}

func isDirectoryExist(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func isFileExist(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
