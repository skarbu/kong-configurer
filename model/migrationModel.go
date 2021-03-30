package model

import (
	"reflect"
)

type Migration struct {
	Services []Service `json:"routing"`
}

type Service struct {
	ServiceName string   `json:"serviceName"`
	URL         string   `json:"url"`
	Routes      []Route  `json:"routes"`
	Plugins     []Plugin `json:"plugins"`
}

type Route struct {
	AddRouteRequest
	Plugins []Plugin `json:"plugins"`
}

type Plugin struct {
	Name      *string      `json:"name"`
	Consumer  *string      `json:"consumer"`
	Protocols *[]string    `json:"protocols"`
	Config    *interface{} `json:"config"`
}

func (p1 Plugin) Equals(p2 KongPluginsResponseModel) bool {
	if *p1.Name == *p2.Name &&
		consumerEquals(p1, p2) &&
		reflect.DeepEqual(*p1.Config, *p2.Config) {
		return true
	}
	return false
}

func consumerEquals(c1 Plugin, c2 KongPluginsResponseModel) bool {
	if c1.Consumer == nil && c2.Consumer == nil {
		return true
	}
	if c1.Consumer != nil && c2.Consumer == nil {
		return false
	}
	if *c1.Consumer == c2.Consumer.ID {
		return true
	}
	return false
}
