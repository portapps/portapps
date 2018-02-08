# Changelog

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
