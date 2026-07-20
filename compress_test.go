package protocache

import (
	"bytes"
	"testing"
)

func TestCompressRoundTrip(t *testing.T) {
	t.Parallel()

	cases := [][]byte{
		nil,
		{},
		{0},
		{0xff},
		{1, 2, 3, 4, 5, 6, 7},
		bytes.Repeat([]byte{0}, 32),
		bytes.Repeat([]byte{0xff}, 32),
		[]byte("mixed\x00payload\xffwith runs\x00\x00\x00"),
	}
	for _, input := range cases {
		encoded := Compress(input)
		decoded, err := Decompress(encoded)
		if err != nil {
			t.Fatalf("Decompress(Compress(%x)): %v", input, err)
		}
		if !bytes.Equal(decoded, input) {
			t.Fatalf("round trip mismatch: got %x, want %x", decoded, input)
		}
	}
}

func TestDecompressRejectsMalformedInput(t *testing.T) {
	t.Parallel()

	cases := [][]byte{
		{0x80},
		{0x80, 0x80, 0x80, 0x80, 0x10},
		{1},
		{4, 0x01},
		{8, 0x11, 1},
		{4, 0x00},
		{1, 0x0b},
		{1, 0x0f},
		{0xff, 0xff, 0xff, 0xff, 0x0f, 0},
	}
	for _, input := range cases {
		if output, err := Decompress(input); err == nil {
			t.Errorf("Decompress(%x) unexpectedly succeeded with %x", input, output)
		}
	}
}

func FuzzDecompress(f *testing.F) {
	seeds := [][]byte{
		nil,
		{0},
		{0x80},
		Compress([]byte("hello")),
		Compress(bytes.Repeat([]byte{0}, 32)),
		Compress(bytes.Repeat([]byte{0xff}, 32)),
	}
	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, encoded []byte) {
		decoded, err := Decompress(encoded)
		if err != nil {
			return
		}
		roundTrip, err := Decompress(Compress(decoded))
		if err != nil || !bytes.Equal(roundTrip, decoded) {
			t.Fatalf("successful decode did not round trip: err=%v", err)
		}
	})
}
