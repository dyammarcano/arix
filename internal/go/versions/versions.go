package versions

import (
	"encoding/json"
	"fmt"
	"github.com/blang/semver"
	"github.com/mholt/archiver/v4"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

const (
	linuxPath   = "/usr/local/go"
	windowsPath = "C:\\go"
)

const (
	goUrl = "https://go.dev/dl/?mode=json&include=all"
)

type Versions struct {
	ID      int    `json:"id,omitempty"`
	Version string `json:"version,omitempty"`
	Stable  bool   `json:"stable,omitempty"`
	Files   []File `json:"files,omitempty"`
}

type File struct {
	ID       int    `json:"id,omitempty"`
	Version  string `json:"version,omitempty"`
	Stable   bool   `json:"stable,omitempty"`
	Filename string `json:"filename,omitempty"`
	Os       string `json:"os,omitempty"`
	Arch     string `json:"arch,omitempty"`
	Sha256   string `json:"sha256,omitempty"`
	Size     int    `json:"size,omitempty"`
	Kind     string `json:"kind,omitempty"`
}

type GoVersion struct {
	StableVersion    string     `json:"stable,omitempty"`
	Versions         []Versions `json:"versions,omitempty"`
	ReleaseCandidate string     `json:"release_candidate,omitempty"`
}

// NewGoVersion returns a new GoVersion.
func NewGoVersion() (*GoVersion, error) {
	goVer, err := getJSON(goUrl)
	if err != nil {
		return nil, err
	}

	releaseCandidate := goVer.Versions[0].Version

	for i := range goVer.Versions {
		for j := range goVer.Versions[i].Files {
			if goVer.Versions[i].Files[j].Kind == "source" {
				goVer.Versions[i].Files[j].Os = "any"
				goVer.Versions[i].Files[j].Arch = "any"
			}
		}
	}

	sort.Slice(goVer.Versions, func(i, j int) bool {
		verI, _ := semver.Make(strings.TrimPrefix(goVer.Versions[i].Version, "go"))
		verJ, _ := semver.Make(strings.TrimPrefix(goVer.Versions[j].Version, "go"))
		return verI.GT(verJ)
	})

	return &GoVersion{
		StableVersion:    goVer.Versions[0].Version,
		ReleaseCandidate: releaseCandidate,
		Versions:         goVer.Versions,
	}, nil
}

func (g *GoVersion) GetStableVersion() string {
	return g.StableVersion
}

func (g *GoVersion) InstallVersion(version string) string {
	for i := range g.Versions {
		if g.Versions[i].Version == version {
			return g.Versions[i].Version
		}
	}
	return ""
}

func (g *GoVersion) InstallStableVersion() error {
	file := choseWitchGoPackage(g.StableVersion, g)
	if file == "" {
		return fmt.Errorf("no suitable Go package found for stable version")
	}

	if err := downloadAndInstall(file, linuxPath); err != nil {
		return fmt.Errorf("error downloading and installing Go file: %w", err)
	}
	return nil
}

// getJSON returns the GoVersion struct from the Go website
func getJSON(url string) (*GoVersion, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			fmt.Println(err)
		}
	}(r.Body)

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var goVer GoVersion
	if err = json.Unmarshal(data, &goVer.Versions); err != nil {
		return nil, err
	}
	return &goVer, nil
}

func downloadGoVersion(goUrl, filePath string) error {
	r, err := http.Get(goUrl)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			fmt.Println(err)
		}
	}(r.Body)

	fileName := strings.Split(goUrl, "/")
	filePath = filepath.Clean(fmt.Sprintf("%s/%s", filePath, fileName[len(fileName)-1]))

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		if err = out.Close(); err != nil {
			fmt.Println(err)
		}
	}(out)

	if _, err = io.Copy(out, r.Body); err != nil {
		return err
	}

	return handleInstallation(filePath)
}

func handleInstallation(filePath string) error {
	if strings.HasSuffix(filePath, ".msi") {
		return installMsi(filePath)
	} else if strings.HasSuffix(filePath, ".zip") {
		return unzipFile(filePath, windowsPath)
	} else {
		return fmt.Errorf("unsupported file type")
	}
}

func installMsi(filePath string) error {
	cmd := exec.Command("msiexec", "/i", filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func unzipFile(zipFile, dest string) error {
	return archiver.RegisterFormat(archiver.FormatZip, archiver.DefaultZip).Unarchive(zipFile, dest)
}

// choseWitchGoPackage selects the appropriate Go package based on OS, Arch, Filename, and Kind.
func choseWitchGoPackage(goVersion string, goVer *GoVersion) string {
	for _, version := range goVer.Versions {
		if version.Version == goVersion {
			for _, file := range version.Files {
				if file.Kind == "archive" && (file.Os == "linux" || file.Os == "windows") && (file.Arch == "amd64" || file.Arch == "386") {
					return file.Filename
				}
			}
		}
	}
	return ""
}

func downloadAndInstall(filename, path string) error {
	goUrl := fmt.Sprintf("https://go.dev/dl/%s", filename)
	return downloadGoVersion(goUrl, path)
}
