package test

import (
	"os"
	"testing"

	"github.com/peterrk/protocache-go"
	"github.com/peterrk/protocache-go/test/pb"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func loadBenchmarkMain(tb testing.TB) *pb.Main {
	tb.Helper()

	raw, err := os.ReadFile("test.json")
	if err != nil {
		tb.Fatal(err)
	}
	message := &pb.Main{}
	if err := protojson.Unmarshal(raw, message); err != nil {
		tb.Fatal(err)
	}
	return message
}

func loadBenchmarkProtobuf(tb testing.TB) []byte {
	tb.Helper()

	raw, err := (proto.MarshalOptions{Deterministic: true}).Marshal(loadBenchmarkMain(tb))
	if err != nil {
		tb.Fatal(err)
	}
	return raw
}

func loadBenchmarkProtoCache(tb testing.TB) []byte {
	tb.Helper()

	raw, err := protocache.Serialize(loadBenchmarkMain(tb))
	if err != nil {
		tb.Fatal(err)
	}
	return raw
}
