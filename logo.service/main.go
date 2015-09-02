package logo

import (
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
	"gopkg.in/olivere/elastic.v2"

	"github.com/moo-mou/maingo/redisManager"

	"bytes"
	"crypto/sha1"
	"strconv"
	"strings"

	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

const ServiceName = "logo"

const ErrorJSON = `{
	"err": "Unable to process request"
}`

type handler func(http.ResponseWriter, *http.Request)

func handleErr(res http.ResponseWriter, err error) {
	debug.PrintStack()
	rawJSON := json.RawMessage(ErrorJSON)
	reply, _ := rawJSON.MarshalJSON()
	res.Header().Set("Content-Type", "application/json")
	res.Write(reply)
	return
}

func getUrls(logos []string) []string {
	srs := make([]*elastic.SearchRequest, len(logos))
	for i, logo := range logos {
		srs[i] = NewLogoSearchRequest(logo, "")
	}
	return searchLogos(srs...)
}

func parseCanvasQuery(query string) ([]string, string) {
	categories := strings.Split(query, "|")
	logos := []string{}

	for _, cat := range categories {
		var all []string
		result := strings.Split(cat, ":")
		if len(result) == 2 {
			all = strings.Split(result[1], ",")
		} else if len(result) == 1 {
			all = strings.Split(result[0], ",")
		}
		logos = append(logos, all...)
	}

	var buffer bytes.Buffer
	h := sha1.New()
	for _, s := range logos {
		buffer.WriteString(s)
	}
	h.Write([]byte(buffer.String()))
	canvasKey := h.Sum(nil)

	return logos, fmt.Sprintf("%x", canvasKey)
}

func renderImg(logos []string, width string) string {
	urls := getUrls(logos)
	base64String := renderViaPJ(urls, width, true)
	return string(base64String)
}

func canvasHandler() handler {
	return func(res http.ResponseWriter, req *http.Request) {
		client := redisManager.GetClient()
		defer client.Close()

		canvas := req.URL.Query().Get("q")
		noCache := req.URL.Query().Get("noCache")
		width := req.URL.Query().Get("width")

		logos, canvasKey := parseCanvasQuery(canvas)

		var img string
		if noCache != "true" {
			img, _ = redis.String(client.Do("GET", fmt.Sprintf("%s:%s", "maingo", canvasKey)))
		}

		if img == "" {
			img = renderImg(logos, width)
			client.Do("SET", fmt.Sprintf("%s:%s", "maingo", canvasKey), img)
		}

		outputBytes, _ := b64.StdEncoding.DecodeString(img)
		res.Header().Set("Content-Type", "image/png")
		res.Header().Set("Content-Length", strconv.Itoa(len(outputBytes)))
		res.Write(outputBytes)
	}
}

func logoHandler() handler {
	return func(res http.ResponseWriter, req *http.Request) {
		q := req.URL.Query().Get("q")
		urls := getUrls([]string{q})
		if len(urls) != 0 {
			http.Redirect(res, req, urls[0], 302)
		} else {
			http.NotFound(res, req)
		}
	}
}

func Register(router *mux.Router) {
	router.HandleFunc("/canvas", canvasHandler()).Methods("GET")
	router.HandleFunc("/direct", logoHandler()).Methods("GET")
	log.Printf(" - %s registered\n", ServiceName)
}
