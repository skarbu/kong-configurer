package migrationService

import (
	"bytes"
	"encoding/json"
	"fmt"
	. "kong-configurer/converter"
	. "kong-configurer/logging"
	"kong-configurer/model"
	"net/http"
)

func (ms *migrationService) ensureService(r model.Service) {
	exist, modified := ms.checkService(r)
	switch exist {
	case true:
		if modified {
			ms.modifyService(r)
		} else {
			LogMsgToFile(fmt.Sprintf("Service `%s` already exist and does not changed", r.ServiceName))
		}
	case false:
		ms.addService(r)
	}
}

func (ms *migrationService) addService(r model.Service) {
	addServiceUrl := fmt.Sprintf("%s/services", ms.connCfg.KongPath)
	data := model.AddServiceRequest{
		Name: r.ServiceName,
		Url:  r.URL,
	}
	body, _ := json.Marshal(data)
	request, _ := http.NewRequest(http.MethodPost, addServiceUrl, bytes.NewBuffer(body))
	ms.doWithAuth(request)
}

func (ms *migrationService) modifyService(r model.Service) {
	modifyServiceUrl := fmt.Sprintf("%s/services/%s", ms.connCfg.KongPath, r.ServiceName)
	data := model.AddServiceRequest{
		Name: r.ServiceName,
		Url:  r.URL,
	}
	body, _ := json.Marshal(data)
	request, _ := http.NewRequest(http.MethodPatch, modifyServiceUrl, bytes.NewBuffer(body))
	ms.doWithAuth(request)
}

func (ms *migrationService) getService(routing model.Service) (service *model.KongServiceResponseModel) {
	getServiceUrl := fmt.Sprintf("%s/services/%s", ms.connCfg.KongPath, routing.ServiceName)
	request, _ := http.NewRequest(http.MethodGet, getServiceUrl, nil)
	rawResponse := ms.doWithAuth(request)
	var serviceResponse model.KongServiceResponseModel
	err := json.NewDecoder(rawResponse.Body).Decode(&serviceResponse)
	LogOnError(err, "Error during parsing json")
	return &serviceResponse
}

func (ms *migrationService) checkService(requestService model.Service) (exist bool, modify bool) {
	existingService := ms.getService(requestService)
	if existingService.ID == nil {
		return false, true
	}
	if isServiceModify(existingService, requestService) {
		return true, true
	}
	return true, false
}

func isServiceModify(existingService *model.KongServiceResponseModel, requestedService model.Service) bool {
	existingURL := fmt.Sprintf("%v://%v:%v%v",
		StrToStr(existingService.Protocol),
		StrToStr(existingService.Host),
		IntToStr(existingService.Port),
		StrToStr(existingService.Path))
	return requestedService.URL != existingURL
}
