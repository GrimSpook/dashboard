package data

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Section struct {
	Title string
	Icon  string
	List  []Workspace
}

type Workspace struct {
	Title  string
	Path   string
	Branch string
	Id     string
	Status string
}

type weztermCliJson struct {
	Window_id float64 `json:"window_id"`
	Tab_id    float64 `json:"tab_id"`
	Pane_id   float64 `json:"pane_id"`
	Workspace string  `json:"workspace"`
	// Size              map[string]size `json:"size"`
	Title             string  `json:"title"`
	Cwd               string  `json:"cwd"`
	Cursor_x          float64 `json:"cursor_x"`
	Cursor_y          float64 `json:"cursor_y"`
	Cursor_shape      string  `json:"cursor_shape"`
	Cursor_visibility string  `json:"cursor_visibility"`
	Left_col          float64 `json:"left_col"`
	Top_row           float64 `json:"top_row"`
	Tab_title         string  `json:"tab_title"`
	Window_title      string  `json:"window_title"`
	Is_active         bool    `json:"is_active"`
	Is_zoomed         bool    `json:"is_zoomed"`
	Tty_name          string  `json:"tty_name"`
}

func GenerateSections() []Section {
	openWorkspaceList, err := getOpenWorkspaces()
	configWorkspacesList, err := getConfigDirs()
	personalWorkspaceList, err := fdSearch(filepath.Join("dev", "personal"))
	schoolWorkspaceList, err := fdSearch(filepath.Join("dev", "school"))
	// zoxideList := getZoxidPaths()
	if err != nil {
		log.Fatal(err)
	}

	return []Section{
		{
			Title: "Open",
			Icon:  "",
			List:  openWorkspaceList,
		},
		{
			Title: "Configs",
			Icon:  "",
			List:  configWorkspacesList,
		},
		{
			Title: "Personal",
			Icon:  "",
			List:  personalWorkspaceList,
		},
		{
			Title: "School",
			Icon:  "",
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

func MergeSectionWorkspaces(sections []Section) []Workspace {
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
