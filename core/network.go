package core

import (
	"io"
	"net"
)

type Message struct {
	From       string
	Payload    []byte
	IsStream   bool
	SourcePeer Peer
}

type Peer interface {
	net.Conn
	Send(Message) error
}

type Transport interface {
	Addr() string           // Addr returns the listening address
	Dial(addr string) error // Dial	creates a new outgoing connection
	ListenAndAccept() error
	Consume() <-chan Message
	Close() error
}

const (
	MessageTypeCommand = 0x01
	MessageTypeStream  = 0x02
)

type Decoder interface {
	Decode(io.Reader, *Message) error
}

type Encoder interface {
	Encode(io.Writer, Message) error
}
