### What it's for
Building unreal engine code plugins in a batch, targeting the specified versions

### What extra do I need
 - a `config.json` file **in the same folder as the exe** that defines 
   - `engineBaseDirectory`: the engine base directory (until and without the version name)
   - `buildScriptPath`: the build script path within the engine directory
   - `outputBaseDirectory`: a base directory for the output
   - `pluginPath`: the path to the uplugin file to build
   - `docsPath`: (optional) the documentation path
 - a `FilterPlugin.ini` file **in the same folder as the exe**
   - **ONLY if you also add documentation**
  Contents something like
  ```
  [FilterPlugin]
  /Documentation/My_Plugin_Documentation.pdf
  ```
where this path, under the release dir, will be created, and the provided docs file will be copied into it, with this name.

### How to use
 - build the project if you haven't already
 - invoke the exe file with the 
   - engine versions (comma separated, no whitespace), e.g. `5.1,5.2,5.3`
   - optional `--skip-docs` flag if you don't want to include docs, in spite of having it in the config

### Platform support
The code itself should work on either Windows, Mac or Linux, but the build script and zip are delegated to subprocesses using the exec.Cmd package. 
This is implemented for Windows, but for Unix, go to the `executor` package and implement the interface in `executor_unix.go`. The Executor itself should automatically pick one to inject into the builder component, so there should be no extra work to do there. 

**Example (windows):**  

With documentation, if its path is set in the config file:   
```
.\unreal-plugin-release.exe --engine-versions=5.2,5.3,5.4,5.5,5.6 
```

Without documentation, even if the path it set in the config file:
```
.\unreal-plugin-release.exe --engine-versions=5.2,5.3,5.4,5.5,5.6 --skip-docs
```