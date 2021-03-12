package converter

import "kong-configurer/model"

func ToAddRouteRequest(r ...model.KongRouteModel) []model.AddRouteRequest {
	routes := make([]model.AddRouteRequest, len(r))
	for i, r := range r {
		routes[i] = model.AddRouteRequest{
			Name:         r.Name,
			Paths:        r.Paths,
			PreserveHost: r.PreserveHost,
			StripPath:    r.StripPath,
			Methods:      r.Methods,
		}
	}
	return routes
}
