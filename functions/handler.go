package functions

import (
	"fmt"
	"net/http"
	"os"
	"proxy-reverso-golang/global"
	"proxy-reverso-golang/handlers"
	loadbalancers "proxy-reverso-golang/load_balancers"
	"strings"
)

func MeuHandler(writer http.ResponseWriter, request *http.Request) {
	redirects := global.ProxiesConfig.Proxies
	fmt.Println("Recebendo Request:", request.URL.String())
	fmt.Println("Proxies carregados:", len(redirects))
	for _, redirect := range redirects {
		fmt.Println("Verificando prefixo:", redirect.Prefix, "para URL:", request.URL.String())
		if strings.HasPrefix(request.URL.String(), redirect.Prefix) {
			global.BalancerMutex.Lock()
			balancer, exists := global.LoadBalancers[redirect.Prefix]
			if !exists {
				balancer = loadbalancers.NewRoundRobinBalancer()
				global.LoadBalancers[redirect.Prefix] = balancer
			}
			global.BalancerMutex.Unlock()

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
	conteudo404, err := os.ReadFile("./html/index.html")
	if err != nil {
		fmt.Println("\033[31m Erro ao ler o arquivo 404:\033[0m", err)
		return
	}
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(http.StatusNotFound)
	writer.Write(conteudo404)
}
