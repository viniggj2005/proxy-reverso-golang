package structs

type ProxyConfigStruct struct {
	LoadBalancer string   `json:"loadBalancer"`
	Urls         []string `json:"urls"`
	Prefix       string   `json:"prefix"`
}

type AllProxiesConfigStruct struct {
	Proxies []ProxyConfigStruct
}
