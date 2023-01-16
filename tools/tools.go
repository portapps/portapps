//go:build tools
// +build tools

// Package tools tracks dependencies on binaries not referenced in this codebase.
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
package tools

import (
	_ "github.com/josephspurrier/goversioninfo/cmd/goversioninfo"
)
