package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	. "kong-configurer/logging"
	"kong-configurer/model"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		FailOnError(errors.New("illegal number of parameters"), "Check a syntax")
	}
	file, err := ioutil.ReadFile(os.Args[1])
	FailOnError(err, "Can't reach configuration, ensure filename is correct")
	var cfg model.Config
	err = json.Unmarshal(file, &cfg)
	FailOnError(err, "Invalid config file structure")
	errs := cfg.Validate()
	FailOnErrors(errs, "Invalid config file structure")
	//todo iterate over config and configure routs
}
