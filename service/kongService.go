package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	. "kong-configurer/logging"
	"kong-configurer/model"
	"net/http"
)

var (
	client = &http.Client{}
)

type kongService struct {
	cfg               model.Config
	rootPath          string
	kongAdminPassword string
}

func NewKongService(cfg model.Config, kongAdminPassword string) kongService {
	rootPath := fmt.Sprintf("%s:%d", cfg.Global.KongHost, cfg.Global.KongPort)
	return kongService{
		cfg:               cfg,
		rootPath:          rootPath,
		kongAdminPassword: kongAdminPassword,
	}
}

func (ks kongService) ProcessMigration() {
	for _, r := range ks.cfg.Routing {
		ks.ensureService(r)
		if r.RemoveAllRoutes {
			ks.removeAllRoutesFromService(r.ServiceName)
		}
		ks.ensureRoutes(r.Routes, r.ServiceName)
	}
}

func (ks kongService) ensureService(r model.Routing) {
	serviceCfgUrl := fmt.Sprintf("%s/services", ks.rootPath)
	data := model.ServiceRequest{
		Name: r.ServiceName,
		Url:  r.URL,
	}
	body, _ := json.Marshal(data)
	request, _ := http.NewRequest(http.MethodPost, serviceCfgUrl, bytes.NewBuffer(body))
	ks.doWithAuth(request)
}

func (ks kongService) ensureRoutes(routes []model.Route, serviceName string) {
	for _, r := range routes {
		switch op := r.Operation; op {
		case model.OperationAdd:
			ks.addRoute(r, serviceName)
		case model.OperationRemove:
			ks.removeRoute(r.RouteName, serviceName)
		}
	}
}

func (ks kongService) addRoute(r model.Route, serviceName string) {
	routeCfgUrl := fmt.Sprintf("%s/services/%s/routes", ks.rootPath, serviceName)
	paths := make([]string, 1)
	paths[0] = r.Path
	data := model.AddRouteRequest{
		Name:         r.RouteName,
		Paths:        paths,
		PreserveHost: false,
		StripPath:    false,
		Methods:      r.Methods,
	}
	body, _ := json.Marshal(data)
	request, _ := http.NewRequest(http.MethodPost, routeCfgUrl, bytes.NewBuffer(body))
	ks.doWithAuth(request)
}

func (ks kongService) removeAllRoutesFromService(serviceName string) {
	routeNames := ks.getAllRouteNamesForService(serviceName)
	for _, n := range routeNames {
		ks.removeRoute(n, serviceName)
	}

}

func (ks kongService) removeRoute(routeName string, serviceName string) {
	removeRouteUrl := fmt.Sprintf("%s/services/%s/routes/%s", ks.rootPath, serviceName, routeName)
	request, _ := http.NewRequest(http.MethodDelete, removeRouteUrl, nil)
	ks.doWithAuth(request)
}

func (ks kongService) getAllRouteNamesForService(serviceName string) (routesNames []string) {
	getAllRoutesUrl := fmt.Sprintf("%s/services/%s/routes", ks.rootPath, serviceName)
	request, _ := http.NewRequest(http.MethodGet, getAllRoutesUrl, nil)
	rawResponse := ks.doWithAuth(request)
	var routeResponse model.GetRouteResponse
	err := json.NewDecoder(rawResponse.Body).Decode(&routeResponse)
	LogOnError(err, "Error during parsing json")
	routesNames = make([]string, len(routeResponse.Data))
	for i, r := range routeResponse.Data {
		routesNames[i] = r.Name
	}
	return
}

func (ks kongService) doWithAuth(request *http.Request) (response *http.Response) {
	request.SetBasicAuth(ks.cfg.Global.KongUser, ks.kongAdminPassword)
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		LogOnError(err, "cannot process request")
	}
	LogToFile(request, response)
	return
}
