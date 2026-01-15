package functions

import (
	"net/http"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func ServerInit() {
	http2Server := &http2.Server{}
	server := &http.Server{
		Addr:         ":8080",
		Handler:      h2c.NewHandler(http.HandlerFunc(MeuHandler), http2Server),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		println("Error starting server:", err.Error())
	}
}
