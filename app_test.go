package main

import (
	"net"
	"testing"
)

func TestListenLocalUsesPreferredPort(t *testing.T) {
	port := freePort(t)

	ln, err := listenLocal(port)
	if err != nil {
		t.Fatalf("listenLocal(%d): %v", port, err)
	}
	defer ln.Close()

	got := ln.Addr().(*net.TCPAddr).Port
	if got != port {
		t.Fatalf("port = %d, want %d", got, port)
	}
}

func TestListenLocalFallsBackWhenPreferredPortBusy(t *testing.T) {
	blocker, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen blocker: %v", err)
	}
	defer blocker.Close()
	blockedPort := blocker.Addr().(*net.TCPAddr).Port

	ln, err := listenLocal(blockedPort)
	if err != nil {
		t.Fatalf("listenLocal(%d): %v", blockedPort, err)
	}
	defer ln.Close()

	got := ln.Addr().(*net.TCPAddr).Port
	if got == blockedPort {
		t.Fatalf("port = %d, want fallback", got)
	}
}

func freePort(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen free port: %v", err)
	}
	defer ln.Close()
	return ln.Addr().(*net.TCPAddr).Port
}

