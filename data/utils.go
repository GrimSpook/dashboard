package data

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

func fdSearch(dir string) ([]Workspace, error) {
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


func base(slice []string) []string {
	s := []string{}
	for _, str := range slice {
		s = append(s, filepath.Base(str))
	}
	return s
}

func removeDuplicates(slice []string) []string {
	dedupMap := make(map[string]bool)
	result := []string{}
	for _, value := range slice {
		if !dedupMap[value] {
			dedupMap[value] = true
			result = append(result, value)
		}
	}
	return result
}

func GetSelectedData(sections []Section, id string) Workspace {
	l := MergeSectionWorkspaces(sections)
	if len(l) != 0 {
		filter := Find(l, func(w Workspace) bool {
			return w.Id == id
		})
		return filter[0]
	}
	return Workspace{}
}

func Find[T any](items []T, predicate func(T) bool) []T {
	var result []T
	for _, item := range items {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

//git functions 
func branch(path string) (string, error) {
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


