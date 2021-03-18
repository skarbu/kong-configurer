package model

type Config struct {
	Connection ConnectionConfig
	Routing    []Routing `json:"routing"`
}

type Routing struct {
	ServiceName string            `json:"serviceName"`
	URL         string            `json:"url"`
	Routes      []AddRouteRequest `json:"routes"`
}

type ConnectionConfig struct {
	KongPath     string
	KongUser     string
	KongPassword string
}
