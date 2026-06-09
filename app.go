package main

import (
	"context"
	"io/fs"
	"log"
	"net"
	"net/http"
	goruntime "runtime"
	"os"
	"path/filepath"
	"time"

	dataseai "github.com/conray/dataseai"
)

type App struct {
	ctx    context.Context
	port   int
	srv    *http.Server
	dsServ *dataseai.Server
}

func NewApp() *App { return &App{} }

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	dataDir, err := appDataDir()
	if err != nil {
		log.Fatalf("data dir: %v", err)
	}

	webSub, err := fs.Sub(dataseai.WebFS, "web/dist")
	if err != nil {
		log.Fatalf("embed sub: %v", err)
	}

	a.dsServ, err = dataseai.NewServer(dataseai.ServerConfig{
		DBPath:       filepath.Join(dataDir, "dataseai.db"),
		KeyPath:      filepath.Join(dataDir, "master.key"),
		Version:      "gui-1.0.0",
		Registration: "open",
		WebFS:        webSub,
	})
	if err != nil {
		log.Fatalf("server init: %v", err)
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	a.port = ln.Addr().(*net.TCPAddr).Port
	log.Printf("dataseai-gui on :%d (data: %s)", a.port, dataDir)

	a.srv = &http.Server{Handler: a.dsServ.Handler()}
	go func() {
		if err := a.srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			log.Printf("server: %v", err)
		}
	}()
}

func (a *App) shutdown(ctx context.Context) {
	if a.srv != nil {
		shutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		_ = a.srv.Shutdown(shutCtx)
	}
	if a.dsServ != nil {
		a.dsServ.Close()
	}
}

func appDataDir() (string, error) {
	var base string
	switch goruntime.GOOS {
	case "darwin":
		home, _ := os.UserHomeDir()
		base = filepath.Join(home, "Library", "Application Support", "DataseAI")
	case "windows":
		base = filepath.Join(os.Getenv("APPDATA"), "DataseAI")
	default:
		if xdg := os.Getenv("XDG_DATA_HOME"); xdg != "" {
			base = filepath.Join(xdg, "dataseai")
		} else {
			home, _ := os.UserHomeDir()
			base = filepath.Join(home, ".local", "share", "dataseai")
		}
	}
	return base, os.MkdirAll(base, 0700)
}
