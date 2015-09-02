package main

import (
	"github.com/gorilla/mux"
	logo "github.com/moo-mou/maingo/logo.service"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type register func(*mux.Router)

type ServiceConfig struct {
	Name   string `json:name`
	Active bool   `json:bool`
}

var (
	availableServices = map[string]register{
		"logo": logo.Register,
	}
)

func ParseServiceConfig(path *string) ([]ServiceConfig, error) {
	var sc []ServiceConfig
	fmt.Printf("sc path:%s\n", *path)
	content, _ := ioutil.ReadFile(*path)
	fmt.Printf("sc content:%s\n", string(content))
	err := json.Unmarshal(content, &sc)
	return sc, err
}

func RegisterServices(sc []ServiceConfig) *mux.Router {
	router := mux.NewRouter()
	services := make([]string, len(availableServices))

	for _, s := range sc {
		if s.Active {
			serviceName := fmt.Sprintf("/%s", s.Name)
			services = append(services, serviceName)
			log.Printf("Registering %s\n", serviceName)
			serviceRouter := router.PathPrefix(serviceName).Subrouter()
			availableServices[s.Name](serviceRouter)
		}
	}

	router.HandleFunc("/index", func(res http.ResponseWriter, req *http.Request) {
		result, _ := json.Marshal(services)
		res.Write(result)
	})

	return router
}
