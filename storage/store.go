package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/someerror/destore/core"
)

type StoreConf struct {
	Root          string
	PathGenerator PathGenerator
}

type Store struct {
	StoreConf
}

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

// Check implementation of Store interface
var _ core.Store = (*Store)(nil)

func NewStore(conf StoreConf) *Store {
	if len(conf.Root) == 0 {
		conf.Root = "destore_storage"
	}

	if conf.PathGenerator == nil {
		conf.PathGenerator = DefaultPathGenerator
	}

	return &Store{StoreConf: conf}
}

type ResolvedPath struct {
	PathName string
	FileName string
}

func (p ResolvedPath) FullPath() string {
	return filepath.Join(p.PathName, p.FileName)
}

func (p ResolvedPath) PathByIndex(pathIndex int) string {
	paths := strings.Split(p.PathName, string(filepath.Separator))

	pathLen := len(paths)

	if pathLen == 0 || pathIndex > pathLen {
		return ""
	}

	return paths[pathIndex]
}

// PathGenerator function defines the strategy for generating filepath.
type PathGenerator func(key string) ResolvedPath

// DefaultPathGenerator create path base on SHA256 hash.
//
// Example: hash a1b432... => a1/b4/32/.../file.svg
func DefaultPathGenerator(key string) ResolvedPath {
	hash := sha256.Sum256([]byte(key))
	hashString := hex.EncodeToString(hash[:])

	pathSegmentLength := 2
	pathDepth := len(hashString) / pathSegmentLength
	paths := make([]string, pathDepth)

	for i := range pathDepth {
		from := i * pathSegmentLength
		to := from + pathSegmentLength

		paths[i] = hashString[from:to]
	}

	return ResolvedPath{
		PathName: strings.Join(paths, string(filepath.Separator)),
		FileName: hashString,
	}
}
