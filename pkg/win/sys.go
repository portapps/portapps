package win

import (
	"os"
	"strings"
)

// Is64Arch detects if program running on 64bits architecture
func Is64Arch() bool {
	return strings.ContainsAny("64", os.Getenv("PROCESSOR_ARCHITECTURE"))
}
