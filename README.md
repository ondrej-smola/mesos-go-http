# Go HTTP client for Apache Mesos

Goal of this project is to provide low and high level API for Apache Mesos using HTTP API

## Status

Current state of project is *alpha* as it is only lightly tested and need more adoption/feedback.

### Mesos versions

Current support is for Mesos 1.1.0 proto.

### Features

- Leading master detection
- Scheduler API
- Operator API
- High level Flow API
- Low level Client API
- Examples

## Installing

Users of this library are encouraged to vendor it. API stability isn't guaranteed
at this stage.

## Logging

Project uses [go-kit/log] (https://github.com/go-kit/kit/tree/v0.3.0/log) compatible interfaces for logging

## Dependencies
   
### Runtime   

 ```
 make install-dependencies
 ``` 
### Testing

```
make install-test-dependencies
```
 
### Notice

This project uses code from [mesos-go](https://github.com/mesos/mesos-go) licensed under the Apache Licence 2.0
This project uses code from [go-kit](https://github.com/go-kit/kit) licensed under the MIT Licence

## License

This project is [Apache License 2.0](LICENSE).


