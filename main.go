package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/someerror/destore/network"
	"github.com/someerror/destore/server"
	"github.com/someerror/destore/storage"
)

func main() {
	tcpOpts := network.TCPTransportOpts{
		ListenAddr:    ":4000",
		Decoder:       network.NewGOBDecoder(),
		Encoder:       network.NewGOBEncoder(),
		HandshakeFunc: network.NopHandshakeFunc,
	}

	tcp := network.NewTCPTransport(tcpOpts)
	store := storage.NewStore(storage.StoreConf{})
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	srv, err := server.NewServer(server.ServerOpts{
		Store:      store,
		Transport:  tcp,
		ListenAddr: ":4000",
		Logger:     logger,
	})
	if err != nil {
		slog.Error("Failed to create server", "err", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err = srv.Start(ctx); err != nil {
		slog.Error("Server exited with error", "err", err)
		os.Exit(1)
	}

	slog.Info("Shutting down..")
}
