# ProtoCache Go

Alternative flat binary format for [Protobuf schema](https://protobuf.dev/programming-guides/proto3/). It' works like FlatBuffers, but it's usually smaller and surpports map. Flat means no deserialization overhead. [A benchmark](test/bench_test.go) shows the Protobuf has considerable deserialization overhead and significant reflection overhead. FlatBuffers is fast but wastes space. ProtoCache takes balance of data size and read speed, so it's useful in data caching.

|  | Protobuf | vtprotobuf | ProtoCache | FlatBuffers | Fory |
|:-------|----:|----:|----:|----:|----:|
| Data Size | 574B | 574B | 780B | 1296B | 615B |
| Compressed Size | 566B | 566B | 571B | 856B | 611B |
| Decompress | 248ns | 248ns | 487ns | 898ns | 295ns |
| Decode + Traverse | 5024ns | 2472ns | 654ns | 1208ns | 4894ns |
| Decode + Traverse(reflection) | 11301ns | 8658ns | 1327ns | No Go API | No Go API |

See detail in [C++ version](https://github.com/peterrk/protocache).

## Code Gen
```sh
protoc --pcgo_out=. [--pcgo_opt=extra,relative] test.proto
```
A protobuf compiler plugin called `protoc-gen-pcgo` is [available](cmd/protoc-gen-pcgo) to generate Go file.
Use the `relative` option to emit generated files relative to the input proto path instead of recreating the `go_package` directory hierarchy.

## Basic APIs
```go
raw, err := protocache.Serialize(pbMessage)
assert(t, err == nil)

root := pc.AS_Main(raw)
assert(t, root.IsValid())
```
Serializing a protobuf message with `protocache.Serialize` is the only way to create protocache binary at present. It's easy to access by wrapping the data with generated code.

## Mutable APIs

The `extra` option generates `EX` types for partial updates and re-serialization without protobuf reflection.

```go
raw, err := protocache.Serialize(pbMessage)
assert(t, err == nil)

root := pc.TO_MainEX(raw)

root.SetStr("patched")
obj := root.GetObject()
obj.SetI32(7)

out, err := root.Serialize()
assert(t, err == nil)
```

## Reflection
```go
std::string err;
raw, err := os.ReadFile("test.proto")
assert(t, err == nil)

proto, err := compiler.ParseProto(raw) //CGO
assert(t, err == nil)

var pool reflect.DescriptorPool
assert(t, pool.Register(proto) == nil)

root := pool.Find("test.Main")
assert(t, root != nil)

field := root.Lookup("f64")
assert(t, field != nil)
```
The reflection apis are simliar to C++ version. An example can be found in the [test](test/reflect_test.go). CGO is needed to parse schema.
