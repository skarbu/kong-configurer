package migrationService

import (
	"bytes"
	"encoding/json"
	"fmt"
	. "kong-configurer/logging"
	"kong-configurer/model"
	"net/http"
	"reflect"
)

func (ms *migrationService) ensureRoutes(requestedRoutes []model.Route, serviceName string) {
	existingRoutes := ms.getAllRouteNamesForService(serviceName)
	for _, existingRoute := range existingRoutes {
		if contains, _ := routesContainsRoute(requestedRoutes, existingRoute); !contains {
			ms.removeRoute(existingRoute.Name, serviceName)
		}
	}
	for _, requestedRoute := range requestedRoutes {
		contains, modified := containRoute(existingRoutes, requestedRoute)
		switch contains {
		case false:
			ms.addRoute(&requestedRoute, serviceName)
		case true:
			if modified {
				ms.modifyRoute(&requestedRoute, serviceName)
			} else {
				LogMsgToFile(fmt.Sprintf("Route `%s` already exist and does not changed", requestedRoute.Name))
			}
		}
	}
}

//TODO investigate why GO cannot accept []intrefce. It cause that method is duplicated for type []Route and []AddRouteRequest
func routesContainsRoute(routes []model.Route, expectingRoute model.AddRouteRequestInterface) (contain bool, modified bool) {
	expectingRouteRequest := expectingRoute.ToAddRouteRequest()
	for _, actualRoute := range routes {
		actualRouteRequest := actualRoute.ToAddRouteRequest()
		if actualRouteRequest.Name == expectingRouteRequest.Name {
			if reflect.DeepEqual(actualRouteRequest, expectingRouteRequest) {
				return true, false
			}
			return true, true
		}
	}
	return false, true
}

func containRoute(routes []model.KongRouteResponseModel, expectingRoute model.AddRouteRequestInterface) (contain bool, modified bool) {
	expectingRouteRequest := expectingRoute.ToAddRouteRequest()
	for _, actualRoute := range routes {
		actualRouteRequest := actualRoute.ToAddRouteRequest()
		if actualRouteRequest.Name == expectingRouteRequest.Name {
			if reflect.DeepEqual(actualRouteRequest, expectingRouteRequest) {
				return true, false
			}
			return true, true
		}
	}
	return false, true
}

func (ms *migrationService) addRoute(route *model.Route, serviceName string) {
	routeCfgUrl := fmt.Sprintf("%s/services/%s/routes", ms.connCfg.KongPath, serviceName)
	body, _ := json.Marshal(route.ToAddRouteRequest())
	request, _ := http.NewRequest(http.MethodPost, routeCfgUrl, bytes.NewBuffer(body))
	ms.doWithAuth(request)
}

func (ms *migrationService) removeRoute(routeName string, serviceName string) {
	removeRouteUrl := fmt.Sprintf("%s/services/%s/routes/%s", ms.connCfg.KongPath, serviceName, routeName)
	request, _ := http.NewRequest(http.MethodDelete, removeRouteUrl, nil)
	ms.doWithAuth(request)
}

func (ms *migrationService) modifyRoute(route *model.Route, serviceName string) {
	modifyRouteUrl := fmt.Sprintf("%s/services/%s/routes/%s", ms.connCfg.KongPath, serviceName, route.Name)
	body, _ := json.Marshal(route.AddRouteRequest)
	request, _ := http.NewRequest(http.MethodPatch, modifyRouteUrl, bytes.NewBuffer(body))
	ms.doWithAuth(request)
}

func (ms *migrationService) getAllRouteNamesForService(serviceName string) (routes []model.KongRouteResponseModel) {
	getAllRoutesUrl := fmt.Sprintf("%s/services/%s/routes", ms.connCfg.KongPath, serviceName)
	request, _ := http.NewRequest(http.MethodGet, getAllRoutesUrl, nil)
	rawResponse := ms.doWithAuth(request)
	var routeResponse model.GetRoutesResponse
	err := json.NewDecoder(rawResponse.Body).Decode(&routeResponse)
	LogOnError(err, "Error during parsing json")
	return routeResponse.Data
}
