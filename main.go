package main

import (
	"github.com/duck8823/webhook-proxy/proxy/handlers"
	"github.com/google/logger"
	"net/http"
	"os"
)

func main() {
	logger.Init("webhook-proxy", false, false, os.Stdout)

	http.Handle("/", &handlers.DangerOnDocker{})

	http.ListenAndServe(":8080", nil)
}
