package model

type AddServiceRequest struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type AddRouteRequest struct {
	Name         string   `json:"name"`
	Paths        []string `json:"paths"`
	PreserveHost bool     `json:"preserve_host"`
	StripPath    bool     `json:"strip_path"`
	Methods      []string `json:"methods"`
}

type GetRoutesResponse struct {
	Data []KongRouteModel `json:"data"`
	Next string           `json:"next"`
}

type KongRouteModel struct {
	ID        string   `json:"id"`
	CreatedAt int      `json:"created_at"`
	UpdatedAt int      `json:"updated_at"`
	Name      string   `json:"name"`
	Protocols []string `json:"protocols"`
	Methods   []string `json:"methods,omitempty"`
	Hosts     []string `json:"hosts,omitempty"`
	Paths     []string `json:"paths,omitempty"`
	Headers   struct {
		XAnotherHeader []string `json:"x-another-header"`
		XMyHeader      []string `json:"x-my-header"`
	} `json:"headers,omitempty"`
	HTTPSRedirectStatusCode int      `json:"https_redirect_status_code"`
	RegexPriority           int      `json:"regex_priority"`
	StripPath               bool     `json:"strip_path"`
	PathHandling            string   `json:"path_handling"`
	PreserveHost            bool     `json:"preserve_host"`
	Tags                    []string `json:"tags"`
	Service                 struct {
		ID string `json:"id"`
	} `json:"service"`
	Snis    []string `json:"snis,omitempty"`
	Sources []struct {
		IP   string `json:"ip,omitempty"`
		Port int    `json:"port,omitempty"`
	} `json:"sources,omitempty"`
	Destinations []struct {
		IP   string `json:"ip,omitempty"`
		Port int    `json:"port,omitempty"`
	} `json:"destinations,omitempty"`
}

type KongServiceModel struct {
	ID                string   `json:"id"`
	CreatedAt         int      `json:"created_at"`
	UpdatedAt         int      `json:"updated_at"`
	Name              string   `json:"name"`
	Retries           int      `json:"retries"`
	Protocol          string   `json:"protocol"`
	Host              string   `json:"host"`
	Port              int      `json:"port"`
	Path              string   `json:"path"`
	ConnectTimeout    int      `json:"connect_timeout"`
	WriteTimeout      int      `json:"write_timeout"`
	ReadTimeout       int      `json:"read_timeout"`
	Tags              []string `json:"tags"`
	ClientCertificate struct {
		ID string `json:"id"`
	} `json:"client_certificate"`
}
