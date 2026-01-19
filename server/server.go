package server

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/someerror/destore/core"
)

type ServerOpts struct {
	ListenAddr string
	Store      core.Store
	Transport  core.Transport
	Logger     *slog.Logger
}

type Server struct {
	ServerOpts
}

func NewServer(opts ServerOpts) (*Server, error) {
	if opts.Transport == nil {
		return nil, fmt.Errorf("server options: transport is missing")
	}
	if opts.Store == nil {
		return nil, fmt.Errorf("server options: store is missing")
	}

	return &Server{
		ServerOpts: opts,
	}, nil
}

func (srv *Server) Start(ctx context.Context) error {
	srv.Logger.Info("Server starting on", "addr", srv.ListenAddr)

	err := srv.Transport.ListenAndAccept()
	if err != nil {
		return err
	}

	srv.loop(ctx)

	srv.Logger.Info("Server stopped")

	return nil
}

func (srv *Server) loop(ctx context.Context) {
	defer func() {
		srv.Logger.Info("Shutting down transport")
		srv.Transport.Close()
	}()

	select {
	case msg := <-srv.Transport.Consume():
		err := srv.handleMessage(msg)

		if err != nil {
			srv.Logger.Error("Handling message", "error", err)

			response := core.Message{
				Payload: []byte(fmt.Sprintf("Server error: %q", err)),
			}

			if msg.SourcePeer != nil {
				msg.SourcePeer.Send(response)
			}
		}
	case <-ctx.Done():
		srv.Logger.Info("Context done")
		return
	}
}

func (srv *Server) handleMessage(msg core.Message) error {
	srv.Logger.Info("message handling")

	if len(msg.Payload) == 0 {
		return fmt.Errorf("empty payload received")
	}

	return nil
}
