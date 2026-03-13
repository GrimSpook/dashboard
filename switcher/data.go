package switcher

import (
	"crypto/rand"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
)

func generateUUID() string {
	b := make([]byte, 16)
	rand.Read(b)

	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

func FdSearch(dir string) ([]Workspace, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	baseDir := filepath.Join(home, dir)

	var dirs []string

	out, err := exec.Command("fd", "-Hs", "^.git$", "-td", "--prune", "--format", "{//}", baseDir).Output()
	if err != nil {
		return nil, err
	}

	dirs = strings.Split(string(out), "\n")
	dirs = dirs[:len(dirs)-1]

	var wsl []Workspace

	for _, workspaceDirs := range dirs {

		bn, err := branch(workspaceDirs)
		if err != nil {
			fmt.Println(workspaceDirs)
		}

		wsName := strings.ReplaceAll(filepath.Base(workspaceDirs), "-", " ")
		wsl = append(wsl, Workspace{Title: wsName, Path: workspaceDirs, Branch: bn, Id: generateUUID()})
	}

	return wsl, nil
}

func branch(path string) (string, error) {

	// out, err := exec.Command("cd ", path, "&&", "git", "--no-bager", "diff", "--stat").Output()

	repo, err := git.PlainOpen(path)
	if err != nil {
		return "", err
	}

	head, err := repo.Head()
	if err != nil {
		return "", err
	}

	return head.Name().Short(), nil
}
