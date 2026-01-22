package main

import (
	"proxy-reverso-golang/functions"
)

func main() {
	config, _ := functions.GetMainConfig("main.json", "")
	functions.GetConfig()
	go functions.WatchConfigs()

	functions.StartHealthCheck()

	if config.HttpsOn {
		go config.HttpsServerInit()
	}
	config.HttpServerInit()
}
