package model

import (
	"errors"
	"fmt"
)

const (
	OperationAdd    = "ADD"
	OperationRemove = "REMOVE"
)

type Config struct {
	Global struct {
		KongHost string `json:"kongHost"`
		KongPort int    `json:"kongPort"`
		KongUser string `json:"kongUser"`
	} `json:"config"`
	Routing []Routing `json:"routing"`
}

type Routing struct {
	ServiceName     string  `json:"serviceName"`
	URL             string  `json:"url"`
	RemoveAllRoutes bool    `json:"removeAllRoutesBeforeMigration"`
	Routes          []Route `json:"routes"`
}

type Route struct {
	Operation string   `json:"operation"`
	RouteName string   `json:"routeName"`
	Path      string   `json:"path"`
	Methods   []string `json:"methods"`
}

func (c Config) Validate() (errs []error) {
	errs = make([]error, 0)
	if c.Global.KongHost == "" {
		errs = append(errs, errors.New("kong host cannot be empty"))
	}
	if c.Global.KongPort == 0 {
		errs = append(errs, errors.New("kong host cannot be empty"))
	}
	if c.Global.KongUser == "" {
		errs = append(errs, errors.New("kong user cannot be empty"))
	}
	if len(c.Routing) == 0 {
		errs = append(errs, errors.New("no routing defined in file"))
	}
	for _, routing := range c.Routing {
		errs = append(errs, routing.Validate()...)
	}
	return errs
}

func (r Routing) Validate() (errs []error) {
	if r.ServiceName == "" {
		errs = append(errs, errors.New("service name cannot be empty"))
	}
	if r.URL == "" {
		errs = append(errs, errors.New("url cannot be empty"))
	}
	if len(r.Routes) == 0 {
		errs = append(errs, errors.New("routes cannot be empty"))
	}
	for _, route := range r.Routes {
		errs = append(errs, route.Validate()...)
	}
	return errs
}

func (r Route) Validate() (errs []error) {
	if r.Operation == "" {
		errs = append(errs, errors.New("route.operation cannot be empty"))
	}
	if r.Operation != OperationAdd && r.Operation != OperationRemove {
		msg := fmt.Sprintf("Wrong route.operation: %s allowed: %s, %s", r.Operation, OperationAdd, OperationRemove)
		errs = append(errs, errors.New(msg))
	}
	if r.RouteName == "" {
		errs = append(errs, errors.New("route.name cannot be empty"))
	}
	if r.Operation == OperationAdd {
		if r.Path == "" {
			errs = append(errs, errors.New("route.path cannot be empty"))
		}
		if len(r.Methods) == 0 {
			errs = append(errs, errors.New("route.method cannot be empty"))
		}
	}
	return errs
}
