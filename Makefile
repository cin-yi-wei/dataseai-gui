.PHONY: build build-mac build-win build-linux dev clean

# Build for current platform
build:
	wails build

# macOS universal binary (arm64 + amd64)
build-mac:
	wails build -platform darwin/universal -o DataseAI

# Windows
build-win:
	wails build -platform windows/amd64 -nsis

# Linux
build-linux:
	wails build -platform linux/amd64

# Dev mode (hot reload)
dev:
	wails dev

clean:
	rm -rf build/bin
