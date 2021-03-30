package model

import "errors"

type AddRouteRequestInterface interface {
	ToAddRouteRequest() AddRouteRequest
}

func (r KongRouteResponseModel) ToAddRouteRequest() AddRouteRequest {
	return AddRouteRequest{
		Name:         r.Name,
		Paths:        r.Paths,
		PreserveHost: r.PreserveHost,
		StripPath:    r.StripPath,
		Methods:      r.Methods,
	}
}

func (r Route) ToAddRouteRequest() AddRouteRequest {
	return AddRouteRequest{
		Name:         r.Name,
		Paths:        r.Paths,
		PreserveHost: r.PreserveHost,
		StripPath:    r.StripPath,
		Methods:      r.Methods,
	}
}

func (p1 Plugin) ToAddPluginRequest(scope string, scopeId string) (*AddPluginRequest, error) {
	var consumer *IdModel
	if p1.Consumer != nil {
		consumer = &IdModel{ID: *p1.Consumer}
	}
	switch scope {
	case ROUTE:
		return &AddPluginRequest{
			Enabled:   true,
			Protocols: p1.Protocols,
			Name:      p1.Name,
			Consumer:  consumer,
			Service:   nil,
			Route:     &IdModel{ID: scopeId},
			Config:    p1.Config,
		}, nil
	case SERVICE:
		return &AddPluginRequest{
			Enabled:   true,
			Protocols: p1.Protocols,
			Name:      p1.Name,
			Consumer:  consumer,
			Service:   &IdModel{ID: scopeId},
			Route:     nil,
			Config:    p1.Config,
		}, nil
	default:
		return nil, errors.New("invalid scope passed to func")
	}

}
