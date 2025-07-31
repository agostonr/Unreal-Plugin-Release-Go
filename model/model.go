// DTOs and model structs go here that are used by other packages.
package model

// represents the json configuration file
type Config struct {
	EngineBaseDirectory string `json:"engineBaseDirectory"`
	BuildScriptPath     string `json:"buildScriptPath"`
	OutputBaseDirectory string `json:"outputBaseDirectory"`
	PluginPath          string `json:"pluginPath"`
	DocsPath            string `json:"docsPath"`
}

type CmdInput struct {
	EngineVersions string
	SkipDocs       bool
}
