.PHONY: build build-mac build-win build-linux dev setup setup-linux clean

DATASEAI := ../dataseai

# 確保 dataseai 前端已 build（embed 需要 web/dist）
frontend:
	cd $(DATASEAI)/web && npm ci && npm run build

# 當前平台
build: frontend
	wails build

# macOS universal binary (Apple Silicon + Intel)
build-mac: frontend
	wails build -platform darwin/universal -o DataseAI

# Windows
build-win: frontend
	wails build -platform windows/amd64

# Linux (需要 libwebkit2gtk-4.0-dev)
build-linux: frontend
	wails build -platform linux/amd64

# 安裝 Linux 開發依賴（Ubuntu/Debian）
setup-linux:
	sudo apt-get update && sudo apt-get install -y --no-install-recommends \
		libwebkit2gtk-4.0-dev libgtk-3-dev pkg-config

# 安裝 Go 工具依賴
setup:
	go install github.com/wailsapp/wails/v2/cmd/wails@v2.12.0

# 開發模式（hot reload loading page；app 直接跑 HTTP server）
dev:
	wails dev

clean:
	rm -rf build/bin DataseAI DataseAI.exe
