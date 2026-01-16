package structs

type ConfigStruct struct {
	Port                int `json:"port"`
	ReadTimeoutSeconds  int `json:"readTimeoutSeconds"`
	WriteTimeoutSeconds int `json:"writeTimeoutSeconds"`
}

type ProxyConfigStruct struct {
	LoadBalancer string   `json:"loadBalancer"`
	Urls         []string `json:"urls"`
	Prefix       string   `json:"prefix"`
}

type AllProxiesConfigStruct struct {
	Proxies []ProxyConfigStruct
}
