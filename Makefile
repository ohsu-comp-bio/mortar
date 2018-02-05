
proto:
	protoc \
		-I ./ \
    -I $$GOPATH/src/github.com/golang/protobuf/ptypes/struct/ \
    -I $$GOPATH/src/github.com/golang/protobuf/ptypes/timestamp/ \
    events/events.proto
