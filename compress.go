package protocache

import (
	"errors"
	"unsafe"
)

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
		x := int8(ch)
		if x == (x >> 1) {
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

func Decompress(src []byte) ([]byte, error) {
	if len(src) == 0 {
		return nil, nil
	}

	k := 0
	size := uintptr(0)
	decodeLen := func() bool {
		for sft := 0; sft < 32; sft += 7 {
			if k >= len(src) {
				return false
			}
			b := uintptr(src[k])
			k++
			if (b & 0x80) != 0 {
				size |= (b & 0x7f) << sft
			} else {
				size |= b << sft
				return true
			}
		}
		return false
	}

	if !decodeLen() {
		return nil, errors.New("broken header")
	}
	out := make([]byte, size, size+7) // extra space for write

	s := uintptr(unsafe.Pointer(&src[k]))
	end := s + uintptr(len(src)-k)
	dest := uintptr(unsafe.Pointer(&out[0]))
	tail := dest + size

	unpack := func(mark uint8) bool {
		if (mark & 8) != 0 {
			cnt := uintptr(mark&3) + 1
			if dest+cnt > tail {
				return false
			}
			*(*uint32)(unsafe.Pointer(dest)) = uint32(0) - uint32((mark>>2)&1)
			dest += cnt
		} else if mark != 0 {
			l := uintptr(mark)
			if s+l > end || dest+l > tail {
				return false
			}
			if s+8 <= end {
				*(*uint64)(unsafe.Pointer(dest)) = *(*uint64)(unsafe.Pointer(s))
			} else {
				for i := uintptr(0); i < l; i++ {
					*(*byte)(unsafe.Pointer(dest + i)) = *(*byte)(unsafe.Pointer(s + i))
				}
			}
			s += l
			dest += l
		}
		return true
	}

	for s < end {
		mark := *(*byte)(unsafe.Pointer(s))
		s++
		if !unpack(mark&0xf) || !unpack(mark>>4) {
			return nil, errors.New("broken data")
		}
	}
	if dest != tail {
		return nil, errors.New("size mismatch")
	}
	return out, nil
}
