### What it is?
Small app develop with Go that automate loading of configuration to kong

### How to use it?

```go run kong-configurer.go {configFile}```
eg.
```go run kong-configurer.go routing.json```
or with build binaries
```
go build kong-configurer.go
kong-configurer {configFile}
```

