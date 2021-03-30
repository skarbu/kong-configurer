package migrationService

import (
	"bytes"
	"encoding/json"
	"fmt"
	. "kong-configurer/logging"
	"kong-configurer/model"
	"net/http"
)

func (ms *migrationService) ensurePlugin(plugins []model.Plugin, scope string, scopeId string) {
	existingPluginsInService := ms.getExistingPlugins(scope, scopeId)
	for _, p := range *existingPluginsInService {
		if !containsByName(plugins, p) {
			ms.removePlugin(scope, scopeId, *p.ID)
		}
	}
	for _, p := range plugins {
		contains, modified, pluginId := containsPlugin(existingPluginsInService, p)
		switch contains {
		case false:
			ms.addPlugin(scope, scopeId, p)
		case true:
			if modified {
				ms.modifyPlugin(scope, scopeId, p, *pluginId)
			} else {
				LogMsgToFile(fmt.Sprintf("Plugin `%s` for `%s`: `%s` already exist and does not changed", *p.Name, scope, scopeId))
			}
		}
	}
}

func containsPlugin(plugins *[]model.KongPluginsResponseModel, plugin model.Plugin) (exist bool, modified bool, pluginId *string) {
	for _, actualPlugin := range *plugins {
		if *plugin.Name == *actualPlugin.Name {
			if plugin.Equals(actualPlugin) {
				return true, false, actualPlugin.ID
			}
			return true, true, actualPlugin.ID
		}
	}
	return false, true, nil
}

func containsByName(plugins []model.Plugin, plugin model.KongPluginsResponseModel) bool {
	for _, p := range plugins {
		if *p.Name == *plugin.Name {
			return true
		}
	}
	return false
}

func (ms *migrationService) removePlugin(scope string, scopeId string, pluginId string) {
	removePluginUrl := fmt.Sprintf("%s/%s/%s/plugins/%s", ms.connCfg.KongPath, scope, scopeId, pluginId)
	request, _ := http.NewRequest(http.MethodDelete, removePluginUrl, nil)
	rawResponse := ms.doWithAuth(request)
	var serviceResponse model.PluginsResponse
	err := json.NewDecoder(rawResponse.Body).Decode(&serviceResponse)
	LogOnError(err, "Error during parsing json")
}

func (ms *migrationService) addPlugin(scope string, scopeId string, plugin model.Plugin) {
	addPluginUrl := fmt.Sprintf("%s/%s/%s/plugins", ms.connCfg.KongPath, scope, scopeId)
	addPluginRequest, err := plugin.ToAddPluginRequest(scope, scopeId)
	FailOnError(err, "error during adding plugins")
	body, _ := json.Marshal(addPluginRequest)
	request, _ := http.NewRequest(http.MethodPost, addPluginUrl, bytes.NewBuffer(body))
	ms.doWithAuth(request)
}

func (ms *migrationService) modifyPlugin(scope string, scopeId string, plugin model.Plugin, pluginId string) {
	addPluginUrl := fmt.Sprintf("%s/plugins/%s", ms.connCfg.KongPath, pluginId)
	addPluginRequest, err := plugin.ToAddPluginRequest(scope, scopeId)
	FailOnError(err, "error during adding plugins")
	body, _ := json.Marshal(addPluginRequest)
	request, _ := http.NewRequest(http.MethodPatch, addPluginUrl, bytes.NewBuffer(body))
	ms.doWithAuth(request)
}

func (ms *migrationService) getExistingPlugins(scope string, scopeId string) *[]model.KongPluginsResponseModel {
	getPluginsUrl := fmt.Sprintf("%s/%s/%s/plugins", ms.connCfg.KongPath, scope, scopeId)
	request, _ := http.NewRequest(http.MethodGet, getPluginsUrl, nil)
	rawResponse := ms.doWithAuth(request)
	var serviceResponse model.PluginsResponse
	err := json.NewDecoder(rawResponse.Body).Decode(&serviceResponse)
	LogOnError(err, "Error during parsing json")
	return &serviceResponse.Data
}
