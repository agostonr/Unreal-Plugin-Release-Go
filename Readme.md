### What's it for
Building unreal engine plugins in a batch, targeting the specified versions

### Notice
The SHA-256 hash of the built `PluginBuilder.exe` is `5c5c1a8def7ace270c09451e6cf9d2132651914e92f9a3aeb28ca21911e9f2a3`  
Before you use the pre-built executable, please verify its contents using A SHA-256 checksum calculator [like this one](https://emn178.github.io/online-tools/sha256_checksum.html). 

Alternatively, build the project yourself by [installing go](https://go.dev/doc/install)
and in the project root, running
```
go build -ldflags "-s -w" -o PluginBuilder.exe
```

### Platform support

The app was originally written for Windows, but it can be made to support Unix operating systems too.
The code itself is should run on either platform, but the app uses subprocesses (through `exec.Cmd`) on two occasions:
 - to invoke Unreal's build script `RunUAT.bat`
 - to zip the released plugin for upload (powershell's compress archive).

To enable this on Unix, implement `executor/executor_unix.go`. This executor interface is added to the plugin builder via constructor injection, and it should automatically select the appropriate OS based on GOOS, so there should be no additional tasks to do.

### What extra do I need
 - a `config.json` file **in the same folder as the exe** that defines 
   - `engineBaseDirectory`: the engine base directory (until and without the version name)
   - `buildScriptPath`: the build script path within the engine directory
   - `outputBaseDirectory`: a base directory for the output
   - `pluginPath`: the path to the uplugin file to build
   - `docsPath`: (optional) the documentation path

Example `config.json`:  
```
{
  "engineBaseDirectory": "D:\\Games",
  "buildScriptPath": "Engine\\Build\\BatchFiles\\RunUAT.bat",
  "outputBaseDirectory": "D:\\ProjectFiles\\unreal\\Release\\MyPlugin",
  "pluginPath": "D:\\ProjectFiles\\unreal\\MyProject\\Plugins\\MyPlugin\\MyPlugin.uplugin",
  "docsPath": "D:\\ProjectFiles\\paperwork\\MyPluginDocs\\My_Plugin_Docs.pdf"
}
```
  
 - a `FilterPlugin.ini` file **in the same folder as the exe**
   - **ONLY if you also add documentation**  
  Contents something like
  ```
  [FilterPlugin]
  /Documentation/My_Plugin_Documentation.pdf
  ```
where this path, under the release dir, will be created, and the provided docs file will be copied into it, with this name.  

The executable was meant to be one per plugin, because it has its own unique config file that concerns that plugin only. The versions are added as command line arguments.  

### How to use
 - build the project if you haven't already
 - invoke the exe file with the 
   - engine versions (comma separated, no whitespace), e.g. `5.1,5.2,5.3`
   - optional `--skip-docs` flag if you don't want to include docs, in spite of having it in the config

**Example (windows):**  

With documentation, if its path is set in the config file:   
```
.\PluginBuilder.exe --engine-versions=5.2,5.3,5.4,5.5,5.6 
```

Without documentation, even if the path it set in the config file:
```
.\PluginBuilder.exe --engine-versions=5.2,5.3,5.4,5.5,5.6 --skip-docs
```
