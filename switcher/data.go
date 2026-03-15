package switcher

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/sahilm/fuzzy"
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

func GenerateSections() []Section {
	openWorkspaceList, err := getOpenWorkspaces()
	configWorkspacesList, err := getConfigDirs()
	personalWorkspaceList, err := FdSearch(filepath.Join("dev", "personal"))
	schoolWorkspaceList, err := FdSearch(filepath.Join("dev", "school"))
	// zoxideList := getZoxidPaths()
	if err != nil {
		log.Fatal(err)
	}

	return []Section{
		{
			Title: "Open",
			List:  openWorkspaceList,
		},
		{
			Title: "Configs",
			List:  configWorkspacesList,
		},
		{
			Title: "Personal",
			List:  personalWorkspaceList,
		},
		{
			Title: "School",
			List:  schoolWorkspaceList,
		},
		// {
		// 	Title: "Zoxide",
		// 	List:  zoxideList,
		// },
	}
}

func getConfigDirs() ([]Workspace, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	weztermDir := filepath.Join(home, ".config", "wezterm")
	nvimDir := filepath.Join(home, "AppData", "Local", "nvim")

	dirs := []string{weztermDir, nvimDir}

	var wsl []Workspace

	for _, configDirs := range dirs {
		wsName := strings.ReplaceAll(filepath.Base(configDirs), "-", " ")
		wsl = append(wsl, Workspace{Title: wsName, Path: configDirs, Id: generateUUID()})
	}

	return wsl, nil
}

func mergeSectionWorkspaces(sections []Section) []Workspace {
	var workspaces []Workspace
	for _, section := range sections {
		workspaces = append(workspaces, section.List...)
	}

	return workspaces
}

func getOpenWorkspaces() ([]Workspace, error) {
	args := []string{"cli", "list", "--format", "json"}

	out, err := exec.Command("wezterm", args...).Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute command: %w", err)
	}

	var temp []weztermCliJson

	if err := json.Unmarshal(out, &temp); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	var wsl []Workspace
	var prev string
	for _, value := range temp {

		wsName := strings.ReplaceAll(filepath.Base(value.Workspace), "-", " ")

		if wsName == "default" {

			if prev != wsName {
				wsl = append(wsl, Workspace{Title: wsName, Path: home, Id: generateUUID()})
			}

		} else {
			if prev != wsName {
				wsl = append(wsl, Workspace{Title: wsName, Path: value.Workspace, Id: generateUUID()})
			}
		}

		prev = wsName
	}

	return wsl, nil
}

func (m Model) Filter(query string, workspaces []Workspace) []Workspace {
	if query == "" {
		return workspaces
	}

	titles := make([]string, len(workspaces))
	for i, w := range workspaces {
		titles[i] = w.Title
	}

	matches := fuzzy.Find(query, titles)

	wsMap := make(map[string]Workspace)
	for _, w := range workspaces {
		wsMap[w.Title] = w
	}

	filtered := make([]Workspace, 0, len(matches))
	for _, match := range matches {
		filtered = append(filtered, wsMap[match.Str])
	}

	return filtered
}

func getZoxidPaths() []Workspace {
	out, err := exec.Command("zoxide", "query", "-l").Output()
	if err != nil {
		log.Fatal(err)
	}

	strs := strings.Split(string(out), "\n")

	basestr := base(strs)
	newStrs := removeDuplicates(basestr)

	var ws []Workspace

	for _, str := range newStrs {
		ws = append(ws, Workspace{Title: filepath.Base(str), Path: str, Id: generateUUID()})
	}

	return ws
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
