package win

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

// GetConsoleTitle sets windows console title
func GetConsoleTitle() (string, error) {
	buffer := make([]uint16, 256)
	ret, _, err := kernel32.NewProc("GetConsoleTitleW").Call(
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(len(buffer)),
	)
	if ret == 0 {
		return "", err
	}
	return windows.UTF16ToString(buffer), nil
}

// SetConsoleTitle sets windows console title
func SetConsoleTitle(title string) error {
	ret, _, err := kernel32.NewProc("SetConsoleTitleW").Call(
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(title))),
	)
	if ret == 0 {
		return err
	}
	return nil
}
