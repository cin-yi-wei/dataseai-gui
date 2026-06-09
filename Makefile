.PHONY: build build-mac build-win build-linux dev setup clean

MYSQLWEB := ../mysqlweb

# 確保 mysqlweb 前端已 build（embed 需要 web/dist）
frontend:
	cd $(MYSQLWEB)/web && npm ci && npm run build

# 當前平台
build: frontend
	wails build

# macOS universal binary (Apple Silicon + Intel)
build-mac: frontend
	wails build -platform darwin/universal -o DataseAI

# Windows
build-win: frontend
	wails build -platform windows/amd64 -nsis

# Linux (需要 libwebkit2gtk)
build-linux: frontend
	wails build -platform linux/amd64

# 開發模式（hot reload loading page；app 直接跑 HTTP server）
dev:
	wails dev

clean:
	rm -rf build/bin DataseAI DataseAI.exe
