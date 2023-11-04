package win

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

// Locale returns the locale set for the user and falls back to the locale
// set for the system if unset. Returns "en-US" if all of them are unset.
func Locale() string {
	buffer := make([]uint16, 128)
	ret, _, _ := kernel32.NewProc("GetUserDefaultLocaleName").Call(
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(len(buffer)),
	)
	if ret == 0 {
		return "en-US"
	}
	return windows.UTF16ToString(buffer)
}
