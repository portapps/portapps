package win

import (
	"syscall"
	"unsafe"
)

// Locale returns the locale set for the user and falls back to the locale
// set for the system if unset. Returns "en-US" if all of them are unset.
func Locale() string {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	proc := kernel32.NewProc("GetUserDefaultLocaleName")
	buffer := make([]uint16, 128)
	ret, _, _ := proc.Call(uintptr(unsafe.Pointer(&buffer[0])), uintptr(len(buffer)))
	if ret == 0 {
		proc = kernel32.NewProc("GetSystemDefaultLocaleName")
		ret, _, _ = proc.Call(uintptr(unsafe.Pointer(&buffer[0])), uintptr(len(buffer)))
		if ret == 0 {
			return "en-US"
		}
	}
	return syscall.UTF16ToString(buffer)
}
