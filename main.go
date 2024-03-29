package main

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/pion/mdns"
	"golang.org/x/net/ipv4"
)

func main() {
	err := setupMDNS([]string{"talos-config.local"})
	if err != nil {
		slog.Error("Failed to initialize mdns", "error", err)
	}
	err = setupFileServer("./configs")
	if err != nil {
		slog.Error("Failed to start http server", "error", err)
	}
	select {}
}

func setupMDNS(localnames []string) error {
	addr, err := net.ResolveUDPAddr("udp", mdns.DefaultAddress)
	if err != nil {
		return fmt.Errorf("Failed to acquire UDP addr: %w", err)
	}

	slog.Info("Broadcasting:", "ip", addr.IP)

	l, err := net.ListenUDP("udp4", addr)
	if err != nil {
		return fmt.Errorf("Failed to setup UDP listener: %w", err)
	}

	_, err = mdns.Server(ipv4.NewPacketConn(l), &mdns.Config{
		LocalNames: []string{"talos-config.local"},
	})
	if err != nil {
		return fmt.Errorf("Failed to start mdns listener: %w", err)

	}
	return nil
}

func setupFileServer(configdir string) error {
	fs := http.FileServer(http.Dir(configdir))
	http.Handle("/", fs)

	err := http.ListenAndServe(":8423", nil)
	if err != nil {
		return fmt.Errorf("failed to start http listener: %w", err)
	}
	return nil
}
