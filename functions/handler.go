package functions

import (
	"fmt"
	"net/http"
	"os"
	"proxy-reverso-golang/handlers"
	"proxy-reverso-golang/structs"
	"strings"
)

func MeuHandler(writer http.ResponseWriter, request *http.Request) {
	redirects := []structs.Redirects{
		{Prefix: "/api", Url: "http://localhost:50051"},
	}
	for _, redirect := range redirects {
		if strings.HasPrefix(request.URL.String(), redirect.Prefix) {

			if request.Header.Get("Upgrade") == "websocket" {
				handlers.HandleWebSocket(writer, request, redirect)
			} else if request.Header.Get("Content-Type") == "application/grpc" || request.Header.Get("Content-Type") == "application/grpc+proto" {
				handlers.HandleGrpc(writer, request, redirect)
			} else {
				handlers.HandleHttp(writer, request, redirect)
			}

			return
		}
	}
	render404(writer)
}

func render404(writer http.ResponseWriter) {
	conteudo404, err := os.ReadFile("./html/index.html")
	if err != nil {
		fmt.Println("Erro ao ler o arquivo 404:", err)
		return
	}
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(http.StatusNotFound)
	writer.Write(conteudo404)
}
