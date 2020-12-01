package win

import (
	"golang.org/x/sys/windows"
)

// Version identifies a Windows version by major, minor, and build number.
type Version struct {
	Major int
	Minor int
	Build int
}

// GetVersion returns the Windows version information. Applications not
// manifested for Windows 8.1 or Windows 10 will return the Windows 8 OS version
// value (6.2).
// For a table of version numbers see:
// https://msdn.microsoft.com/en-us/library/windows/desktop/ms724833(v=vs.85).aspx
func GetVersion() (Version, error) {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/ms724439(v=vs.85).aspx
	ver, err := windows.GetVersion()
	if err != nil {
		return Version{}, err
	}

	return Version{
		Major: int(ver & 0xFF),
		Minor: int(ver >> 8 & 0xFF),
		Build: int(ver >> 16),
	}, nil
}
