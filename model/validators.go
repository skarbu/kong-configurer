package model

import (
	"errors"
)

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

func (r AddRouteRequest) Validate() (errs []error) {
	if r.Name == "" {
		errs = append(errs, errors.New("route.name cannot be empty"))
	}
	if len(r.Paths) == 0 {
		errs = append(errs, errors.New("route.paths cannot be empty"))
	}
	if len(r.Methods) == 0 {
		errs = append(errs, errors.New("route.method cannot be empty"))
	}
	return errs
}
