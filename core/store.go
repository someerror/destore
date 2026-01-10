package core

import "io"

type Store interface {
	Has(key string) (bool, error)
	Clear() error
	Write(key string, r io.Reader) (int64, error)
	Read(key string) (io.ReadCloser, error)
	Deete(key string) error
}
