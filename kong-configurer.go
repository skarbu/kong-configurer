package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	. "kong-configurer/logging"
	"kong-configurer/model"
	"kong-configurer/service"
	"syscall"
)

var (
	configFileParam string
	userNameParam   string
	passwordParam   string
	hostNameParam   string
)

func init() {
	flag.StringVar(&configFileParam, "f", "", "relative path to json config file")
	flag.StringVar(&userNameParam, "u", "", "kong user with admin privileges")
	flag.StringVar(&passwordParam, "p", "", "password for kong user")
	flag.StringVar(&hostNameParam, "h", "", "kong host with port")
	flag.Parse()
}
func main() {
	cfg := loadConfigFromFile()
	cfg.Connection = loadConnectionCfgFromArgs()
	cfg.Connection.KongPassword = getPassword(cfg.Connection.KongUser)
	errs := cfg.Validate()
	FailOnErrors(errs, "Invalid config file structure")
	kongService := service.NewKongService(cfg)
	fmt.Print("processing... \n")
	kongService.ProcessMigration()
	fmt.Print("done")
}

func loadConnectionCfgFromArgs() model.ConnectionConfig {
	return model.ConnectionConfig{
		KongUser: userNameParam,
		KongPath: hostNameParam,
	}
}

func loadConfigFromFile() (cfg model.Config) {
	if configFileParam == "" {
		FailOnError(errors.New("config filename cannot be empty"), "Check a syntax")
	}
	file, err := ioutil.ReadFile(configFileParam)
	FailOnError(err, "Can't reach configuration, ensure filename is correct")
	err = json.Unmarshal(file, &cfg)
	FailOnError(err, "Invalid config file structure")
	return
}

func getPassword(usr string) string {
	if passwordParam == "" {
		fmt.Printf("Enter password for kong user: '%s' \n", usr)
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		FailOnError(err, "error during reading password")
		passwordParam = string(bytePassword)
	}
	return passwordParam
}
