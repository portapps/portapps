package mutex

import (
	"strings"
	"syscall"
	"unsafe"

	"github.com/rs/zerolog/log"
)

// Mutex is a mutex instance
type Mutex struct {
	handle uintptr
}

// New creates a mutex object for ensuring that only one instance is open.
// https://docs.microsoft.com/en-us/windows/desktop/api/synchapi/nf-synchapi-createmutexw
func New(name string) (*Mutex, error) {
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
	log.Info().Msgf("Mutex created: %d, err%d", int(r), int(err.(syscall.Errno)))

	if int(err.(syscall.Errno)) == 0 {
		return &Mutex{
			handle: r,
		}, nil
	}

	return nil, err
}

// Release releases previously created mutex based on id.
// https://docs.microsoft.com/en-us/windows/desktop/api/synchapi/nf-synchapi-releasemutex
func (m *Mutex) Release() error {
	log.Info().Msgf("Releasing mutex %d", int(m.handle))
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
