package model

type Config struct {
	Global struct {
		KongHost string `json:"kongHost"`
		KongPort int    `json:"kongPort"`
		KongUser string `json:"kongUser"`
	} `json:"config"`
	Routing []Routing `json:"routing"`
}

type Routing struct {
	ServiceName string            `json:"serviceName"`
	URL         string            `json:"url"`
	Routes      []AddRouteRequest `json:"routes"`
}
