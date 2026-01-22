package functions

import (
	"fmt"
	"net/http"
	"proxy-reverso-golang/global"
	"time"
)

func StartHealthCheck() {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for range ticker.C {
			checkHealth()
		}
	}()
}

func checkHealth() {
	client := http.Client{
		Timeout: 2 * time.Second,
	}

	global.ProxyMutex.Lock()
	defer global.ProxyMutex.Unlock()

	for proxyIndex, proxy := range global.ProxiesConfig.Proxies {
		for serverIndex, server := range proxy.Servers {
			available := true
			response, err := client.Head(server.Url)
			if err != nil || response.StatusCode >= 400 {
				available = false
			}
			if response != nil {
				response.Body.Close()
			}

			if available != server.Available {
				status := "ONLINE"
				if !available {
					status = "OFFLINE"
				}
				fmt.Printf("\033[33m[HealthCheck]\033[0m Servidor %s agora está %s\n", server.Url, status)

				global.ProxiesConfig.Proxies[proxyIndex].Servers[serverIndex].Available = available

				global.BalancerMutex.Lock()
				delete(global.LoadBalancers, proxy.Prefix)
				global.BalancerMutex.Unlock()
			}
		}
	}
}
