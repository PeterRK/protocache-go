package protocache_test

import (
	"testing"

	"github.com/peterrk/protocache-go"
)

type externalHashKeySource struct {
	keys [][]byte
	curr int
}

func (s *externalHashKeySource) Reset() {
	s.curr = 0
}

func (s *externalHashKeySource) Total() int {
	return len(s.keys)
}

func (s *externalHashKeySource) Next() []byte {
	key := s.keys[s.curr]
	s.curr++
	return key
}

func TestPerfectHashTableExportedAPI(t *testing.T) {
	keys := [][]byte{
		[]byte("alpha"),
		[]byte("bravo"),
		[]byte("charlie"),
		[]byte("delta"),
	}

	table := protocache.BuildPerfectHashTable(&externalHashKeySource{keys: keys})
	if !table.IsValid() {
		t.Fatal("perfect hash table is invalid")
	}
	if table.Size() != uint32(len(keys)) {
		t.Fatalf("size = %d, want %d", table.Size(), len(keys))
	}

	seen := make([]bool, len(keys))
	for _, key := range keys {
		slot := table.Lookup(key)
		if slot >= uint32(len(keys)) {
			t.Fatalf("slot = %d, want < %d", slot, len(keys))
		}
		if seen[slot] {
			t.Fatalf("duplicate slot %d", slot)
		}
		seen[slot] = true
	}

	var decoded protocache.PerfectHashTable
	if !decoded.InitFromEncoded(table.EncodedBytes()) {
		t.Fatal("failed to decode perfect hash table")
	}
	for _, key := range keys {
		if got, want := decoded.Lookup(key), table.Lookup(key); got != want {
			t.Fatalf("decoded lookup = %d, want %d", got, want)
		}
	}
}
