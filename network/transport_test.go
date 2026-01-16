package network

import (
	"net"
	"testing"
	"time"

	"github.com/someerror/destore/core"
)

const DEST_PORT = ":4000"

func TestTCPTrans(t *testing.T) {
	opts := TCPTransportOpts{
		ListenAddr:    DEST_PORT,
		HandshakeFunc: NopHandshakeFunc,
		Decoder:       &GOBDecoder{},
	}

	transport := NewTCPTransport(opts)
	defer transport.Close()

	err := transport.ListenAndAccept()
	if err != nil {
		t.Fatalf("tcp transport error: %q", err)
	}

	time.Sleep(2 * time.Second)

	t.Run("transport:send_message", func(t *testing.T) {
		conn, err := net.Dial("tcp", DEST_PORT)
		if err != nil {
			t.Fatalf("peer connection error")
		}
		defer conn.Close()

		msg := core.Message{
			IsStream: false,
			Payload:  []byte("Command Payload"),
		}

		enc := &GOBEncoder{}
		err = enc.Encode(conn, msg)
		if err != nil {
			t.Fatalf("Encoding failed: %q", err)
		}

		select {
		case msg := <-transport.Consume():
			if msg.IsStream {
				t.Errorf("expected msg with IsStream=false, got IsStream=true")
			}

			if string(msg.Payload) != "Command Payload" {
				t.Errorf("expected msg.Payload 'Command Payload', got %q", msg.Payload)
			}
		case <-time.After(2 * time.Second):
			t.Errorf("timeout waiting for regular message")
		}

	})

	t.Run("transport:send_stream", func(t *testing.T) {
		conn, err := net.Dial("tcp", DEST_PORT)
		if err != nil {
			t.Fatalf("peer connection error")
		}
		defer conn.Close()

		msg := core.Message{
			IsStream: true,
		}

		enc := &GOBEncoder{}
		err = enc.Encode(conn, msg)
		if err != nil {
			t.Fatalf("Encoding failed: %q", err)
		}

		_, err = conn.Write([]byte("ROW_BINARY_DATA_OF_THE_FILE"))
		if err != nil {
			t.Fatalf("binary data writing failed with stream, err: %q", err)
		}

		select {
		case msg := <-transport.Consume():
			if !msg.IsStream {
				t.Errorf("expected msg with IsStream=true, got IsStream=false")
			}

			if msg.SourcePeer == nil {
				t.Errorf("expected SourcePeer is not nil")
			}
			t.Log("Stream header received successfully")
		case <-time.After(2 * time.Second):
			t.Fatal("timeout waiting for stream message")
		}

	})
}
