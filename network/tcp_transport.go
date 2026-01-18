package network

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"

	"github.com/someerror/destore/core"
)

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       core.Decoder
	Encoder       core.Encoder
	OnPeer        func(core.Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	rpcch    chan core.Message
}

var _ core.Transport = (*TCPTransport)(nil)

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch:            make(chan core.Message, 1024), // buffer so that the server has time to read RPC packs
	}
}

func (t *TCPTransport) Addr() string {
	return t.ListenAddr
}

// Consume returns read-only RPC chan for incoming messages
func (t *TCPTransport) Consume() <-chan core.Message {
	return t.rpcch
}

func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

// Create connection to Peer
func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	go t.handleConn(conn, true)

	return nil
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	slog.Info("TCP transport started", "port", t.ListenAddr)

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	// TODO
	for {
		conn, err := t.listener.Accept()

		// t.Close
		if errors.Is(err, net.ErrClosed) {
			return
		}
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
			continue
		}

		go t.handleConn(conn, false)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) {
	var err error

	defer func() {
		if !errors.Is(err, io.EOF) {
			fmt.Printf("peer connection dropped with error: %q", err)
		}
		conn.Close()
	}()

	peer := NewTCPPeer(conn, outbound, t.Encoder)

	if t.HandshakeFunc != nil {
		err = t.HandshakeFunc(peer)
		if err != nil {
			return
		}
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}

	for {
		msg := core.Message{}

		err = t.Decoder.Decode(conn, &msg)
		if err != nil {
			return
		}

		msg.From = conn.RemoteAddr().String()
		msg.SourcePeer = peer

		t.rpcch <- msg
	}
}
