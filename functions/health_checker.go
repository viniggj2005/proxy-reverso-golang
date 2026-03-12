package functions

import (
	"fmt"
	"net/http"
	"proxy-reverso-golang/providers"
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

	for proxyIndex, proxy := range providers.GetProxies() {
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

				providers.SetProxyAvailability(available, proxyIndex, serverIndex)

				providers.DeleteLoadBalancer(proxy.Prefix)
			}
		}
	}
}
