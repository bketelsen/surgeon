package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bketelsen/surgeon/codemods"
	"github.com/go-git/go-git/v5"
)

type Patient struct {
	Config       Config
	ForkRoot     string
	UpsreamRoot  string
	forkRepo     *git.Repository
	upstreamRepo *git.Repository
}

func NewPatient(config Config) *Patient {
	dir, _ := os.Getwd()
	return &Patient{
		Config:   config,
		ForkRoot: dir,
	}
}

func (p *Patient) Operate() error {

	r, err := git.PlainOpen(p.ForkRoot)
	if err != nil {
		return fmt.Errorf("opening git repository: %w", err)
	}
	err = p.sanityCheck()
	if err != nil {
		return fmt.Errorf("sanity check failed: %w", err)
	}
	p.forkRepo = r

	// update the local fork
	err = p.updateLocalFork()
	if err != nil {
		return fmt.Errorf("updating local fork: %w", err)
	}

	// clone the upstream repository
	err = p.Clone()
	if err != nil {
		return fmt.Errorf("cloning upstream repository: %w", err)
	}
	for _, mod := range p.Config.CodeMods {
		cm, ok := codemods.Mods[mod.Mod]
		if !ok {
			return fmt.Errorf("code mod %s not found", mod.Mod)
		}
		fmt.Printf("Applying codemod %s: %s\n", mod.Mod, mod.Description)
		err = cm.Validate(p.UpsreamRoot, p.ForkRoot, mod.Match, mod.Args...)
		if err != nil {
			return fmt.Errorf("validating code mod: %w", err)
		}
		err = cm.Apply(p.UpsreamRoot, p.ForkRoot, mod.Match, mod.Args...)
		if err != nil {
			return fmt.Errorf("applying code mod: %w", err)
		}

	}

	missing, err := compareDirs(p.UpsreamRoot, p.ForkRoot)
	if err != nil {
		return fmt.Errorf("comparing directories: %w", err)
	}
	if len(missing) > 0 {
		fmt.Println("The following files are missing from the fork")
		for _, m := range missing {
			if !strings.HasPrefix(m, ".git") {
				fmt.Println(m)
				if !p.IsIgnored(m) {
					err = copyFile(m, p.UpsreamRoot, p.ForkRoot)
					if err != nil {
						return fmt.Errorf("copying file: %w", err)
					}
				}
			}
		}
		return nil
	}

	// get the changed files in the upstream repository
	w, err := p.upstreamRepo.Worktree()
	if err != nil {
		return err
	}
	status, err := w.Status()
	if err != nil {
		return err
	}
	fmt.Println(status)
	for s := range status {
		fmt.Println(s)
		// copy the file from the upstream repository to the fork
		if !p.IsIgnored(s) {
			err = copyFile(s, p.UpsreamRoot, p.ForkRoot)
			if err != nil {
				return fmt.Errorf("copying file: %w", err)
			}
		}
	}

	var clean bool
	// if the user wants to stage the changes, do so
	if p.Config.Stage {
		fmt.Println("Staging changes")
		w, err := p.forkRepo.Worktree()
		if err != nil {
			return fmt.Errorf("getting git worktree: %w", err)
		}
		status, err := w.Status()
		if err != nil {
			return fmt.Errorf("getting worktree status: %w", err)
		}
		if status.IsClean() {
			fmt.Println("No changes to stage")
			clean = true

		} else {
			for s := range status {
				fmt.Println(s)
				_, err = w.Add(s)
				if err != nil {
					return fmt.Errorf("staging file in git: %w", err)
				}
			}
		}
		// if the user wants to commit the changes, do so
		if p.Config.Commit && !clean {
			fmt.Println("Committing changes")
			_, err = w.Commit("chore: Surgeon changes", &git.CommitOptions{})
			if err != nil {
				return fmt.Errorf("committing changes: %w", err)
			}
			// // if the user wants to push the changes, do so
			// if p.Config.Push {
			// 	fmt.Println("Pushing changes")
			// 	err = p.forkRepo.Push(&git.PushOptions{})
			// 	if err != nil {
			// 		return fmt.Errorf("pushing changes: %w", err)
			// 	}
			// }
		}
	}
	// clean up the temporary directory
	defer os.RemoveAll(p.UpsreamRoot)

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
	fmt.Println("Cloning upstream repository to", p.UpsreamRoot)
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
	fmt.Printf("Copying %s to %s\n", sourcePath, targetPath)
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
		if strings.HasPrefix(path, i.Prefix) {
			fmt.Println("Ignoring", path)
			return true
		}
	}
	return false
}
