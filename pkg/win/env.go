package win

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows/registry"
)

const (
	HWND_BROADCAST   = 0xFFFF
	WM_SETTINGCHANGE = 0x001A
	SMTO_ABORTIFHUNG = 0x0002
)

// RefreshEnv refresh Windows environment
// https://docs.microsoft.com/en-us/windows/desktop/api/winuser/nf-winuser-sendmessagetimeoutw
func RefreshEnv() {
	env, _ := syscall.UTF16PtrFromString("Environment")
	syscall.NewLazyDLL("user32.dll").NewProc("SendMessageTimeoutW").Call(
		uintptr(HWND_BROADCAST),
		uintptr(WM_SETTINGCHANGE),
		0,
		uintptr(unsafe.Pointer(env)),
		uintptr(SMTO_ABORTIFHUNG),
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
