### What it's for
Building unreal engine code plugins in a batch, targeting the specified versions

### What extra do I need
 - a `config.json` file **in the same folder as the exe** that defines 
   - the engine base directory (until and without the version name)
   - the build script path within the engine directory
   - a base directory for the output

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
   - engine versions (comma separated, no whitespace)
   - plugin location (full path incl. `uplugin` file) 
   - _optional_ docs path (full path, also incl. file)

**Example:**
.\unreal-plugin-release.exe --engine-versions=5.2,5.3,5.4,5.5,5.6 \
--plugin-location="E:\ProjectFiles\unreal\MyPluginProject\Plugins\MyPlugin\MyPlugin.uplugin" \
--docs-path="E:\ProjectFiles\paperwork\MyPluginDocs\My_Plugin_Documentation.pdf"