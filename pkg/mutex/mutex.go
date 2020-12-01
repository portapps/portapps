package mutex

import (
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/sys/windows"
)

// Create creates a mutex object for ensuring that only one instance is open.
// https://docs.microsoft.com/en-us/windows/desktop/api/synchapi/nf-synchapi-createmutexw
func Create(name string) (windows.Handle, error) {
	var sbName strings.Builder
	sbName.WriteString("Portapps")
	sbName.WriteString(name)

	muName, err := windows.UTF16PtrFromString(sbName.String())
	if err != nil {
		return 0, err
	}

	handle, err := windows.OpenMutex(windows.MUTEX_ALL_ACCESS, false, muName)
	if err == nil {
		windows.CloseHandle(handle)
		return 0, errors.Errorf("already running")
	}

	return windows.CreateMutex(nil, false, muName)
}

// Release releases previously created mutex based on id.
// https://docs.microsoft.com/en-us/windows/desktop/api/synchapi/nf-synchapi-releasemutex
func Release(handle windows.Handle) error {
	return windows.CloseHandle(handle)
}
