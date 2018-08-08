NAME=ils
VERSION=$(shell git describe)

clean:
	rm -rf build/

build: clean
	mkdir -p build
	GOOS=darwin GOARCH=amd64 packr build -ldflags "-s -X main.version=$(VERSION) -X main.ThemePath=$(ILS_THEME_PATH)" -o build/$(NAME)-$(VERSION)
