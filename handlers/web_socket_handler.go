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

	// tls.Dial("tcp", connectionUrl.Host, &tls.Config{InsecureSkipVerify: true})
	dialConnection, err := connectToBackend(redirect)
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor:", err)
		return
	}
	forwardHanshake(request, redirect, dialConnection)

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
	tunnel(hijackConnection, dialConnection)

}

func tunnel(hijackConnection net.Conn, dialConnection net.Conn) {
	go io.Copy(hijackConnection, dialConnection)
	io.Copy(dialConnection, hijackConnection)
}

func forwardHanshake(request *http.Request, redirect structs.Redirects, dialConnection net.Conn) error {
	request.URL.Host = getHost(redirect)
	request.URL.Path = strings.TrimPrefix(request.URL.Path, redirect.Prefix)
	request.URL.Scheme = "ws"
	request.RequestURI = ""

	err := request.Write(dialConnection)
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor:", err)
		return err
	}
	return nil
}

func connectToBackend(redirect structs.Redirects) (net.Conn, error) {
	dialConnection, err := net.Dial("tcp", getHost(redirect))
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor:", err)
		return nil, err
	}
	return dialConnection, nil
}

func getHost(redirect structs.Redirects) string {
	connectionUrl, err := url.Parse(redirect.Url)
	if err != nil {
		fmt.Println("Erro ao fazer parse da URL:", err)
		return ""
	}
	return connectionUrl.Host
}
