package win

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows/registry"
)

// RefreshEnv refresh Windows environment
// https://docs.microsoft.com/en-us/windows/desktop/api/winuser/nf-winuser-sendmessagetimeoutw
func RefreshEnv() {
	env, _ := syscall.UTF16PtrFromString("Environment")
	syscall.NewLazyDLL("user32.dll").NewProc("SendMessageTimeoutW").Call(
		uintptr(0xFFFF),
		uintptr(0x001A),
		0,
		uintptr(unsafe.Pointer(env)),
		uintptr(0x0002),
		uintptr(5000))
}

// SetPermEnv set an environment variable permanently on Windows
func SetPermEnv(env registry.Key, name string, value string) error {
	key, err := registry.OpenKey(env, "Environment", registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer key.Close()

	return key.SetStringValue(name, value)
}
