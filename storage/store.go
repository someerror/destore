package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"path/filepath"
	"strings"
)

type PathKey struct {
	PathName string
	FileName string
}

func (p PathKey) FullPath() string {
	return filepath.Join(p.PathName, p.FileName)
}

func (p PathKey) PathByIndex(pathIndex int) string {
	paths := strings.Split(p.PathName, string(filepath.Separator))

	pathLen := len(paths)

	if pathLen == 0 || pathIndex > pathLen {
		return ""
	}

	return paths[pathIndex]
}

type PathGenerator func(key string) PathKey

type StoreConf struct {
	Root          string
	PathGenerator PathGenerator
}

type Store struct {
	StoreConf
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

/*
DefaultPathGenerator create path base on SHA256 hash.
Example a1/b4/32/.../file.svg
*/
func DefaultPathGenerator(key string) PathKey {
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

	return PathKey{
		PathName: strings.Join(paths, string(filepath.Separator)),
		FileName: hashString,
	}
}
