package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"path/filepath"
	"strings"
)

// Struct resolving by [PathGenerator]
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
