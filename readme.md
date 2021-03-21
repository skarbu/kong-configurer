### What it is?
Small app develop with Go that automate loading of routing configuration to kong

### Main features

Batch operations including:

- add services
- add/remove/modify routes

### How to use it?

```go run kong-configurer.go -f {configFile} -h {kongHost}  -u {user} -p {password}```
```
  -f string - relative path to json config file
  -h string - kong host with port
  -p string - password for kong user
  -u string - kong user with admin privileges
```
param -p is optional, if not passed app ask for password with blind password prompt

You can also build and run app as binaries
```
go build kong-configurer.go
kong-configurer  -f {configFile} -h {kongHost}  -u {user} -p {password}
```

### Config file structure
```
{
  "routing": [
    {
      "serviceName": "{service name}",
      "url": "{service root url}",
      "routes": [
        {
          "name": "{route-name}",
          "paths": [ "{route path}" ],
          "methods": [ "{methods}" ] // list of methods for routed path: POST, GET etc.
          "preserve_host" : false,
          "strip_path" : false
        }
      ]
    }
  ]
}
```

### TODO
- support for plugins configuration
- work with bigger set of routes (kong rest api provide pagination)

