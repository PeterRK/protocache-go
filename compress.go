package protocache

import (
	"errors"
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
		switch ch {
		case 0:
			for k < len(src) && cnt < 4 && src[k] == 0 {
				k++
				cnt++
			}
			return 0x8 | (cnt - 1)
		case 0xff:
			for k < len(src) && cnt < 4 && src[k] == 0xff {
				k++
				cnt++
			}
			return 0xC | (cnt - 1)
		default:
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
	out := make([]byte, 0, len(src))

	k := 0
	size := uint32(0)
	for sft := 0; sft < 32; sft += 7 {
		if k >= len(src) {
			return nil, errors.New("broken header")
		}
		b := uint32(src[k])
		k++
		if (b & 0x80) != 0 {
			size |= (b & 0x7f) << sft
		} else {
			size |= b << sft
			break
		}
	}

	unpack := func(mark uint8) bool {
		if (mark & 8) != 0 {
			cnt := (mark & 3) + 1
			ch := uint8(0)
			if (mark & 4) != 0 {
				ch = 0xff
			}
			for ; cnt != 0; cnt-- {
				out = append(out, ch)
			}
		} else {
			l := int(mark & 7)
			if k+l > len(src) {
				return false
			}
			out = append(out, src[k:k+l]...)
			k += l
		}
		return true
	}

	for k < len(src) {
		mark := src[k]
		k++
		if !unpack(mark&0xf) || !unpack(mark>>4) {
			return nil, errors.New("broken data")
		}
	}
	if len(out) != int(size) {
		return nil, errors.New("size mismatch")
	}
	return out, nil
}
