// +build windows

package mutex

import (
	"strings"
	"syscall"
	"unsafe"

	"github.com/google/logger"
)

type Mutex struct {
	handle uintptr
	log    *logger.Logger
}

// Create creates a mutex object for ensuring that only one instance is open.
// https://docs.microsoft.com/en-us/windows/desktop/api/synchapi/nf-synchapi-createmutexw
func Create(name string, logger *logger.Logger) (*Mutex, error) {
	var sbName strings.Builder
	sbName.WriteString("Portapps")
	sbName.WriteString(name)

	handle, err := syscall.LoadLibrary("kernel32.dll")
	if err != nil {
		return nil, err
	}
	defer syscall.FreeLibrary(handle)

	proc, err := syscall.GetProcAddress(handle, "CreateMutexW")
	if err != nil {
		return nil, err
	}

	rName, err := syscall.UTF16PtrFromString(sbName.String())
	if err != nil {
		return nil, err
	}

	r, _, err := syscall.Syscall(proc, 3, 0, 0, uintptr(unsafe.Pointer(rName)))
	logger.Infof("Mutex created: %d, %d", int(r), int(err.(syscall.Errno)))

	if int(err.(syscall.Errno)) == 0 {
		return &Mutex{
			handle: r,
			log:    logger,
		}, nil
	}

	return nil, err
}

// Release releases previously created mutex based on id
// https://docs.microsoft.com/en-us/windows/desktop/api/synchapi/nf-synchapi-releasemutex
func (m *Mutex) Release() error {
	m.log.Infof("Releasing mutex %d", int(m.handle))
	handle, err := syscall.LoadLibrary("kernel32.dll")
	if err != nil {
		return err
	}
	defer syscall.FreeLibrary(handle)

	proc, err := syscall.GetProcAddress(handle, "ReleaseMutex")
	if err != nil {
		return err
	}

	_, _, err = syscall.Syscall(proc, 1, m.handle, 0, 0)
	if int(err.(syscall.Errno)) == 0 {
		return nil
	}

	return err
}
