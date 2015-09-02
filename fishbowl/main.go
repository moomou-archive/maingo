package main

import (
	"fmt"
	"github.com/facebookgo/grace/gracehttp"
	"net/http"
	"os"
)

const DEFAULT_CONFIG_PATH = "./_config/service.json"

var (
	serviceConfigPath = "" + os.Getenv("CONFIG_PATH")
	serverPort        = ":" + os.Getenv("PORT")
)

func main() {
	if serviceConfigPath == "" {
		serviceConfigPath = DEFAULT_CONFIG_PATH
	}
	if serverPort == ":" {
		serverPort = "localhost:9990"
	}

	sc, err := ParseServiceConfig(&serviceConfigPath)

	if err != nil {
		panic(err)
	}

	handler := RegisterServices(sc)

	fmt.Printf("Starting server on %s...\n", serverPort)

	gracehttp.Serve(
		&http.Server{Addr: serverPort, Handler: handler},
	)
}
