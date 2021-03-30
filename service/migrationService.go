package migrationService

import (
	"bytes"
	"io/ioutil"
	. "kong-configurer/logging"
	"kong-configurer/model"
	"net/http"
)

var (
	client = &http.Client{}
)

type migrationService struct {
	services []model.Service
	connCfg  model.ConnectionConfig
}

func ProcessMigration(migration *model.Migration, connCfg *model.ConnectionConfig) {
	ms := migrationService{
		services: migration.Services,
		connCfg:  *connCfg,
	}
	for _, s := range ms.services {
		ms.ensureService(s)
		ms.ensureRoutes(s.Routes, s.ServiceName)
		ms.ensurePlugin(s.Plugins, model.SERVICE, s.ServiceName)
		for _, r := range s.Routes {
			ms.ensurePlugin(r.Plugins, model.ROUTE, r.Name)
		}
	}
}

func (ms *migrationService) doWithAuth(request *http.Request) (response *http.Response) {
	requestBodyString := ""
	if request.Body != nil {
		requestBody, _ := ioutil.ReadAll(request.Body)
		request.Body.Close()
		request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
		requestBodyString = string(requestBody)
	}
	request.SetBasicAuth(ms.connCfg.KongUser, ms.connCfg.KongPassword)
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		LogOnError(err, "cannot process request")
	}
	LogToFile(request, response, requestBodyString)
	return
}
