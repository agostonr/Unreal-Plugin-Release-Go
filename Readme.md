### What it's for

Building unreal engine code plugins in a batch, targeting the specified versions

### What extra do I need
 - a `config.json` file that defines the engine base directory (until and without the version name)
 - the build script path in the engine directory
 - a base directory for the output
 - a path to copy the documentation pdf from

### How to use

.\unreal-plugin-release.exe --engine-versions=5.2,5.3,5.4,5.5 --plugin-location="E:\ProjectFiles\unreal\AimAssistPro\Plugins\AimAssistScripts\AimAssistScripts.uplugin"