
proto:
	protoc \
		-I ./ \
    -I $$GOPATH/src/github.com/golang/protobuf/ptypes/struct/ \
    -I $$GOPATH/src/github.com/golang/protobuf/ptypes/timestamp/ \
    --go_out=Mtes/tes.proto=github.com/ohsu-comp-bio/tes:. \
    events/events.proto
