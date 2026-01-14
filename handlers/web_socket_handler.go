package handlers

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"proxy-reverso-golang/structs"
	"strings"
)

func HandleWebSocket(writer http.ResponseWriter, request *http.Request, redirect structs.Redirects) {
	connectionUrl, err := url.Parse(redirect.Url)
	if err != nil {
		fmt.Println("Erro ao fazer parse da URL:", err)
		return
	}
	fmt.Println(connectionUrl.Host)
	dialConnection, err := net.Dial("tcp", connectionUrl.Host)
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor:", err)
		return
	}
	request.URL.Host = connectionUrl.Host
	request.URL.Path = strings.TrimPrefix(request.URL.Path, redirect.Prefix)
	request.URL.Scheme = "ws"
	err = request.Write(dialConnection)
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor:", err)
		return
	}
	hijack, ok := writer.(http.Hijacker)
	if !ok {
		fmt.Println("Erro ao conectar ao servidor:", err)
		return
	}
	hijackConnection, _, err := hijack.Hijack()
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor:", err)
		return
	}
	go io.Copy(hijackConnection, dialConnection)
	io.Copy(dialConnection, hijackConnection)

}
