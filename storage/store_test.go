package storage

import (
	"bytes"
	"fmt"
	"testing"
)

func TestDefaultPathGenerator(t *testing.T) {
	key := "awesometestfile"
	expectedFileName := "a8d78228c58bcedf8a9e9b8dda583ce8aecb60cdfebdfcb16571195e64a85df0"
	expectedPathName := "a8/d7/82/28/c5/8b/ce/df/8a/9e/9b/8d/da/58/3c/e8/ae/cb/60/cd/fe/bd/fc/b1/65/71/19/5e/64/a8/5d/f0"
	resolvedPath := DefaultPathGenerator(key)

	if resolvedPath.FileName != expectedFileName {
		t.Errorf("expected filename %s not equal with generated filename %s", expectedFileName, resolvedPath.FileName)
	}

	if resolvedPath.PathName != expectedPathName {
		t.Errorf("expected pathname %s not equal with generated pathname %s", expectedPathName, resolvedPath.PathName)
	}
}

func newStore(t *testing.T) *Store {
	conf := StoreConf{
		Root:          t.TempDir(),
		PathGenerator: DefaultPathGenerator,
	}

	return NewStore(conf)
}

func TestStore_Lyfecicle(t *testing.T) {
	store := newStore(t)

	key := "usr/someerror/avatar"
	data := []byte("data_bytes")

	// Stage 1: Write func
	n, err := store.Write(key, bytes.NewReader(data))

	if err != nil {
		t.Fatalf("Write Failed")
	}

	dataLen := int64(len(data))

	if n != dataLen {
		t.Errorf("Short write: expected %d bytes, got %d bytes", dataLen, n)
	}

	// Stage 2: Has func
	exists, err := store.Has(key)

	if err != nil {
		t.Fatalf("Has func Failed")
	}

	if !exists {
		t.Fatalf("Has func error: Expected key %s to exist", key)
	}

	// Stage 3: Read func
}

func TestStore_Parallel_Write(t *testing.T) {
	store := newStore(t)

	for i := range 100 {
		t.Run(fmt.Sprintf("peer_%d", i), func(t *testing.T) {
			t.Parallel()
			key := fmt.Sprintf("key_%d", i)
			data := []byte("data")

			_, err := store.Write(key, bytes.NewReader(data))

			if err != nil {
				t.Errorf("Write error in %d task: %s", i, err)
			}
		})
	}
}
