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
	migrationService "kong-configurer/service"
	"syscall"
	"time"
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
	defer TimeTrack(time.Now(), "migration")
	migration := loadMigrationFromFile()
	connCfg := loadConnCfg()
	fmt.Print("processing... \n")
	migrationService.ProcessMigration(&migration, &connCfg)
	fmt.Print("done\n")
}

func loadConnCfg() model.ConnectionConfig {
	connCfg := model.ConnectionConfig{
		KongUser: userNameParam,
		KongPath: hostNameParam,
	}
	connCfg.KongPassword = loadPassword(connCfg.KongUser)
	errs := connCfg.Validate()
	FailOnErrors(errs, "Invalid config file structure")
	return connCfg
}

func loadMigrationFromFile() (migration model.Migration) {
	if configFileParam == "" {
		FailOnError(errors.New("config filename cannot be empty"), "Check a syntax")
	}
	file, err := ioutil.ReadFile(configFileParam)
	FailOnError(err, "Can't reach configuration, ensure filename is correct")
	err = json.Unmarshal(file, &migration)
	FailOnError(err, "Invalid config file structure")
	return
}

func loadPassword(usr string) string {
	if passwordParam == "" {
		fmt.Printf("Enter password for kong user: '%s' \n", usr)
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		FailOnError(err, "error during reading password")
		passwordParam = string(bytePassword)
	}
	return passwordParam
}
