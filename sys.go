package portapps

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"unsafe"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

// WindowsShortcut the Windows shortcut structure
type WindowsShortcut struct {
	ShortcutPath     string
	TargetPath       string
	Arguments        string
	Description      string
	IconLocation     string
	WorkingDirectory string
}

// CreateShortcut creates a windows shortcut
func CreateShortcut(shortcut WindowsShortcut) error {
	Log.Infof("Create shortcut for %s in %s...", shortcut.TargetPath, shortcut.ShortcutPath)
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_SPEED_OVER_MEMORY)
	defer ole.CoUninitialize()

	oleShellObject, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		return err
	}

	defer oleShellObject.Release()
	wshell, err := oleShellObject.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return err
	}

	defer wshell.Release()
	cs, err := oleutil.CallMethod(wshell, "CreateShortcut", shortcut.ShortcutPath)
	if err != nil {
		return err
	}

	idispatch := cs.ToIDispatch()
	oleutil.PutProperty(idispatch, "TargetPath", shortcut.TargetPath)
	oleutil.PutProperty(idispatch, `"{9F4C2855-9F79-4B39-A8D0-E1D42DE1D5F3}",5`, shortcut.TargetPath)
	if shortcut.Arguments != "" {
		oleutil.PutProperty(idispatch, "Arguments", shortcut.Arguments)
	}
	if shortcut.Description != "" {
		oleutil.PutProperty(idispatch, "Description", shortcut.Description)
	}
	if shortcut.IconLocation != "" {
		oleutil.PutProperty(idispatch, "IconLocation", shortcut.IconLocation)
	}
	if shortcut.WorkingDirectory != "" {
		oleutil.PutProperty(idispatch, "WorkingDirectory", shortcut.WorkingDirectory)
	}
	_, err = oleutil.CallMethod(idispatch, "Save")

	return err
}

// SetFileAttributes set attributes to a file
func SetFileAttributes(path string, attrs uint32) error {
	pointer, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return err
	}

	return syscall.SetFileAttributes(pointer, attrs)
}

// SetConsoleTitle sets windows console title
func SetConsoleTitle(title string) (int, error) {
	handle, err := syscall.LoadLibrary("kernel32.dll")
	if err != nil {
		return 0, err
	}
	defer syscall.FreeLibrary(handle)

	proc, err := syscall.GetProcAddress(handle, "SetConsoleTitleW")
	if err != nil {
		return 0, err
	}

	rTitle, err := syscall.UTF16PtrFromString(title)
	if err != nil {
		return 0, err
	}

	r, _, err := syscall.Syscall(proc, 1, uintptr(unsafe.Pointer(rTitle)), 0, 0)
	return int(r), err
}

// CopyFile copy a file
func CopyFile(src string, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	err = destFile.Sync()
	if err != nil {
		return err
	}

	return nil
}

// CopyFolder copy a folder
func CopyFolder(source string, dest string) (err error) {
	err = os.MkdirAll(dest, 777)
	if err != nil {
		return err
	}

	folder, _ := os.Open(source)
	objects, err := folder.Readdir(-1)
	for _, object := range objects {
		sourceFile := path.Join(source, object.Name())
		destFile := path.Join(dest, object.Name())
		if object.IsDir() {
			err = CopyFolder(sourceFile, destFile)
			if err != nil {
				return err
			}
		} else {
			err = CopyFile(sourceFile, destFile)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// RemoveContents remove contents of a specified directory
func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateFolderCheck to create a folder and get its path and return error
func CreateFolderCheck(path string) (string, error) {
	if err := os.MkdirAll(path, 777); err != nil {
		return "", err
	}
	return path, nil
}

// CreateFolder to create a folder and get its path
func CreateFolder(path string) string {
	Log.Infof("Create folder %s...", path)
	if _, err := CreateFolderCheck(path); err != nil {
		Log.Fatalf("Cannot create folder: %v", err)
	}
	return path
}

// CreateFile creates / overwrites a file with content
func CreateFile(path string, content string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err = file.Sync(); err != nil {
		return err
	}
	return nil
}

// PathJoin to join paths
func PathJoin(elem ...string) string {
	for i, e := range elem {
		if e != "" {
			return strings.Join(elem[i:], `\`)
		}
	}
	return ""
}

// AppPathJoin to join paths from Papp.Path
func AppPathJoin(elem ...string) string {
	return PathJoin(append([]string{Papp.Path}, elem...)...)
}

// FormatUnixPath to format a path for unix
func FormatUnixPath(path string) string {
	return strings.Replace(path, `\`, `/`, -1)
}

// FormatWindowsPath to format a path for windows
func FormatWindowsPath(path string) string {
	return strings.Replace(path, `/`, `\`, -1)
}
