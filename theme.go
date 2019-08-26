package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gobuffalo/packr"
)

const defaultIcon = "default"
const filename = "theme.json"

// iconsCache is a specified type for memo caching loaded icons
type iconsCache map[string]string

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

// Theme is unite types, folders and files icons
type Theme struct {
	extensions, folders, files *iconsLoader
}

// NewTheme is a theme constructor
func NewTheme(location string) (*Theme, error) {
	box := packr.NewBox(location)
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
		extensions: newIconsLoader(location, data.Extensions),
		folders:    newIconsLoader(location, data.Folders),
		files:      newIconsLoader(location, data.Files),
	}

	// Initialize default icons and returns error if it's not found
	for _, l := range [...]*iconsLoader{t.extensions, t.folders, t.files} {
		if _, err := l.Load(defaultIcon); err != nil {
			return nil, err
		}
	}

	return &t, nil
}

// GetIcon is match and returns file icon
func (t *Theme) GetIcon(file os.FileInfo) (icon string) {
	var err error

	if file.IsDir() {
		icon, err = t.folders.Load(file.Name())
		if err != nil {
			icon, _ = t.folders.Load(defaultIcon)
		}
		return icon
	}

	icon, err = t.files.Load(strings.ToLower(file.Name()))
	if err == nil {
		return icon
	}

	if strings.Contains(file.Name(), ".") {
		sl := strings.Split(file.Name(), ".")
		extension := sl[len(sl)-1]
		icon, err = t.extensions.Load(strings.ToLower(extension))
		if err == nil {
			return icon
		}
	}

	icon, _ = t.files.Load(defaultIcon)
	return icon
}

type iconsLoader struct {
	location string
	icons    map[string]Icon
	cache    iconsCache
}

func newIconsLoader(location string, icons []Icon) *iconsLoader {
	return &iconsLoader{
		icons:    hashIcons(icons),
		cache:    make(iconsCache),
		location: location,
	}
}

// LoadIcon returns encoded icon in base64 from cache or file
func (l iconsLoader) Load(filename string) (string, error) {
	// Try to get icon from cache
	if result, ok := l.cache[filename]; ok {
		return result, nil
	}
	// Match and load icon if it's not in cache
	icon, ok := l.icons[filename]
	if !ok {
		// Return error if icon not matched
		return "", fmt.Errorf("Icon '%s' not matched", filename)
	}
	// Load icon in base62 encoding
	// Used default icon if file can't be loaded
	result, err := icon.Load(l.location)
	if err != nil {
		return "", err
	}
	// Set value in cache and return icon
	l.cache[filename] = result
	return result, nil
}

func hashIcons(icons []Icon) map[string]Icon {
	result := make(map[string]Icon)
	for _, icon := range icons {
		for _, name := range icon.Names {
			result[name] = icon
		}
	}
	return result
}
