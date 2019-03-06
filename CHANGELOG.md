# Changelog

## 1.19.2 (2019/03/06)

* Omit app config if empty
* Check underlying app config value

## 1.19.1 (2019/03/06)

* Review and fix mutex
* Display message box on fatal error
* Display message box if main process is not found
* Add Portapps core version in `portapp.json`

## 1.19.0 (2019/03/05)

* Implement global configuration file (Issue #12)
* Add Windows MsgBox helpers
* Reorganize package
* Update libs

## 1.18.0 (2019/02/17)

* Add checksums file to releases (Issue #11)
* Switch to TravisCI

## 1.17.0 (2019/01/08)

* Add `CreateMutex` function

## 1.16.1 (2018/11/16)

* Fix error while compressing using UPX

## 1.16 (2018/10/26)

* Allow to enable CGO
* Review `CreateShortcut` function
* Add functions `Exists`, `IsDirEmpty`, `RawWinver`, `ReplaceByPrefix`, `WriteToFile`
* Add build type single

## 1.15 (2018/10/03)

* Add functions `CreateShortcut`, `SetFileAttributes`, `CopyFile`, `CopyFolder`, `RemoveContents`

## 1.14 (2018/09/30)

* Update InnoSetup to 5.6.1
* Update UPX to 3.95

## 1.13 (2018/09/29)

* `RegAdd` function added to create a registry key
* Allow to pass overwrite option while extracting files (default `-aoa`)
* Remove nsis-muarch build mode (use archive mode instead)

## 1.12 (2018/09/19)

* Assume Yes on all queries while extracting files
* Add go.sum

## 1.11 (2018/09/07)

* Upgrade to Go 1.11
* Use [go mod](https://golang.org/cmd/go/#hdr-Module_maintenance) instead of `dep`
* Fix asar lib version
* Disable strict-ssl on asar
* Update libs

## 1.10 (2018/04/28)

* Duplicated prepare target
* Update 7zip download link

## 1.9 (2018/03/13)

* Update libs
* Add QuickExecCmd function
* Add SetConsoleTitle function to set window console title
* Exclude version.dat from deletion
* Allow to use custom version for electron apps
* Add a check while creating folders
* ldflags can be customized

## 1.8 (2018/02/12)

* Add portapp.json file
* Uncheck run app in setup
* Mutualise release.app task
* App version not set in executable
* New artifact target (atf.win3264) for multi arch apps
* Move ia32/x64 to win32/win64 for arch def
* Add file creation, format unix / windows path
* Remove nupkg file
* Remove unnecessary files if not in debug mode

## 1.7 (2018/02/08)

* Ability to pass custom args to the portable process
* Update to rcedit 1.0.0
* Remove appasar from global electron build
* Add setelectronuserdata macro and load.lib.asar target
* Allow multi replacements in app.asar
* Ability to process a custom asar file

## 1.6 (2018/01/13)

* Replace userData in electron.asar for electron build type

## 1.5 (2017/12/14)

* Remove atf.original from untouched releases (can be found in the untouched repo)
* Add output type (7z or exe)

## 1.4 (2017/11/24)

* Disable UPX in debug mode
* Mutualize InnoSetup tpls
* Add NSIS Multi arch type
* Add prepare and finalize targets

## 1.3 (2017/11/23)

* Add untouched type

## 1.2 (2017/11/21)

* Coding style
* Add type archive, electron and innosetup
* Add debug mode
* Move app to a subfolder
* Manage registry keys

## 1.1 (2017/11/19)

* Disable basepath
* Fix execution stub resolution
* Grab LICENSE file from app repository
* Add CHANGELOG to release archive / setup
* Check artifact arch
* Setup default dir to C:\portapps\app
* Clear default values

## 1.0 (2017/11/16)

* Initial version
