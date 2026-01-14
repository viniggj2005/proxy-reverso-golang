package handlers

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"proxy-reverso-golang/structs"
	"strings"
	"time"
)

var hipByHop = [8]string{
	"Te",
	"Upgrade",
	"Trailers",
	"Connection",
	"Keep-Alive",
	"Transfer-Encoding",
	"Proxy-Authenticate",
	"Proxy-Authorization",
}

func HandleHttp(writer http.ResponseWriter, request *http.Request, redirect structs.Redirects) {
	transport := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}
	preparedRequest, err := prepareProxyRequest(request, redirect)
	if err != nil {
		http.Error(writer, "Erro ao fazer request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := client.Do(preparedRequest)
	if err != nil {
		http.Error(writer, "Erro ao fazer request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	copyHeader(writer, response.Header)

	writer.WriteHeader(response.StatusCode)

	flusher, ok := writer.(http.Flusher)
	if !ok {
		io.Copy(writer, response.Body)
		return
	}

	buffer := make([]byte, 1024)
	for {
		n, err := response.Body.Read(buffer)
		if n > 0 {
			writer.Write(buffer[:n])
			flusher.Flush()
		}
		if err != nil {
			if err != io.EOF {
				fmt.Println("Erro ao ler resposta:", err)
			}
			break
		}
	}
}

func isHipByHop(header string) bool {
	for _, banned := range hipByHop {
		if header == banned {
			return true
		}
	}
	return false

}

func copyHeader(writer http.ResponseWriter, header http.Header) {
	for key, values := range header {
		if !isHipByHop(key) {
			for _, value := range values {
				writer.Header().Add(key, value)
			}
		}
	}
}

func prepareProxyRequest(request *http.Request, redirect structs.Redirects) (*http.Request, error) {
	suffix := strings.TrimPrefix(request.URL.Path, redirect.Prefix)
	targetUrl, err := url.Parse(redirect.Url)
	if err != nil {
		fmt.Println("Erro ao fazer parse da URL:", err)
		return nil, err
	}
	clientIP, _, err := net.SplitHostPort(request.RemoteAddr)
	if err != nil {
		clientIP = request.RemoteAddr
	}

	prior := request.Header.Get("X-Forwarded-For")
	if prior != "" {
		clientIP = prior + ", " + clientIP
	}
	request.Header.Set("X-Forwarded-For", clientIP)

	request.Header.Set("Host", targetUrl.Host)
	request.Host = targetUrl.Host
	targetUrl.Path += suffix
	request.URL = targetUrl
	request.RequestURI = ""

	for _, header := range hipByHop {
		request.Header.Del(header)
	}
	return request, nil
}
