package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	. "kong-configurer/logging"
	"kong-configurer/model"
	"kong-configurer/service"
	"os"
	"syscall"
)

func main() {
	cfg := loadConfigFromFile()
	kongPassword := getPasswordPrompt(cfg.Global.KongUser)
	kongService := service.NewKongService(cfg, kongPassword)
	fmt.Print("processing... \n")
	kongService.ProcessMigration()
	fmt.Print("done")
}

func loadConfigFromFile() (cfg model.Config) {
	if len(os.Args) != 2 {
		FailOnError(errors.New("illegal number of parameters"), "Check a syntax")
	}
	file, err := ioutil.ReadFile(os.Args[1])
	FailOnError(err, "Can't reach configuration, ensure filename is correct")
	err = json.Unmarshal(file, &cfg)
	FailOnError(err, "Invalid config file structure")
	errs := cfg.Validate()
	FailOnErrors(errs, "Invalid config file structure")
	return
}

func getPasswordPrompt(usr string) string {
	fmt.Printf("Enter password for kong user: '%s' \n", usr)
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	FailOnError(err, "error during reading password")
	return string(bytePassword)
}
