//go:build generate

//go:generate go run generate.go

package main

import (
	"encoding/json"
	"fmt"
	"github.com/caarlos0/log"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/afero"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

var fs afero.Fs

func init() {
	fs = afero.NewOsFs()
}

const (
	defaultTag      = "v0.0.0"
	txtName         = "VERSION"
	versionFileName = "version.go"

	templateFile = `// this file version.go was generated with go generate command

package version

import (
	"fmt"
	"strings"
)

var info *Info

type Info struct {
	Version    string   ` + "`json:\"version\"`" + `
	CommitHash string   ` + "`json:\"commitHash\"`" + `
	Date       string   ` + "`json:\"date\"`" + `
}

func init() {
	info = &Info{
		Version:    "{{.Version}}",
		CommitHash: "{{.CommitHash}}",
		Date:       "{{.Date}}",
		Features:   []string{},
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
`
)

type (
	Version struct {
		Version    string   `json:"version"`
		CommitHash string   `json:"commitHash"`
		Date       string   `json:"date"`
		Features   []string `json:"features,omitempty"`
	}

	Generator struct {
		projectPath string
		repo        *git.Repository
		commit      *object.Commit
	}
)

var afs = afero.NewOsFs()

func main() {
	projectPath, err := findGitRoot()
	if err != nil {
		log.Fatalf("error finding git root: %v", err)
		os.Exit(1)
	}

	ver, err := newVersion(projectPath)
	if err != nil {
		log.Fatalf("error creating version: %v", err)
		os.Exit(1)
	}

	if err = ver.Generate(); err != nil {
		log.Fatalf("error generating version: %v", err)
		os.Exit(1)
	}
}

func newVersion(projectPath string) (*Generator, error) {
	repo, err := git.PlainOpen(projectPath)
	if err != nil {
		return nil, err
	}

	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	p := &Generator{
		projectPath: projectPath,
		repo:        repo,
		commit:      commit,
	}

	return p, nil
}

// Generate creates a version.go file in the destination path
func (g *Generator) Generate() error {
	destPath := filepath.Join(g.projectPath, "internal", "version")

	// create folder if not exists
	if _, err := fs.Stat(destPath); os.IsNotExist(err) {
		if err = fs.MkdirAll(destPath, os.ModePerm); err != nil {
			return err
		}
	}

	versionFile := filepath.Join(destPath, versionFileName)

	log.Infof("generating go file: %s", versionFile)

	file, err := fs.Create(versionFile)
	if err != nil {
		return err
	}
	defer func(file afero.File) {
		if err = file.Close(); err != nil {
			log.Fatalf("error closing file: %g", err)
		}
	}(file)

	date := g.commit.Author.When.Format("2006-01-02T15:04:05Z")

	tag, err := g.getTag()
	if err != nil {
		return fmt.Errorf("error getting tag: %w", err)
	}

	data := &Version{
		Version:    tag,
		CommitHash: g.commit.Hash.String(),
		Date:       date,
	}

	tmpl, err := template.New(versionFile).Parse(templateFile)
	if err != nil {
		return fmt.Errorf("error parsing template: %w", err)
	}

	if err = tmpl.ExecuteTemplate(file, versionFile, data); err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	return g.genTxt(data)
}

// getTag returns the tag of the repository
func (g *Generator) getTag() (string, error) {
	tags, err := g.repo.Tags()
	if err != nil {
		return "", fmt.Errorf("error getting tags: %w", err)
	}

	tag := defaultTag

	callback := func(ref *plumbing.Reference) error {
		if ref.Name().IsTag() {
			tag = ref.Name().Short()
		}

		return nil
	}

	if err = tags.ForEach(callback); err != nil {
		return "", fmt.Errorf("error getting tags: %w", err)
	}

	return tag, nil
}

// genTxt creates a VERSION file with the version information
func (g *Generator) genTxt(ver *Version) error {
	txtFile := filepath.Join(g.projectPath, txtName)
	file, err := fs.Create(txtFile)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer func(file afero.File) {
		if err = file.Close(); err != nil {
			log.Fatalf("error closing file: %g", err)
		}
	}(file)

	log.Infof("generating version file: %s", txtFile)

	if err = json.NewEncoder(file).Encode(ver); err != nil {
		return fmt.Errorf("error encoding json: %w", err)
	}

	return nil
}

func findGitRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return filepath.Clean(strings.TrimSpace(string(out))), nil
}
