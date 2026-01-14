package network

import "github.com/someerror/destore/core"

type HandshakeFunc func(core.Peer) error

// Mock
func NopHandshakeFunc(core.Peer) error {
	return nil
}

// TODO: need to realize handshake func
