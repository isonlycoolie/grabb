// Command grabb runs the download server and opens its dashboard.
package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/isonlycoolie/grabb/internal/desktop"
	"github.com/isonlycoolie/grabb/internal/download"
	"github.com/isonlycoolie/grabb/internal/httpapi"
	"github.com/isonlycoolie/grabb/internal/jobs"
	"github.com/isonlycoolie/grabb/internal/media"
	"github.com/isonlycoolie/grabb/internal/toolchain"
	"github.com/isonlycoolie/grabb/web"
)

// version is set at release time with -ldflags "-X main.version=...".
var version = "0.1.0"

func main() {
	if err := run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func run() error {
	home, _ := os.UserHomeDir()
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	addr := flag.String("addr", "127.0.0.1:7878", "address to serve the dashboard on")
	outDir := flag.String("out", filepath.Join(home, "Downloads", "grabb"), "folder for finished files")
	workers := flag.Int("workers", 2, "downloads that run at the same time")
	conns := flag.Int("connections", 16, "parallel connections per file")
	noOpen := flag.Bool("no-open", false, "do not open the dashboard in a browser")
	flag.Parse()

	tools, err := toolchain.Locate()
	if err != nil {
		return err
	}
	cfg := download.Config{OutDir: *outDir, Connections: max(*conns, 1), FFmpeg: tools.FFmpeg, Aria2: tools.Aria2}
	runner := func(ctx context.Context, link string, o download.Options, report func(download.Event)) (string, error) {
		return download.Run(ctx, tools.YtDlp, download.Args(cfg, o, link), report)
	}
	manager, err := jobs.NewManager(filepath.Join(configDir, "grabb", "library.json"), *workers, runner)
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	manager.Start(ctx)

	probe := func(ctx context.Context, link string) (media.Info, error) { return media.Probe(ctx, tools.YtDlp, link) }
	settings := httpapi.Settings{
		Version: version, OutDir: cfg.OutDir, Workers: max(*workers, 1),
		Connections: cfg.Connections, Tools: tools.Versions(ctx),
	}
	srv := &http.Server{
		Handler:           httpapi.New(manager, probe, desktop.Reveal, settings, web.FS()),
		ReadHeaderTimeout: 10 * time.Second,
		BaseContext:       func(net.Listener) context.Context { return ctx },
	}

	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		return err
	}
	url := "http://" + ln.Addr().String()
	slog.Info("grabb is running", "dashboard", url, "saving to", cfg.OutDir)
	if !*noOpen {
		if err := desktop.Open(url); err != nil {
			slog.Warn("could not open a browser", "error", err)
		}
	}

	go func() {
		<-ctx.Done()
		shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdown)
	}()
	if err := srv.Serve(ln); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
