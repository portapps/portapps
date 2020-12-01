package win

import (
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

// RefreshEnv refresh Windows environment
// https://docs.microsoft.com/en-us/windows/desktop/api/winuser/nf-winuser-sendmessagetimeoutw
func RefreshEnv() error {
	buffer, err := windows.UTF16PtrFromString("Environment")
	if err != nil {
		return err
	}

	ret, _, err := user32.NewProc("SendMessageTimeoutW").Call(
		uintptr(0xFFFF),
		uintptr(0x001A),
		0,
		uintptr(unsafe.Pointer(buffer)),
		uintptr(0x0002),
		uintptr(5000))
	if ret == 0 {
		return err
	}

	return nil
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

// DeletePermEnv deletes an environment variable permanently on Windows
func DeletePermEnv(env registry.Key, name string) error {
	key, err := registry.OpenKey(env, "Environment", registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer key.Close()

	return key.DeleteValue(name)
}

// GetPermEnv gets an environment variable value on Windows
func GetPermEnv(env registry.Key, name string) (string, error) {
	key, err := registry.OpenKey(env, "Environment", registry.ALL_ACCESS)
	if err != nil {
		return "", nil
	}
	defer key.Close()

	val, _, err := key.GetStringValue(name)
	if err != nil {
		return "", err
	}

	return val, nil
}
