[![Build Status](https://travis-ci.org/ohsu-comp-bio/mortar.svg?branch=master)](https://travis-ci.org/ohsu-comp-bio/mortar)
[![Gitter](https://badges.gitter.im/ohsu-comp-bio/mortar.svg)](https://gitter.im/ohsu-comp-bio/mortar)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Godoc](https://img.shields.io/badge/godoc-ref-blue.svg)](http://godoc.org/github.com/ohsu-comp-bio/mortar)

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
