package dbmodels

import (
	"path/filepath"
	"strings"
)

const (
	FileCurrentDir    = "."
	FilePathSeparator = string(filepath.Separator)
)

type Branch struct {
	Platform string `json:"platform" required:"true"`
	Org      string `json:"org" required:"true"`
	Repo     string `json:"repo" required:"true"`
	Branch   string `json:"branch" required:"true"`
}

func (b Branch) IsMissingParam() (string, bool) {
	v := true
	if b.Platform == "" {
		return "platform", v
	}

	if b.Org == "" {
		return "org", v
	}

	if b.Repo == "" {
		return "repo", v
	}

	if b.Branch == "" {
		return "branch", v
	}
	return "", false
}

type File struct {
	Path string `json:"path" required:"true"`
	SHA  string `json:"sha" required:"true"`

	// Allow empty file
	Content string `json:"content,omitempty"`
}

func (f File) Name() string {
	return filepath.Base(f.Path)
}

func (f File) Dir() string {
	s := strings.TrimRight(f.Path, FilePathSeparator)
	return filepath.Dir(s)
}

func (f *File) IsMissingParam() (string, bool) {
	if f == nil {
		return "", false
	}

	v := true

	if f.Path == "" {
		return "path", v
	}

	if f.SHA == "" {
		return "sha", v
	}
	return "", false
}
