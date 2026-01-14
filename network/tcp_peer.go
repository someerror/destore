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
}

var _ core.Peer = (*TCPPeer)(nil)

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
		outbound: outbound,
		wg:       &sync.WaitGroup{},
	}
}

func (p *TCPPeer) Send(data []byte) error {
	// TODO optimization add buffio
	_, err := p.Conn.Write(data)

	return err
}
