package protocache

import (
	"errors"
	"unsafe"
)

// Compress returns the ProtoCache compressed representation of src; it returns nil for empty input.
func Compress(src []byte) []byte {
	if len(src) == 0 {
		return nil
	}
	out := make([]byte, 0, len(src))

	n := uint32(len(src))
	for (n & ^uint32(0x7f)) != 0 {
		out = append(out, byte(0x80|(n&0x7f)))
		n >>= 7
	}
	out = append(out, byte(n))

	k := 0
	pick := func() uint8 {
		cnt := uint8(1)
		ch := src[k]
		k++
		if int8(ch) == (int8(ch) >> 1) {
			for k < len(src) && cnt < 4 && src[k] == ch {
				k++
				cnt++
			}
			return 0x8 | (ch & 0x4) | (cnt - 1)
		} else {
			for k < len(src) && cnt < 7 && src[k] != 0 && src[k] != 0xff {
				k++
				cnt++
			}
			return cnt
		}
	}

	for k < len(src) {
		x := k
		a := pick()
		if k == len(src) {
			out = append(out, a)
			if (a & 0x8) == 0 {
				out = append(out, src[x:x+int(a)]...)
			}
			break
		}
		y := k
		b := pick()
		out = append(out, a|(b<<4))
		if (a & 0x8) == 0 {
			out = append(out, src[x:x+int(a)]...)
		}
		if (b & 0x8) == 0 {
			out = append(out, src[y:y+int(b)]...)
		}
	}

	return out
}

// Decompress decodes ProtoCache compressed data and rejects malformed input; it returns nil, nil for empty input.
func Decompress(src []byte) ([]byte, error) {
	if len(src) == 0 {
		return nil, nil
	}

	var size uint32
	k := 0
	for shift := uint(0); ; shift += 7 {
		if k >= len(src) || shift >= 32 {
			return nil, errors.New("broken header")
		}
		b := src[k]
		k++
		if shift == 28 && b&0xf0 != 0 {
			return nil, errors.New("broken header")
		}
		size |= uint32(b&0x7f) << shift
		if b&0x80 == 0 {
			break
		}
	}

	encoded := src[k:]
	if uint64(size) > uint64(len(encoded))*8 {
		return nil, errors.New("broken data")
	}
	maxInt := int(^uint(0) >> 1)
	if uint64(size) > uint64(maxInt-7) {
		return nil, errors.New("data too large")
	}

	tail := int(size)
	out := make([]byte, tail+7)
	source, dest := 0, 0
	unpack := func(mark uint8) bool {
		if mark&8 != 0 {
			count := int(mark&3) + 1
			if count > tail-dest {
				return false
			}
			value := uint32(0)
			if mark&4 != 0 {
				value = ^value
			}
			*(*uint32)(unsafe.Add(unsafe.Pointer(unsafe.SliceData(out)), dest)) = value
			dest += count
		} else if mark != 0 {
			count := int(mark)
			if count > len(encoded)-source || count > tail-dest {
				return false
			}
			if len(encoded)-source >= 8 {
				src := unsafe.Add(unsafe.Pointer(unsafe.SliceData(encoded)), source)
				dst := unsafe.Add(unsafe.Pointer(unsafe.SliceData(out)), dest)
				*(*uint64)(dst) = *(*uint64)(src)
			} else {
				copy(out[dest:dest+count], encoded[source:source+count])
			}
			source += count
			dest += count
		}
		return true
	}

	for source < len(encoded) {
		mark := encoded[source]
		source++
		if !unpack(mark&0xf) || !unpack(mark>>4) {
			return nil, errors.New("broken data")
		}
	}
	if dest != tail {
		return nil, errors.New("size mismatch")
	}
	return out[:tail], nil
}
