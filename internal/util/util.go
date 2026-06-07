package util

import (
	"runtime/debug"
	"strings"
)

func GetModuleName() string {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}

	if len(bi.Deps) > 0 && bi.Deps[0].Path != "" {
		return bi.Deps[0].Path
	}

	return bi.Path
}

func ExtractPackageName(name string) string {
	parts := strings.Split(strings.TrimSpace(name), "/")
	return parts[len(parts)-1]
}
