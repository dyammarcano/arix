// this file version.go was generated with go generate command

package version

import (
	"runtime/debug"
	"strings"
)

var info *Info

type Info struct {
	Version    string `json:"version"`
	CommitHash string `json:"commitHash"`
	Date       string `json:"date"`
	GoVersion  string `json:"goVersion"`
}

func init() {
	nfo, ok := debug.ReadBuildInfo()
	if !ok || (nfo != nil && nfo.Main.Version == "") {
		return
	}

	info = &Info{
		Version:    "v0.0.0",
		CommitHash: "b8febb33e0859cc6526a5f9e95e4253baff56e0f",
		Date:       "2024-07-12T20:33:20Z",
		GoVersion:  nfo.Main.Version,
	}
}

// G returns the Info struct
func G() *Info {
	return info
}

// GetVersionInfo returns the info
func GetVersionInfo() string {
	var sb strings.Builder
	sb.WriteString(info.Version)

	if info.CommitHash != "" {
		sb.WriteString("-")
		sb.WriteString(info.CommitHash)
	}

	if info.Date != "" {
		sb.WriteString("-")
		sb.WriteString(info.Date[:10]) // format date to yyyy-mm-dd
	}

	return sb.String()
}
