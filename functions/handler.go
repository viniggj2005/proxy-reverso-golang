package functions

import (
	"embed"
	"fmt"
	"net/http"
	"proxy-reverso-golang/handlers"
	loadbalancers "proxy-reverso-golang/load_balancers"
	"proxy-reverso-golang/providers"
	"proxy-reverso-golang/structs"
	"strings"
)

//go:embed html/*
var htmlFiles embed.FS

func MeuHandler(writer http.ResponseWriter, request *http.Request) {
	redirects := providers.GetProxies()
	for _, redirect := range redirects {
		if strings.HasPrefix(request.URL.String(), redirect.Prefix) {
			balancer := providers.GetOrCreateBalancer(redirect.Prefix, func() loadbalancers.LoadBalancer {
				return getBalancer(redirect.LoadBalancer, redirect.Servers)
			})
			target := balancer.Next(redirect.Servers)
			if target == nil {
				render404(writer)
				return
			}
			target.Prefix = redirect.Prefix
			if request.Header.Get("Upgrade") == "websocket" {
				handlers.HandleWebSocket(writer, request, *target)
			} else if request.Header.Get("Content-Type") == "application/grpc" || request.Header.Get("Content-Type") == "application/grpc+proto" {
				handlers.HandleGrpc(writer, request, *target)
			} else {
				handlers.HandleHttp(writer, request, *target)
			}

			return
		}
	}
	render404(writer)
}

func render404(writer http.ResponseWriter) {
	conteudo404, err := htmlFiles.ReadFile("html/index.html")
	if err != nil {
		fmt.Println("\033[31m Erro ao ler o arquivo 404 embutido:\033[0m", err)
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte("404 - Not Found"))
		return
	}
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(http.StatusNotFound)
	writer.Write(conteudo404)
}

func getBalancer(balancer string, servers []structs.ServerConfigStruct) loadbalancers.LoadBalancer {
	switch balancer {
	case "default", "", "round-robin":
		return loadbalancers.NewRoundRobinBalancer()
	case "weighted-round-robin":
		return loadbalancers.NewWeightedRoundRobinBalancer(servers)
	case "random":
		return loadbalancers.NewRandomBalancer()
	default:
		return loadbalancers.NewRoundRobinBalancer()
	}
}
