package win

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

// MsgBox items types
var (
	MsgBoxBtnOk                uint = 0x000000
	MsgBoxBtnOkCancel          uint = 0x000001
	MsgBoxBtnAbortRetryIgnore  uint = 0x000002
	MsgBoxBtnYesNoCancel       uint = 0x000003
	MsgBoxBtnYesNo             uint = 0x000004
	MsgBoxBtnRetryCancel       uint = 0x000005
	MsgBoxBtnCancelTryContinue uint = 0x000006
	MsgBoxIconNone             uint = 0x000000
	MsgBoxIconError            uint = 0x000010
	MsgBoxIconQuestion         uint = 0x000020
	MsgBoxIconWarning          uint = 0x000030
	MsgBoxIconInformation      uint = 0x000040
	MsgBoxDefaultButton1       uint = 0x000000
	MsgBoxDefaultButton2       uint = 0x000100
	MsgBoxDefaultButton3       uint = 0x000200
	MsgBoxDefaultButton4       uint = 0x000300
	MsgBoxTopMost              uint = 0x041000
	MsgBoxService              uint = 0x200000

	MsgBoxSelectOk       = 1
	MsgBoxSelectCancel   = 2
	MsgBoxSelectAbort    = 3
	MsgBoxSelectRetry    = 4
	MsgBoxSelectIgnore   = 5
	MsgBoxSelectYes      = 6
	MsgBoxSelectNo       = 7
	MsgBoxSelectTry      = 10
	MsgBoxSelectContinue = 11
)

// MsgBox create message box
func MsgBox(title string, msg string, flag uint) (int, error) {
	rTitle, err := windows.UTF16PtrFromString(title)
	if err != nil {
		return 0, err
	}

	rMsg, err := windows.UTF16PtrFromString(msg)
	if err != nil {
		return 0, err
	}

	ret, _, err := user32.NewProc("MessageBoxW").Call(
		0,
		uintptr(unsafe.Pointer(rMsg)),
		uintptr(unsafe.Pointer(rTitle)),
		uintptr(flag),
	)
	if ret == 0 {
		return 0, err
	}

	return int(ret), nil
}
