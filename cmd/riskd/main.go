package main

import (
	"context"
	"example.com/sbom-risk-graph/internal/observability"
	"example.com/sbom-risk-graph/internal/repository"
	httpapi "example.com/sbom-risk-graph/internal/transport/http"
	"example.com/sbom-risk-graph/pkg/logging"
	"example.com/sbom-risk-graph/pkg/metrics"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func shutdownContext(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	return ctx, cancel
}

func main() {
	addr := flag.String("addr", ":8123", "listen address")
	dir := flag.String("data-dir", "./data", "data directory")
	flag.Parse()
	logger := logging.New("sbom-risk-graph")
	repo, err := repository.Open(*dir)
	if err != nil {
		logger.Error("open", logging.Err(err))
		os.Exit(1)
	}
	service := repository.NewService(repo)
	counters := &metrics.Counters{}
	server := &http.Server{Addr: *addr, Handler: httpapi.New(service, counters).Routes(), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 20 * time.Second, WriteTimeout: 20 * time.Second}
	go func() {
		logger.Info("started", "addr", *addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("serve", logging.Err(err))
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	ctx, cancel := observability.ShutdownContext(context.Background(), 10*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
}
