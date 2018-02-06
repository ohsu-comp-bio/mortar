
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


