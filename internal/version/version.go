// this file version.go was generated with go generate command

package version

import (
	"strings"
)

var info *Info

type Info struct {
	Version    string `json:"version"`
	CommitHash string `json:"commitHash"`
	Date       string `json:"date"`
}

func init() {
	info = &Info{
		Version:    "v0.0.0",
		CommitHash: "9bd0dcc13de29e90b1467d312a881f43aaf8771d",
		Date:       "2024-07-12T20:26:38Z",
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
