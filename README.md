# Mortar

### Usage

Start Mongo:
```
make start-mongodb
```

Start Kafka:
```
make start-kafka
```

Start Arachne:
```
make start-arachne
```

Start Funnel:
```
make start-funnel
```

Start Mortar:
```
go run cmd/mortar/import.go
```

Run a Task:
```
funnel run 'echo hi' --tag hello=mortar -i I=./dev/funnel-kafka.config.yml
```
