//go:build tools
// +build tools

// nolint
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
package tools

import (
	_ "github.com/gojuno/minimock/v3/cmd/minimock"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
)
