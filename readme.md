### What it is?
Small app develop with Go that automate loading of configuration to kong

### Main features

batch operations including:

- add services
- add routes
- removes single routes
- remove multiple route for service

### How to use it?

```go run kong-configurer.go {configFile}```
eg.
```go run kong-configurer.go routing.json```

or with build binaries
```
go build kong-configurer.go
kong-configurer {configFile}
```

### Config file structure
```
{
  "config": {
    "kongHost": "{host}",
    "kongPort": {port},
    "kongUser": "{user with admin role}"
  },
  "routing": [
    {
      "serviceName": "{service name}",
      "url": "{service root url}",
      "routes": [
        {
          "name": "{route-name}",
          "paths": [ "{route path}" ],
          "methods": [ "{methods}" ] // list of methods for routed path
          "preserve_host" : false,
          "strip_path" : false
        }
      ]
    }
  ]
}
```

### TODO
- work with bigger set of routes (kong rest api provide pagination)

