// +build windows

package win

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

type rtlOSVersionInfo struct {
	dwOSVersionInfoSize uint32
	dwMajorVersion      uint32
	dwMinorVersion      uint32
	dwBuildNumber       uint32
	dwPlatformId        uint32
	szCSDVersion        [128]byte
}

// Version returns Windows OS version
// TODO: Replace with `windows.GetVersion()` when this is resolved: https://github.com/golang/go/issues/17835
func Version() (major, minor, build uint32) {
	var verStruct rtlOSVersionInfo

	ntoskrnl := windows.MustLoadDLL("ntoskrnl.exe")
	defer ntoskrnl.Release()

	proc := ntoskrnl.MustFindProc("RtlGetVersion")
	verStruct.dwOSVersionInfoSize = uint32(unsafe.Sizeof(verStruct))
	proc.Call(uintptr(unsafe.Pointer(&verStruct)))

	return verStruct.dwMajorVersion, verStruct.dwMinorVersion, verStruct.dwBuildNumber
}
