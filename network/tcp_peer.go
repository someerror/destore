package network

import (
	"net"
	"sync"

	"github.com/someerror/destore/core"
)

type TCPPeer struct {
	net.Conn
	outbound bool // if true its client - out / if false its server - in
	wg       *sync.WaitGroup
	encoder  core.Encoder
}

var _ core.Peer = (*TCPPeer)(nil)

func NewTCPPeer(conn net.Conn, outbound bool, encoder core.Encoder) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
		outbound: outbound,
		wg:       &sync.WaitGroup{},
		encoder:  encoder,
	}
}

func (p *TCPPeer) Send(msg core.Message) error {
	return p.encoder.Encode(p.Conn, msg)
}
