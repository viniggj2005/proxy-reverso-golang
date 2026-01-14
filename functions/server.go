package functions

import (
	"net/http"
	"time"
)

func ServerInit() {
	server := &http.Server{
		Addr:         ":8080",
		Handler:      http.HandlerFunc(MeuHandler),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		println("Error starting server:", err.Error())
	}
}
