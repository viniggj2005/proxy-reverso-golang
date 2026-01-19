package functions

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type ConfigStruct struct {
	ReadTimeoutSeconds       int `json:"readTimeoutSeconds"`
	WriteTimeoutSeconds      int `json:"writeTimeoutSeconds"`
	IdleTimeoutSeconds       int `json:"idleTimeoutSeconds"`
	ReadHeaderTimeoutSeconds int `json:"readHeaderTimeoutSeconds"`
	MaxHeaderMB              int `json:"maxHeaderMB"`
}

func (config ConfigStruct) HttpServerInit() {
	http2Server := &http2.Server{}
	server := &http.Server{
		Addr:              ":80",
		Handler:           h2c.NewHandler(http.HandlerFunc(MeuHandler), http2Server),
		ReadTimeout:       time.Duration(config.ReadTimeoutSeconds) * time.Second,
		WriteTimeout:      time.Duration(config.WriteTimeoutSeconds) * time.Second,
		IdleTimeout:       time.Duration(config.IdleTimeoutSeconds) * time.Second,
		ReadHeaderTimeout: time.Duration(config.ReadHeaderTimeoutSeconds) * time.Second,
		MaxHeaderBytes:    config.MaxHeaderMB * 1 << 20,
	}
	fmt.Println("\033[32mServer starting on port 80\033[0m")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("\033[31mError starting server:\033[0m", err)
	}
}

func (config ConfigStruct) HttpsServerInit() {
	http2Server := &http2.Server{}
	serverTLS := &http.Server{
		Addr:              ":443",
		Handler:           h2c.NewHandler(http.HandlerFunc(MeuHandler), http2Server),
		ReadTimeout:       time.Duration(config.ReadTimeoutSeconds) * time.Second,
		WriteTimeout:      time.Duration(config.WriteTimeoutSeconds) * time.Second,
		IdleTimeout:       time.Duration(config.IdleTimeoutSeconds) * time.Second,
		ReadHeaderTimeout: time.Duration(config.ReadHeaderTimeoutSeconds) * time.Second,
		MaxHeaderBytes:    config.MaxHeaderMB * 1 << 20,
	}
	fmt.Println("\033[32mServer starting on port 443\033[0m")
	if err := serverTLS.ListenAndServeTLS("cert.pem", "key.pem"); err != nil {
		fmt.Println("\033[31mError starting server:\033[0m", err)
	}
}
