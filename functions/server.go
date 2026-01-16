package functions

import (
	"fmt"
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
	fmt.Println("\033[32mServer starting on port 8080\033[0m")

	if err := server.ListenAndServe(); err != nil {
		fmt.Println("\033[31mError starting server:\033[0m", err)
	}
}
