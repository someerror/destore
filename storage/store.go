package storage

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/someerror/destore/core"
)

func NewStore(conf StoreConf) *Store {
	if len(conf.Root) == 0 {
		conf.Root = "destore_storage"
	}

	if conf.PathGenerator == nil {
		conf.PathGenerator = DefaultPathGenerator
	}

	return &Store{StoreConf: conf}
}

type StoreConf struct {
	Root          string
	PathGenerator PathGenerator
}

type Store struct {
	StoreConf
}

var _ core.Store = (*Store)(nil)

func (s *Store) resolvePath(key string) string {
	resolvedPath := s.PathGenerator(key)
	return filepath.Join(s.Root, resolvedPath.FullPath())
}

// os.Stat Syscall
func (s *Store) Has(key string) (bool, error) {
	fullPath := s.resolvePath(key)

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

// os.Open Syscall returns io.ReadCloser that must be closed mannualy
func (s *Store) Read(key string) (io.ReadCloser, error) {
	fullPath := s.resolvePath(key)

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("store.read:  open file error for key%q:%w", key, err)
	}

	return file, nil
}

// io.Copy
func (s *Store) Write(key string, r io.Reader) (int64, error) {
	f, err := s.createFile(key)
	if err != nil {
		return 0, err
	}

	defer f.Close()

	// TODO: this is basic implementation, for prod we need realize write to temp pattern and rename it after
	return io.Copy(f, r)
}

// os.MkdirAll Syscall os.Create Syscall
func (s *Store) createFile(key string) (*os.File, error) {
	fullPath := s.resolvePath(key)
	dir := filepath.Dir(fullPath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf(
			"store.createFile: failed to create path struct by os.MkdirAll, %w", err,
		)
	}

	return os.Create(fullPath)
}

// clear empty dirs
func (s *Store) cleanUpEmptyDirs(key string) {
	fullPath := s.resolvePath(key)
	// slice out filename
	dir := filepath.Clean(filepath.Dir(fullPath)) 

	for {
		if dir == s.Root || dir == "." {
			return
		}

		if err := os.Remove(dir); err != nil {
			return
		}

		dir = filepath.Dir(dir)
	}
}

// os.Remove Syscall and cleanup empty dirs
func (s *Store) Delete(key string) error {
	fullPath := s.resolvePath(key)

	if err := os.Remove(fullPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return fmt.Errorf("store.delete: os.Remove failed to remove for key %q: %w", key, err)
	}

	s.cleanUpEmptyDirs(key)

	return nil
}
