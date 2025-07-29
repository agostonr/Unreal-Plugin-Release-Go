// DTOs and model structs go here that are used by other packages.
package model

// represents the json configuration file
type Config struct {
	EngineBaseDirectory string `json:"engineBaseDirectory"`
	BuildScriptPath     string `json:"buildScriptPath"`
	OutputBaseDirectory string `json:"outputBaseDirectory"`
}

// command line input, different per plugin release
type CmdInput struct {
	EngineVersions string
	PluginLocation string
	DocsPath       string
}
