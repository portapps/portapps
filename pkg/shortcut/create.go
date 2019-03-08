package shortcut

import (
	"runtime"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

// Shortcut the Windows shortcut structure
type Shortcut struct {
	ShortcutPath     string
	TargetPath       string
	Arguments        Property
	Description      Property
	IconLocation     Property
	WorkingDirectory Property
}

// Property the Windows shortcut property
type Property struct {
	Value string
	Clear bool
}

// Create creates a windows shortcut
func Create(shortcut Shortcut) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if err := ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_SPEED_OVER_MEMORY); err != nil {
		return err
	}
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
	if _, err := oleutil.PutProperty(idispatch, "TargetPath", shortcut.TargetPath); err != nil {
		return err
	}
	if shortcut.Arguments.Value != "" || shortcut.Arguments.Clear {
		if _, err := oleutil.PutProperty(idispatch, "Arguments", shortcut.Arguments.Value); err != nil {
			return err
		}
	}
	if shortcut.Description.Value != "" || shortcut.Description.Clear {
		if _, err := oleutil.PutProperty(idispatch, "Description", shortcut.Description.Value); err != nil {
			return err
		}
	}
	if shortcut.IconLocation.Value != "" || shortcut.IconLocation.Clear {
		if _, err := oleutil.PutProperty(idispatch, "IconLocation", shortcut.IconLocation.Value); err != nil {
			return err
		}
	}
	if shortcut.WorkingDirectory.Value != "" || shortcut.WorkingDirectory.Clear {
		if _, err := oleutil.PutProperty(idispatch, "WorkingDirectory", shortcut.WorkingDirectory.Value); err != nil {
			return err
		}
	}
	_, err = oleutil.CallMethod(idispatch, "Save")

	return err
}
