package structs

type ProxyConfigStruct struct {
	LoadBalancer string               `json:"loadBalancer"`
	ServerName   string               `json:"serverName"`
	Servers      []ServerConfigStruct `json:"servers"`
	Prefix       string               `json:"prefix"`
}

type ServerConfigStruct struct {
	Url       string `json:"url"`
	Weight    int    `json:"weight"`
	Available bool   `json:"available"`
}

type AllProxiesConfigStruct struct {
	Proxies []ProxyConfigStruct
}
