package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/bketelsen/surgeon"
	"github.com/bketelsen/surgeon/codemods"

	"github.com/go-git/go-git/v5"
)

type Patient struct {
	Config       surgeon.Config
	ForkRoot     string
	UpsreamRoot  string
	forkRepo     *git.Repository
	upstreamRepo *git.Repository
}

func NewPatient(config surgeon.Config) *Patient {
	dir, _ := os.Getwd()
	return &Patient{
		Config:   config,
		ForkRoot: dir,
	}
}

func (p *Patient) Operate() error {

	slog.Debug("Opening fork git repository", "path", p.ForkRoot)
	r, err := git.PlainOpen(p.ForkRoot)
	if err != nil {
		slog.Error("opening git repository", "error", err)
		return fmt.Errorf("opening git repository: %w", err)
	}
	err = p.sanityCheck()
	if err != nil {
		slog.Error("sanity check failed", "error", err)
		return fmt.Errorf("sanity check failed: %w", err)
	}
	p.forkRepo = r

	// update the local fork
	slog.Debug("Updating local fork")
	err = p.updateLocalFork()
	if err != nil {
		slog.Error("updating local fork", "error", err)
		return fmt.Errorf("updating local fork: %w", err)
	}

	// clone the upstream repository
	slog.Debug("Cloning upstream repository")
	err = p.Clone()
	if err != nil {
		slog.Error("cloning upstream repository", "error", err)
		return fmt.Errorf("cloning upstream repository: %w", err)
	}
	defer os.RemoveAll(p.UpsreamRoot)

	slog.Debug("Applying code mods")
	for _, mod := range p.Config.CodeMods {
		cm, ok := codemods.Mods[mod.Mod]
		if !ok {
			slog.Error("code mod not found", "mod", mod.Mod)
			return fmt.Errorf("code mod %s not found", mod.Mod)
		}
		slog.Info("Applying codemod", "mod", mod.Mod, "description", mod.Description)
		slog.Debug("Validating code mod")
		err = cm.Validate(p.UpsreamRoot, p.ForkRoot, mod.Match, mod.Args...)
		if err != nil {
			slog.Error("validating code mod", "error", err)
			return fmt.Errorf("validating code mod: %w", err)
		}
		slog.Debug("Applying code mod")
		err = cm.Apply(p.UpsreamRoot, p.ForkRoot, mod.Match, mod.Args...)
		if err != nil {
			slog.Error("applying code mod", "error", err)
			return fmt.Errorf("applying code mod: %w", err)
		}

	}

	slog.Info("Comparing directories")
	missing, err := compareDirs(p.UpsreamRoot, p.ForkRoot)
	if err != nil {
		slog.Error("comparing directories", "error", err)
		return fmt.Errorf("comparing directories: %w", err)
	}
	if len(missing) > 0 {
		slog.Info("Found potential missing files", "count", len(missing))
		for _, m := range missing {
			if !strings.HasPrefix(m, ".git") {
				slog.Debug("Testing file", "file", m)
				if !p.IsIgnored(m) {
					slog.Debug("Copying file", "file", m)
					err = copyFile(m, p.UpsreamRoot, p.ForkRoot)
					if err != nil {
						return fmt.Errorf("copying file: %w", err)
					}
				} else {
					slog.Info("Skipping ignored", "file", m)
				}
			}
		}

	}

	// get the changed files in the upstream repository
	slog.Debug("Getting status of upstream repository")
	w, err := p.upstreamRepo.Worktree()
	if err != nil {
		slog.Error("getting git worktree", "error", err)
		return err
	}
	status, err := w.Status()
	if err != nil {
		slog.Error("getting worktree status", "error", err)
		return err
	}
	//fmt.Println(status)
	for s := range status {
		//	fmt.Println(s)
		// copy the file from the upstream repository to the fork
		if !p.IsIgnored(s) {
			slog.Debug("Copying file", "file", s)
			err = copyFile(s, p.UpsreamRoot, p.ForkRoot)
			if err != nil {
				slog.Error("copying file", "error", err)
				return fmt.Errorf("copying file: %w", err)
			}
		}
	}

	return nil
}

func (p *Patient) sanityCheck() error {

	ok, err := p.isClean()
	if err != nil {
		return fmt.Errorf("error checking if %s is in clean status: %w", p.ForkRoot, err)
	}
	if !ok {
		return fmt.Errorf("git repo %s is dirty", p.ForkRoot)
	}

	return nil
}

// isClean checks if the git repository is clean
func (p *Patient) isClean() (bool, error) {
	r, err := git.PlainOpen(p.ForkRoot)
	if err != nil {
		return false, err
	}
	w, err := r.Worktree()
	if err != nil {
		return false, err
	}
	status, err := w.Status()
	if err != nil {
		return false, err
	}
	return status.IsClean(), nil
}

func (p *Patient) updateLocalFork() error {
	w, err := p.forkRepo.Worktree()
	if err != nil {
		return err
	}
	err = w.Pull(&git.PullOptions{
		RemoteName: "origin",
	})
	if err != nil {
		if err != git.NoErrAlreadyUpToDate {
			return err
		}
	}
	return nil
}

func (p *Patient) Clone() error {
	var err error
	p.UpsreamRoot, err = os.MkdirTemp("", "surgeonupstream")
	if err != nil {
		return fmt.Errorf("creating temporary directory: %w", err)
	}
	slog.Info("Cloning upstream repository ", "location", p.UpsreamRoot)
	slog.Info("Cloning upstream repository ", "url", p.Config.Upstream)
	p.upstreamRepo, err = git.PlainClone(p.UpsreamRoot, false, &git.CloneOptions{
		URL: p.Config.Upstream,
	})
	if err != nil {
		return fmt.Errorf("cloning upstream repository: %w", err)
	}
	return nil
}

func copyFile(path, source, target string) error {
	sourcePath := filepath.Join(source, path)
	targetPath := filepath.Join(target, path)
	slog.Debug("Copying file", "source", sourcePath, "target", targetPath)
	err := os.MkdirAll(filepath.Dir(targetPath), 0755)
	if err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("opening source file: %w", err)
	}
	defer sourceFile.Close()
	targetFile, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("creating target file: %w", err)
	}
	defer targetFile.Close()
	_, err = io.Copy(targetFile, sourceFile)
	if err != nil {
		return fmt.Errorf("copying file: %w", err)
	}
	return nil

}

// compareDirs returns a list of files that are missing from the target directory
// compared to the source directory
func compareDirs(source, target string) ([]string, error) {
	var missing []string
	err := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == ".git" {
			slog.Debug("skipping git directory", "dir", info.Name())
			return filepath.SkipDir
		}
		rel, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(target, rel)
		_, err = os.Stat(targetPath)
		if err != nil {
			if os.IsNotExist(err) {
				fi, err := os.Stat(path)
				if err != nil {
					return err
				}
				if !fi.IsDir() {
					missing = append(missing, rel)
				}

			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return missing, nil
}

func (p *Patient) IsIgnored(path string) bool {
	for _, i := range p.Config.IgnoreList {
		slog.Debug("Checking Ignore List", "path", path, "prefix", i.Prefix)
		if strings.HasPrefix(path, i.Prefix) {
			slog.Debug("Ignoring", "path", path)
			return true
		}
	}
	return false
}
