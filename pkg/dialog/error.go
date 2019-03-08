package dialog

import (
	"fmt"
)

// MsgBoxError represents a message box error
type MsgBoxError uint32

// Error returns a formatted string error
func (e MsgBoxError) Error() string {
	return fmt.Sprintf("SystemErrorCode: %#x", e)
}
