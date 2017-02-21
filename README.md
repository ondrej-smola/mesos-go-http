# Go HTTP client for Apache Mesos

Goal of this project is to provide low and high level API for Apache Mesos using HTTP API

[![Build Status](https://travis-ci.org/ondrej-smola/mesos-go-http.svg?branch=master)](https://travis-ci.org/ondrej-smola/mesos-go-http)

## Status

Current state of project is *alpha* as it needs more adoption/feedback.
Users of this library are encouraged to vendor it. API stability isn't guaranteed at this stage.

### Mesos proto versions

* (MASTER) 1.1.0+ (v1)

### Features

- Leading master detection
- Scheduler API
- Operator API
- High level Flow API
- Low level Client API
- Examples

## Get started

```
make
```
* add $(GOPATH)/bin to your PATH

##### List master tasks
```
operator master tasks -e http://127.0.0.1:5050/api/v1
```
##### Subscribe for master events
```
operator master event-stream -e http://127.0.0.1:5050/api/v1
```
##### Run example scheduler
```
scheduler -e http://127.0.0.1:5050/api/v1/scheduler --cmd=sleep --arg=15 --tasks=5
```

## Local infrastructure

For easy testing of master failover, agent failed ...

#### [Docker compose](https://docs.docker.com/compose/)

* Set DOCKER_IP to IP of docker host (e.g. 10.0.75.2 docker for windows, 127.0.0.1 local docker)
* In project root
```
docker-compose up -d
```
* Creates 3 masters and 2 agents 
* Each agent reports all host resources (means double your CPU count will be reported) 
     * Only for testing purposes (multiple offers, agent disconnects, ...)
     * For launching multiple tasks that actually utilize all resources start only 1 agent or modify per agent resources    

## Logging

Project uses [go-kit/log] (https://github.com/go-kit/kit/tree/v0.3.0/log) compatible interfaces for logging

   
## Testing
```
make install-test-dependencies
make test
```
 
## Notice

This project uses code from [mesos-go](https://github.com/mesos/mesos-go) licensed under the Apache Licence 2.0

## License

This project is [Apache License 2.0](LICENSE).