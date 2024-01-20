# ProtoCache-Go

Alternative flat binary format for [Protobuf schema](https://protobuf.dev/programming-guides/proto3/). It' works like FlatBuffers, but it's usually smaller and surpports map. Flat means no deserialization overhead. [A benchmark](test/benchmark) shows the Protobuf has considerable deserialization overhead and significant reflection overhead. FlatBuffers is fast but wastes space. ProtoCache takes balance of data size and read speed, so it's useful in data caching.

|  | Protobuf | ProtoCache | FlatBuffers |
|:-------|----:|----:|----:|
| Wire format size | 574B | 780B | 1296B |
| Decode + Traverse + Dealloc | 11451ns | 1121ns | 2191ns |
| Decode + Traverse + Dealloc (reflection) | 25458ns | 2266ns | - |

See detail in [C++ implement](https://github.com/PeterRK/protocache).

## Code Gen
```sh
protoc --pcgo_out=. test.proto
```
A protobuf compiler plugin called `protoc-gen-pcgo` is available to generate Go file. The generated file is short and human friendly. Don't mind to edit it if nessasery.

## Basic APIs
```go
raw, err := protocache.Serialize(pbMessage)
assert(t, err == nil)

root := pc.AS_Main(raw)
assert(t, root.IsValid())
```
Serializing a protobuf message with protocache.Serialize is the only way to create protocache binary at present. It's easy to access by wrapping the data with generated code.

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
The reflection apis are simliar to protobuf's. An example can be found in the [test](test/reflect_test.go). CGO is needed to parse schema.