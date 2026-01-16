package main

import (
	"fmt"
	"proxy-reverso-golang/functions"
	"proxy-reverso-golang/global"
)

func main() {
	functions.GetConfig()
	for _, proxy := range global.ProxiesConfig.Proxies {
		fmt.Println(proxy.Prefix)
	}
	// functions.ServerInit()
}
