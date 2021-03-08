package model

type ServiceRequest struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type RouteRequest struct {
	Name         string   `json:"name"`
	Paths        []string `json:"paths"`
	PreserveHost bool     `json:"preserve_host"`
	StripPath    bool     `json:"strip_path"`
	Methods      []string `json:"methods"`
}
