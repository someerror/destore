package storage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/someerror/destore/core"
)

type StoreConf struct {
	Root          string
	PathGenerator PathGenerator
}

type Store struct {
	StoreConf
}

// Check implementation of Store interface
var _ core.Store = (*Store)(nil)

// os.Stat Syscall
func (s *Store) Has(key string) (bool, error) {
	resolvedPath := s.PathGenerator(key)
	fullPath := filepath.Join(s.Root, resolvedPath.FullPath())

	_, err := os.Stat(fullPath)

	// Case 1: file exists
	if err == nil {
		return true, nil
	}

	// Case 2: file not exists
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	// Case3: system error (permission or something else)
	return false, err
}

// os.RemoveAll Syscall
func (s *Store) Clear() error {
	if len(s.Root) < 3 || filepath.IsAbs(s.Root) {
		return fmt.Errorf("store.clear: unsafe root path for clear: %q", s.Root)
	}

	return os.RemoveAll(s.Root)
}

func NewStore(conf StoreConf) *Store {
	if len(conf.Root) == 0 {
		conf.Root = "destore_storage"
	}

	if conf.PathGenerator == nil {
		conf.PathGenerator = DefaultPathGenerator
	}

	return &Store{StoreConf: conf}
}
