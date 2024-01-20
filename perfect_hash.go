package protocache

import (
	"time"
)

func clearAll(data []byte) {
	vec := upCast[byte, uint64](data)
	for i := 0; i < len(vec); i++ {
		vec[i] = 0
	}
	for i := len(vec) * 8; i < len(data); i++ {
		data[i] = 0
	}
}

func setAll(data []byte) {
	vec := upCast[byte, uint64](data)
	for i := 0; i < len(vec); i++ {
		vec[i] = ^uint64(0)
	}
	for i := len(vec) * 8; i < len(data); i++ {
		data[i] = 0xff
	}
}

func bit2(bitmap []byte, pos uint32) byte {
	blk, sft := pos>>2, (pos&3)<<1
	return (bitmap[blk] >> sft) & 3
}

func setBit2on11(bitmap []byte, pos uint32, val byte) {
	blk, sft := pos>>2, (pos&3)<<1
	bitmap[blk] ^= ((^val & 3) << sft)
}

func countValidSlot(v uint64) uint32 {
	v &= (v >> 1)
	v = (v & 0x1111111111111111) + ((v >> 2) & 0x1111111111111111)
	v = v + (v >> 4)
	v = v + (v >> 8)
	v = (v & 0xf0f0f0f0f0f0f0f) + ((v >> 16) & 0xf0f0f0f0f0f0f0f)
	v = v + (v >> 32)
	return 32 - (uint32(v) & 0xff)
}

func setBit(bitmap []byte, pos uint32) {
	blk, sft := pos>>3, pos&7
	bitmap[blk] |= byte(1) << sft
}

func testAndSetBit(bitmap []byte, pos uint32) bool {
	blk, sft := pos>>3, pos&7
	m := byte(1) << sft
	if (bitmap[blk] & m) != 0 {
		return false
	}
	bitmap[blk] |= m
	return true
}

func calcSectionSize(size uint32) uint32 {
	size = uint32((uint64(size)*105 + 255) / 256)
	if size < 10 {
		return 10
	}
	return size
}

func calcBitmapSize(section uint32) uint32 {
	return ((section*3 + 31) & ^uint32(31)) / 4
}

type PerfectHash struct {
	data    []byte
	size    uint32
	section uint32
}

func (h *PerfectHash) IsValid() bool {
	return len(h.data) != 0
}

func (h *PerfectHash) Data() []byte {
	return h.data
}

func (h *PerfectHash) Size() uint32 {
	return h.size
}

func (h *PerfectHash) Init(data []byte) bool {
	if len(data) < 4 {
		return false
	}
	size := getUint32(data) & 0xfffffff
	if size < 2 {
		h.size = size
		h.section = 0
		h.data = data[:4]
		return true
	}
	section := calcSectionSize(size)
	n := calcBitmapSize(section)
	if size > 0xffff {
		n += n / 2
	} else if size > 0xff {
		n += n / 4
	} else if size > 24 {
		n += n / 8
	}
	n += 8
	if len(data) < int(n) {
		return false
	}
	h.size = size
	h.section = section
	h.data = data[:n]
	return true
}

func (h *PerfectHash) Locate(key []byte) uint32 {
	if h.size < 2 {
		return 0
	}
	code := hash96(getUint32(h.data[4:8]), key)
	code[0] %= h.section
	code[1] = code[1]%h.section + h.section
	code[2] = code[2]%h.section + h.section*2

	bitmap := h.data[8:]
	m := bit2(bitmap, code[0]) + bit2(bitmap, code[1]) + bit2(bitmap, code[2])
	slot := code[m%3]

	a, b := slot>>5, slot&31
	table := h.data[8+calcBitmapSize(h.section):]

	off := uint32(0)
	if h.size > 0xffff {
		off = getUint32(table[a*4:])
	} else if h.size > 0xff {
		off = uint32(getUint16(table[a*2:]))
	} else if h.size > 24 {
		off = uint32(table[a])
	}

	block := getUint64(bitmap[a*8:])
	block |= uint64(0xffffffffffffffff) << (b << 1)
	return off + countValidSlot(block)
}

type unsigned interface {
	uint8 | uint16 | uint32 | uint64
}

type vertex[T unsigned] struct {
	slot T
	next T
}

type graph[T unsigned] struct {
	edges [][3]vertex[T]
	nodes []T
	sizes []uint8
}

type KeySource interface {
	Reset()
	Total() int
	Next() []byte
}

func (g *graph[T]) init(seed uint32, src KeySource) bool {
	setAll(castToBytes(g.nodes))
	clearAll(castToBytes(g.sizes))

	section := uint32(len(g.nodes) / 3)
	total := src.Total()
	src.Reset()
	for i := 0; i < total; i++ {
		key := src.Next()
		code := hash96(seed, key)
		code[0] %= section
		code[1] = code[1]%section + section
		code[2] = code[2]%section + section*2
		for j := 0; j < 3; j++ {
			v := &g.edges[i][j]
			v.slot = T(code[j])
			v.next = g.nodes[v.slot]
			g.nodes[v.slot] = T(i)
			g.sizes[v.slot]++
			if g.sizes[v.slot] > 50 {
				return false
			}
		}
	}
	return true
}

func (g *graph[T]) tear(free []T, book []byte) []T {
	clearAll(book)
	free = free[:0]
	for i := len(g.edges) - 1; i >= 0; i-- {
		edge := g.edges[i]
		for j := 0; j < 3; j++ {
			v := edge[j]
			if g.sizes[v.slot] == 1 && testAndSetBit(book, uint32(i)) {
				free = append(free, T(i))
			}
		}
	}
	for head := 0; head < len(free); head++ {
		curr := free[head]
		for j := 0; j < 3; j++ {
			v := &g.edges[curr][j]
			p := &g.nodes[v.slot]
			for *p != curr {
				p = &g.edges[*p][j].next
			}
			*p = v.next
			v.next = ^T(0)
			g.sizes[v.slot]--
			i := g.nodes[v.slot]
			if g.sizes[v.slot] == 1 && testAndSetBit(book, uint32(i)) {
				free = append(free, i)
			}
		}
	}
	return free
}

func (g *graph[T]) mapping(free []T, book []byte, bitmap []byte) {
	clearAll(book)
	setAll(bitmap)
	for i := len(free) - 1; i >= 0; i-- {
		edge := g.edges[free[i]]
		a := uint32(edge[0].slot)
		b := uint32(edge[1].slot)
		c := uint32(edge[2].slot)
		switch {
		case testAndSetBit(book, a):
			setBit(book, b)
			setBit(book, c)
			sum := bit2(bitmap, b) + bit2(bitmap, c)
			setBit2on11(bitmap, a, (6-sum)%3)
		case testAndSetBit(book, b):
			setBit(book, c)
			sum := bit2(bitmap, a) + bit2(bitmap, c)
			setBit2on11(bitmap, b, (7-sum)%3)
		case testAndSetBit(book, c):
			sum := bit2(bitmap, a) + bit2(bitmap, b)
			setBit2on11(bitmap, c, (8-sum)%3)
		default:
			panic("all nodes are occupied")
		}
	}
}

func build[T unsigned](src KeySource) []byte {
	total := src.Total()
	if total <= 1 || total > 0xfffffff {
		return nil
	}
	size := uint32(total)
	section := calcSectionSize(size)
	bmsz := calcBitmapSize(section)
	bytes := 8 + bmsz
	if bmsz > 8 {
		bytes += (bmsz / 8) * uint32(sizeof[T]())
	}
	out := make([]byte, bytes)
	putUint32(out, size)
	bitmap := out[8 : 8+bmsz]

	var table []T
	if bmsz > 8 {
		table = upCast[byte, T](out[8+bmsz:])
	}

	slotCnt := section * 3

	g := graph[T]{
		edges: make([][3]vertex[T], size),
		nodes: make([]T, slotCnt),
		sizes: make([]uint8, slotCnt),
	}
	free := make([]T, 0, size)
	book := make([]byte, (slotCnt+7)/8)

	chance := 16
	if sizeof[T]() == 1 {
		chance = 40
	}
	var xs xorshift
	xs.init(uint32(time.Now().UnixNano()))
	for ; chance >= 0; chance-- {
		seed := xs.next()
		putUint32(out[4:8], seed)
		if !g.init(seed, src) {
			continue
		}
		free = g.tear(free, book)
		if len(free) != total {
			continue
		}
		g.mapping(free, book, bitmap)
		if bmsz > 8 {
			vec := upCast[byte, uint64](bitmap)
			cnt := uint32(0)
			for i := 0; i < len(table); i++ {
				table[i] = T(cnt)
				cnt += countValidSlot(vec[i])
			}
			if cnt != size {
				panic("item lost")
			}
		}
		return out
	}
	return nil
}

func Build(src KeySource) PerfectHash {
	total := src.Total()
	if total > 0xfffffff {
		return PerfectHash{}
	}
	out := PerfectHash{size: uint32(total)}
	if total > 0xffff {
		out.data = build[uint32](src)
	} else if total > 0xff {
		out.data = build[uint16](src)
	} else if total > 1 {
		out.data = build[uint8](src)
	} else {
		out.data = make([]byte, 4)
		putUint32(out.data, out.size)
		return out
	}
	if out.data == nil {
		return PerfectHash{}
	}
	out.section = calcSectionSize(out.size)
	return out
}

type xorshift struct {
	x, y, z, w uint32
}

func (xs *xorshift) init(seed uint32) {
	xs.x, xs.y, xs.z = 0x6c078965, 0x9908b0df, 0x9d2c5680
	xs.w = seed
}

func (xs *xorshift) next() uint32 {
	t := xs.x ^ (xs.x << 11)
	xs.x, xs.y, xs.z = xs.y, xs.z, xs.w
	xs.w ^= (xs.w >> 19) ^ t ^ (t >> 8)
	return xs.w
}
