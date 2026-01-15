package handlers

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httputil"
	"proxy-reverso-golang/shared"
	"proxy-reverso-golang/structs"
	"strings"

	"golang.org/x/net/http2"
)

func HandleGrpc(writer http.ResponseWriter, request *http.Request, redirect structs.Redirects) {
	transport := getTransport(redirect.Url)
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = "http"
			if shared.VerifyTlsConnection(redirect.Url) {
				req.URL.Scheme = "https"
			}
			req.URL.Host = getHost(redirect)
			req.Header.Set("TE", "trailers")
			req.URL.Path = strings.TrimPrefix(req.URL.Path, redirect.Prefix)
		},
		FlushInterval: -1,
		Transport:     transport,
	}
	proxy.ServeHTTP(writer, request)
}

func getTransport(url string) *http2.Transport {
	isTls := shared.VerifyTlsConnection(url)
	if isTls {
		return &http2.Transport{}
	}
	return &http2.Transport{
		AllowHTTP: true,
		DialTLS: func(network, address string, config *tls.Config) (net.Conn, error) {
			return net.Dial(network, address)
		},
	}
}
