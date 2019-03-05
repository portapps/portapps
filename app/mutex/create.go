// +build windows

package mutex

import (
	"strings"
	"syscall"
	"unsafe"
)

// Create creates a mutex object for ensuring that only one instance is open.
func Create(name string) (int, error) {
	var sbName strings.Builder
	sbName.WriteString("Portapps")
	sbName.WriteString(name)

	handle, err := syscall.LoadLibrary("kernel32.dll")
	if err != nil {
		return 0, err
	}
	defer syscall.FreeLibrary(handle)

	proc, err := syscall.GetProcAddress(handle, "CreateMutexW")
	if err != nil {
		return 0, err
	}

	rName, err := syscall.UTF16PtrFromString(sbName.String())
	if err != nil {
		return 0, err
	}

	r, _, err := syscall.Syscall(proc, 1, uintptr(unsafe.Pointer(rName)), 0, 0)
	return int(r), err
}
