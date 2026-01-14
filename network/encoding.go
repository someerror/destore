package network

import (
	"encoding/gob"
	"io"

	"github.com/someerror/destore/core"
)

type GOBDecoder struct{}

var _ core.Decoder = (*GOBDecoder)(nil)

func (d *GOBDecoder) Decode(r io.Reader, msg *core.Message) error {
	peekBuf := make([]byte, 1)
	if _, err := r.Read(peekBuf); err != nil {
		return err
	}

	// if its a raw data give control flow to server
	if peekBuf[0] == core.MessageTypeStream {
		msg.IsStream = true

		return nil
	}

	return gob.NewDecoder(r).Decode(msg)
}

type GOBEncoder struct {}

var _ core.Encoder = (*GOBEncoder)(nil)

func (e *GOBEncoder) Encode(w io.Writer, msg core.Message) error {
	var headerByte byte = core.MessageTypeCommand

	if (msg.IsStream) {
		headerByte = core.MessageTypeStream
	}

	_, err := w.Write([]byte{headerByte})
	if (err != nil) {
		return err
	}

	if (msg.IsStream) {
		return nil
	}

	return gob.NewEncoder(w).Encode(msg)
}


