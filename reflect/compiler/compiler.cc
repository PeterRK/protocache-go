#include <cstdbool>
#include <cstdint>
#include <cstring>
#include <string>
#include <sstream>
#include <google/protobuf/io/tokenizer.h>
#include <google/protobuf/compiler/parser.h>
#include <google/protobuf/io/zero_copy_stream_impl.h>


class ErrorCollector : public google::protobuf::io::ErrorCollector {
public:
	void AddError(int line, int column, const std::string& message) override {
		core_ << "[error] line " << line << ", column " << column << ": " << message << std::endl;
	}

	void AddWarning(int line, int column, const std::string& message) override {
		core_ << "[warning] line " << line << ", column " << column << ": " << message << std::endl;
	}

	std::string String() noexcept {
		return core_.str();
	}

private:
	std::ostringstream core_;
};


extern "C" {
    
struct Buffer {
    const uint8_t* data;
    size_t size;
};

void FreeBuffer(Buffer* buf) {
    if (buf->data == nullptr) {
        return;
    }
    delete[] buf->data;
    buf->data = nullptr;
    buf->size = 0;
}

bool ParseProto(const Buffer* input, Buffer* output) {
    google::protobuf::io::ArrayInputStream stream(input->data, input->size);
    
    google::protobuf::FileDescriptorProto proto;
    
    ErrorCollector collector;
	google::protobuf::io::Tokenizer tokenizer(&stream, &collector);
	google::protobuf::compiler::Parser parser;
	bool done = parser.Parse(&tokenizer, &proto);
    std::string out;
	if (done) {
        out = proto.SerializeAsString();
    } else {
        out = collector.String();
	}
    uint8_t* buf = new uint8_t[out.size()];
    memcpy(buf, out.data(), out.size());
    output->data = buf;
    output->size = out.size();
    return done;
}

}