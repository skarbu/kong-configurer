package model

import (
	"errors"
	logger "kong-configurer/logging"
)

func (connCfg ConnectionConfig) Validate() (errs []error) {
	errs = make([]error, 0)
	if connCfg.KongPath == "" {
		logger.LogOnError(errors.New("invalid configuration"), "kong path cannot be empty")
		logger.FatalInvalidArgs()
	}
	if connCfg.KongUser == "" {
		logger.LogOnError(errors.New("invalid configuration"), "kong user cannot be empty")
		logger.FatalInvalidArgs()
	}
	if connCfg.KongPassword == "" {
		logger.LogOnError(errors.New("invalid configuration"), "kong password cannot be empty")
		logger.FatalInvalidArgs()
	}
	return errs
}

func (r Service) Validate() (errs []error) {
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
