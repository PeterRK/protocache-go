package protocache

import (
	"strconv"
	"testing"
)

type reader struct {
	keys [][]byte
	curr int
}

func (r *reader) Reset() {
	r.curr = 0
}

func (r *reader) Total() int {
	return len(r.keys)
}

func (r *reader) Next() []byte {
	key := r.keys[r.curr]
	r.curr++
	return key
}

func assert(t *testing.T, state bool) {
	if !state {
		t.FailNow()
	}
}

func testPerfectHash(t *testing.T, size int) {
	keys := make([][]byte, size)
	for i := 0; i < size; i++ {
		keys[i] = castStrToBytes(strconv.Itoa(i))
	}
	table := Build(&reader{keys: keys})
	assert(t, table.IsValid())

	mark := make([]bool, size)
	for i := 0; i < size; i++ {
		pos := table.Locate(keys[i])
		assert(t, pos < uint32(size))
		assert(t, !mark[pos])
		mark[pos] = true
	}
}

func TestPerfectHashTiny(t *testing.T) {
	testPerfectHash(t, 0)
	testPerfectHash(t, 1)
	testPerfectHash(t, 2)
	testPerfectHash(t, 24)
}

func TestPerfectHashSmall(t *testing.T) {
	testPerfectHash(t, 200)
	testPerfectHash(t, 1000)
}

func TestPerfectHashBig(t *testing.T) {
	testPerfectHash(t, 100000)
}
