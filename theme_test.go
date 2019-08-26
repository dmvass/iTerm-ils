package main

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const location = "theme"
const wrongLocation = "wrong_theme_path"

// FileInfo base mocked object
type FileMocked struct{}

func (fm *FileMocked) Name() string       { return "default" }
func (fm *FileMocked) IsDir() bool        { return false }
func (fm *FileMocked) Size() int64        { return 0 }
func (fm *FileMocked) Mode() os.FileMode  { return os.FileMode(0700) }
func (fm *FileMocked) ModTime() time.Time { return time.Time{} }
func (fm *FileMocked) Sys() interface{}   { return nil }

// FileInfo extension mocked objects
type FileMockedExtension struct {
	FileMocked
}

func (fm *FileMockedExtension) Name() string { return "default.go" }

// FileInfo file mocked objects
type FileMockedFile struct {
	FileMocked
}

type FileMockedWrongFile struct {
	FileMocked
}

func (fm *FileMockedWrongFile) Name() string { return "wrong_file_name" }

// FileInfo folder mocked objects
type FileMockedFolder struct {
	FileMocked
}

func (fm *FileMockedFolder) IsDir() bool { return true }

type FileMockedWrongFolder struct {
	FileMockedFolder
}

func (fm *FileMockedWrongFolder) Name() string { return "wrong_folder_name" }

func TestNewTheme(t *testing.T) {
	_, err := NewTheme(location)
	assert.Nil(t, err)

	_, err = NewTheme(wrongLocation)
	assert.NotNil(t, err)
}

func TestHashIcons(t *testing.T) {
	actual := []Icon{
		Icon{Names: []string{"test1", "test2"}, Filename: "first"},
		Icon{Names: []string{"test3", "test4"}, Filename: "second"},
	}
	expected := map[string]Icon{
		"test1": actual[0],
		"test2": actual[0],
		"test3": actual[1],
		"test4": actual[1],
	}
	assert.Equal(t, expected, hashIcons(actual))
}

func TestGetIcon(t *testing.T) {
	theme, err := NewTheme(location)
	assert.Nil(t, err)

	tables := []struct {
		actual   os.FileInfo
		expected Icon
	}{
		{new(FileMockedExtension), Icon{[]string{defaultIcon}, "extension_go.png"}},
		{new(FileMockedFile), Icon{[]string{defaultIcon}, "file_default.png"}},
		{new(FileMockedWrongFile), Icon{[]string{defaultIcon}, "file_default.png"}},
		{new(FileMockedFolder), Icon{[]string{defaultIcon}, "folder_default.png"}},
		{new(FileMockedWrongFolder), Icon{[]string{defaultIcon}, "folder_default.png"}},
	}

	for _, table := range tables {
		icon, err := table.expected.Load(location)
		assert.Nil(t, err)
		assert.Equal(t, icon, theme.GetIcon(table.actual))
	}
}

func TestLoad(t *testing.T) {
	icon := Icon{Filename: "file_default.png"}

	_, err := icon.Load(location)
	assert.Nil(t, err)

	_, err = icon.Load(wrongLocation)
	assert.NotNil(t, err)
}

func TestLoadIcon(t *testing.T) {
	loader := makeIconsLoader()

	icon, err := loader.Load("test1")
	assert.Nil(t, err)

	actul, ok := loader.cache["test1"]
	assert.True(t, ok)
	assert.Equal(t, icon, actul)
}

func TestLoadIconError(t *testing.T) {
	loader := makeIconsLoader()

	_, err := loader.Load("test2")
	assert.NotNil(t, err)

	_, err = loader.Load("wrong")
	assert.NotNil(t, err)
}

func TestLoadIconFromCache(t *testing.T) {
	expected := "expected"

	loader := makeIconsLoader()
	loader.cache["test1"] = expected

	icon, err := loader.Load("test1")
	assert.Nil(t, err)
	assert.Equal(t, expected, icon)
}

func makeIconsLoader() *iconsLoader {
	icons := []Icon{
		Icon{Names: []string{"test1"}, Filename: "file_default.png"},
		Icon{Names: []string{"test2"}, Filename: "wrong_file_path.png"},
	}
	return newIconsLoader(location, icons)
}
