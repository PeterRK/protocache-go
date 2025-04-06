# ProtoCache Go

Alternative flat binary format for [Protobuf schema](https://protobuf.dev/programming-guides/proto3/). It' works like FlatBuffers, but it's usually smaller and surpports map. Flat means no deserialization overhead. [A benchmark](test/benchmark) shows the Protobuf has considerable deserialization overhead and significant reflection overhead. FlatBuffers is fast but wastes space. ProtoCache takes balance of data size and read speed, so it's useful in data caching.

|  | Data Size | Decode + Traverse | Decode + Traverse(reflection) |
|:-----|----:|-----:|-----:|
| Protobuf | 574B | 5998ns | 12480ns |
| vtprotobuf | 574B | 3359ns | 10294ns |
| ProtoCache  | 780B  | 586ns | 1260ns |
| FlatBuffers | 1296B | 1129ns | No Go API |

See detail in [C++ version](https://github.com/peterrk/protocache).

## Code Gen
```sh
protoc --pcgo_out=. test.proto
```
A protobuf compiler plugin called `protoc-gen-pcgo` is [available](cmd/protoc-gen-pcgo) to generate Go file.

## Basic APIs
```go
raw, err := protocache.Serialize(pbMessage)
assert(t, err == nil)

root := pc.AS_Main(raw)
assert(t, root.IsValid())
```
Serializing a protobuf message with `protocache.Serialize` is the only way to create protocache binary at present. It's easy to access by wrapping the data with generated code.

## Reflection
```go
std::string err;
raw, err := os.ReadFile("test.proto")
assert(t, err == nil)

proto, err := compiler.ParseProto(raw) //CGO
assert(t, err == nil)

var pool reflect.DescriptorPool
assert(t, pool.Register(proto))

root := pool.Find("test.Main")
assert(t, root != nil)

field := root.Lookup("f64")
assert(t, field != nil)
```
The reflection apis are simliar to C++ version. An example can be found in the [test](test/reflect_test.go). CGO is needed to parse schema.