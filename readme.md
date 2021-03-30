### What it is?
Small app develop with Go that automate loading of routing configuration to kong

### Main features

Batch operations including:

- add services
- add/remove/modify routes
- add plugins on service and route level
- all operation are appling only when app find difference between actual configuration and applied configuration.

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
      "plugins": [
              {
                "name": "{pluginName}",
                "consumer": {consumerName},
                "protocols": [ "{protocol}"], // e.g. [ "http" ]
                "config": {
                //plugin config
                }
              }
            ],
      "routes": [
        {
          "name": "{route-name}",
          "paths": [ "{route path}" ],
          "methods": [ "{methods}" ] // list of methods for routed path: POST, GET etc.
          "preserve_host" : false,
          "strip_path" : false
           "plugins": [
                        {
                          "name": "{pluginName}",
                          "consumer": {consumerName},
                          "protocols": [ "{protocol}"], // e.g. [ "http" ]
                          "config": {
                          //plugin config
                          }
                        }
                      ],
        }
      ]
    }
  ]
}
```

### LOGS
After execution, logs are available as a file in working dir. 

### TODO
- work with bigger set of routes (kong rest api provide pagination)

