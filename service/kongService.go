package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"kong-configurer/converter"
	. "kong-configurer/logging"
	"kong-configurer/model"
	"net/http"
	"reflect"
)

var (
	client = &http.Client{}
)

type kongService struct {
	cfg               model.Config
	rootPath          string
	kongAdminPassword string
}

func NewKongService(cfg model.Config) kongService {
	return kongService{
		cfg:               cfg,
		rootPath:          cfg.Connection.KongPath,
		kongAdminPassword: cfg.Connection.KongPassword,
	}
}

func (ks kongService) ProcessMigration() {
	for _, s := range ks.cfg.Routing {
		ks.ensureService(s)
		ks.ensureRoutes(s.Routes, s.ServiceName)
	}
}

func (ks kongService) ensureService(r model.Routing) {
	exist, modified := ks.checkService(r)
	switch exist {
	case true:
		if modified {
			ks.modifyService(r)
		} else {
			LogMsgToFile(fmt.Sprintf("Service `%s` already exist and does not changed", r.ServiceName))
		}
	case false:
		ks.addService(r)
	}
}

func (ks kongService) addService(r model.Routing) {
	addServiceUrl := fmt.Sprintf("%s/services", ks.rootPath)
	data := model.AddServiceRequest{
		Name: r.ServiceName,
		Url:  r.URL,
	}
	body, _ := json.Marshal(data)
	request, _ := http.NewRequest(http.MethodPost, addServiceUrl, bytes.NewBuffer(body))
	ks.doWithAuth(request)
}

func (ks kongService) modifyService(r model.Routing) {
	modifyServiceUrl := fmt.Sprintf("%s/services/%s", ks.rootPath, r.ServiceName)
	data := model.AddServiceRequest{
		Name: r.ServiceName,
		Url:  r.URL,
	}
	body, _ := json.Marshal(data)
	request, _ := http.NewRequest(http.MethodPatch, modifyServiceUrl, bytes.NewBuffer(body))
	ks.doWithAuth(request)
}

func (ks kongService) ensureRoutes(requestedRoutes []model.AddRouteRequest, serviceName string) {
	existingRoutes :=
		converter.ToAddRouteRequest(ks.getAllRouteNamesForService(serviceName)...)
	for _, existingRoute := range existingRoutes {
		if contains, _ := containRoute(requestedRoutes, existingRoute); !contains {
			ks.removeRoute(existingRoute.Name, serviceName)
		}
	}
	for _, requestedRoute := range requestedRoutes {
		contains, modified := containRoute(existingRoutes, requestedRoute)
		switch contains {
		case false:
			ks.addRoute(requestedRoute, serviceName)
		case true:
			if modified {
				ks.modifyRoute(requestedRoute, serviceName)
			} else {
				LogMsgToFile(fmt.Sprintf("Route `%s` already exist and does not changed", requestedRoute.Name))
			}
		}
	}
}

func containRoute(routes []model.AddRouteRequest, expectingRouteName model.AddRouteRequest) (contain bool, modified bool) {
	for _, actualRoute := range routes {
		if actualRoute.Name == expectingRouteName.Name {
			if reflect.DeepEqual(actualRoute, expectingRouteName) {
				return true, false
			}
			return true, true
		}
	}
	return false, true
}

func (ks kongService) addRoute(addRouteBody model.AddRouteRequest, serviceName string) {
	routeCfgUrl := fmt.Sprintf("%s/services/%s/routes", ks.rootPath, serviceName)
	body, _ := json.Marshal(addRouteBody)
	request, _ := http.NewRequest(http.MethodPost, routeCfgUrl, bytes.NewBuffer(body))
	ks.doWithAuth(request)
}

func (ks kongService) removeRoute(routeName string, serviceName string) {
	removeRouteUrl := fmt.Sprintf("%s/services/%s/routes/%s", ks.rootPath, serviceName, routeName)
	request, _ := http.NewRequest(http.MethodDelete, removeRouteUrl, nil)
	ks.doWithAuth(request)
}

func (ks kongService) modifyRoute(addRouteBody model.AddRouteRequest, serviceName string) {
	modifyRouteUrl := fmt.Sprintf("%s/services/%s/routes/%s", ks.rootPath, serviceName, addRouteBody.Name)
	body, _ := json.Marshal(addRouteBody)
	request, _ := http.NewRequest(http.MethodPatch, modifyRouteUrl, bytes.NewBuffer(body))
	ks.doWithAuth(request)
}

func (ks kongService) getService(routing model.Routing) (service *model.KongServiceModel) {
	getServiceUrl := fmt.Sprintf("%s/services/%s", ks.rootPath, routing.ServiceName)
	request, _ := http.NewRequest(http.MethodGet, getServiceUrl, nil)
	rawResponse := ks.doWithAuth(request)
	var serviceResponse model.KongServiceModel
	err := json.NewDecoder(rawResponse.Body).Decode(&serviceResponse)
	LogOnError(err, "Error during parsing json")
	return &serviceResponse
}

func (ks kongService) getAllRouteNamesForService(serviceName string) (routes []model.KongRouteModel) {
	getAllRoutesUrl := fmt.Sprintf("%s/services/%s/routes", ks.rootPath, serviceName)
	request, _ := http.NewRequest(http.MethodGet, getAllRoutesUrl, nil)
	rawResponse := ks.doWithAuth(request)
	var routeResponse model.GetRoutesResponse
	err := json.NewDecoder(rawResponse.Body).Decode(&routeResponse)
	LogOnError(err, "Error during parsing json")
	return routeResponse.Data
}

func (ks kongService) doWithAuth(request *http.Request) (response *http.Response) {
	request.SetBasicAuth(ks.cfg.Connection.KongUser, ks.kongAdminPassword)
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		LogOnError(err, "cannot process request")
	}
	LogToFile(request, response)
	return
}

func (ks kongService) checkService(requestService model.Routing) (exist bool, modify bool) {
	existingService := ks.getService(requestService)
	if existingService == nil {
		return false, true
	}
	if isServiceModify(existingService, requestService) {
		return true, true
	}
	return true, false
}

func isServiceModify(existingService *model.KongServiceModel, requestedService model.Routing) bool {
	existingURL := fmt.Sprintf("%s://%s:%d%s", existingService.Protocol, existingService.Host, existingService.Port, existingService.Path)
	return requestedService.URL != existingURL
}
