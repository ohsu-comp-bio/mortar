VERSION = 0.1.0

GO_CODE=$(shell find . -name '*go' -not -name '*pb.go')

V=github.com/ohsu-comp-bio/mortar/version
VERSION_LDFLAGS=\
 -X "$(V).BuildDate=$(shell date)" \
 -X "$(V).GitCommit=$(shell git rev-parse --short HEAD)" \
 -X "$(V).GitBranch=$(shell git symbolic-ref -q --short HEAD)" \
 -X "$(V).GitUpstream=$(shell git remote get-url $(shell git config branch.$(shell git symbolic-ref -q --short HEAD).remote) 2> /dev/null)" \
 -X "$(V).Version=$(VERSION)"


# Build the code
install:
	@touch version/version.go
	@go install -ldflags '$(VERSION_LDFLAGS)' github.com/ohsu-comp-bio/mortar

# Compile protobuf to Go
proto:
	protoc \
		-I ./ \
    -I $$GOPATH/src/github.com/golang/protobuf/ptypes/struct/ \
    -I $$GOPATH/src/github.com/golang/protobuf/ptypes/timestamp/ \
    events/events.proto

start-mongodb:
	@docker rm -f mortar-mongodb-test > /dev/null 2>&1 || echo
	@docker run -d --name mortar-mongodb-test -p 27017:27017 docker.io/mongo:3.5.13 > /dev/null

start-kafka:
	@docker rm -f mortar-kafka > /dev/null 2>&1 || echo
	@docker run -d --name mortar-kafka -p 2181:2181 -p 9092:9092 --env ADVERTISED_HOST="localhost" --env ADVERTISED_PORT=9092 spotify/kafka

start-arachne:
	arachne server --db arachne --port 8082 --rpc 5757 --mongo localhost:27017

start-funnel:
	funnel server run --config dev/funnel-kafka.config.yml

# Build binaries for all OS/Architectures
cross-compile:
	@echo '=== Cross compiling... ==='
	@mkdir -p build/bin
	@for GOOS in darwin linux; do \
		for GOARCH in amd64; do \
			GOOS=$$GOOS GOARCH=$$GOARCH go build -a \
				-ldflags '$(VERSION_LDFLAGS)' \
				-o build/bin/mortar-$$GOOS-$$GOARCH .; \
		done; \
	done

# Automatially update code formatting
tidy:
	@goimports -w $(GO_CODE)
	@gofmt -w $(GO_CODE)

# Run code style and other checks
lint:
	@gometalinter -e ".*pb.go" ./...

# Install dev. utils
deps:
	@go get github.com/alecthomas/gometalinter
