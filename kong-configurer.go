package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	. "kong-configurer/logging"
	"kong-configurer/model"
	"net/http"
	"os"
	"time"
)

var (
	logFileName = "log_" + timestamp() + ".txt"
	cfg         model.Config
	rootPath    string
	client      = &http.Client{}
)

func main() {
	if len(os.Args) != 2 {
		FailOnError(errors.New("illegal number of parameters"), "Check a syntax")
	}
	file, err := ioutil.ReadFile(os.Args[1])
	FailOnError(err, "Can't reach configuration, ensure filename is correct")
	err = json.Unmarshal(file, &cfg)
	FailOnError(err, "Invalid config file structure")
	errs := cfg.Validate()
	FailOnErrors(errs, "Invalid config file structure")
	rootPath = fmt.Sprintf("%s:%d", cfg.Config.KongHost, cfg.Config.KongPort)
	for _, r := range cfg.Routing {
		ensureService(r.ServiceName, r.URL)
		ensureRoutes(r.Routes, r.ServiceName)
	}
}

func ensureService(serviceName string, url string) {
	serviceCfgUrl := fmt.Sprintf("%s/services", rootPath)
	data := model.ServiceRequest{
		Name: serviceName,
		Url:  url,
	}
	body, _ := json.Marshal(data)
	request, _ := http.NewRequest(http.MethodPost, serviceCfgUrl, bytes.NewBuffer(body))
	doWithAuth(request)
}

func ensureRoutes(routes []model.Route, serviceName string) {
	for _, r := range routes {
		routeCfgUrl := fmt.Sprintf("%s/services/%s/routes", rootPath, serviceName)
		paths := make([]string, 1)
		paths[0] = r.Path
		data := model.RouteRequest{
			Name:         r.RouteName,
			Paths:        paths,
			PreserveHost: false,
			StripPath:    false,
			Methods:      r.Methods,
		}
		body, _ := json.Marshal(data)
		request, _ := http.NewRequest(http.MethodPost, routeCfgUrl, bytes.NewBuffer(body))
		doWithAuth(request)
	}
}

func doWithAuth(request *http.Request) {
	//TODO password should be read while app start from console
	request.SetBasicAuth(cfg.Config.KongUser, "password")
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		LogError(err, "cannot process request")
		return
	}
	logToFile(response)
}

func logToFile(response *http.Response) {
	body, _ := ioutil.ReadAll(response.Body)
	ioutil.WriteFile(logFileName, body, fs.ModeAppend)
	newLine := bytes.NewBufferString("\n").Bytes()
	ioutil.WriteFile(logFileName, newLine, fs.ModeAppend)
}

func timestamp() string {
	now := time.Now()
	return fmt.Sprintf("%d%d%d%d%d%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
}
