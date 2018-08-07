package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gobuffalo/packr"
)

const (
	// Default is a default icon name
	Default = "default"
	// Filename is a default theme configuration filename
	filename = "theme.json"
)

// cache is a specified type for memo caching loaded icons
type cache map[string]string

func hashIcons(icons []Icon) map[string]Icon {
	result := make(map[string]Icon)
	for _, icon := range icons {
		for _, name := range icon.Names {
			result[name] = icon
		}
	}
	return result
}

// NewTheme is a theme constructor
func NewTheme(path string) (*Theme, error) {
	// Read specified theme file
	box := packr.NewBox(path)
	file, err := box.MustBytes(filename)
	if err != nil {
		return nil, err
	}

	data := struct {
		Extensions []Icon `json:"extensions"`
		Folders    []Icon `json:"folders"`
		Files      []Icon `json:"files"`
	}{}

	json.Unmarshal(file, &data)

	// Create new Theme with types, folders and files icons lists
	t := Theme{
		Extensions: &IconList{
			Icons: hashIcons(data.Extensions),
			path:  path,
			c:     make(cache),
		},
		Folders: &IconList{
			Icons: hashIcons(data.Folders),
			path:  path,
			c:     make(cache),
		},
		Files: &IconList{
			Icons: hashIcons(data.Files),
			path:  path,
			c:     make(cache),
		},
	}

	// Initialize default icons or return error if it's not found
	// Initialize default types icon
	for _, list := range [...]*IconList{t.Extensions, t.Folders, t.Files} {
		if _, err := list.GetIcon(Default); err != nil {
			return nil, err
		}
	}

	return &t, nil
}

// Theme is unite types, folders and files icons
type Theme struct {
	Extensions, Folders, Files *IconList
}

// GetIcon is match and returns file icon
func (t *Theme) GetIcon(file os.FileInfo) string {
	var icon string
	var err error

	if file.IsDir() {
		icon, err = t.Folders.GetIcon(file.Name())
		if err != nil {
			icon, _ = t.Folders.GetIcon(Default)
		}
		return icon
	}

	icon, err = t.Files.GetIcon(strings.ToLower(file.Name()))
	if err != nil {
		if strings.Contains(file.Name(), ".") {
			sl := strings.Split(file.Name(), ".")
			extension := sl[len(sl)-1]
			icon, err = t.Extensions.GetIcon(strings.ToLower(extension))
			if err == nil {
				return icon
			}
		}
		icon, _ = t.Files.GetIcon(Default)
	}

	return icon
}

// Icon structure with name and filename properties
type Icon struct {
	Names    []string `json:"names"`
	Filename string   `json:"filename"`
}

// Load and encode icon in base64 from a file
func (i Icon) Load(path string) (string, error) {
	box := packr.NewBox(path)

	data, err := box.MustBytes(i.Filename)
	if err != nil {
		return "", err
	}

	b64data := base64.StdEncoding.EncodeToString(data)
	return b64data, nil
}

// IconList base structure for all icon types
type IconList struct {
	Icons map[string]Icon

	path string
	c    cache
}

// GetIcon returns encoded icon in base64 from cache or file
func (l IconList) GetIcon(filename string) (string, error) {
	// Try to get icon from cache
	icon, ok := l.c[filename]
	if !ok {
		// Match and load icon if it's not in cache
		item, ok := l.Icons[filename]
		if !ok {
			// Return error if icon not matched
			return "", fmt.Errorf("Icon '%s' not matched", filename)
		}

		// Load icon in base62 encoding
		// Used default icon if file can't be loaded
		icon, err := item.Load(l.path)
		if err != nil {
			return "", err
		}

		// Set value in cache and return icon
		l.c[filename] = icon
		return icon, nil
	}

	return icon, nil
}