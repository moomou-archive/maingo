package logo

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"gopkg.in/olivere/elastic.v2"

	es "github.com/moo-mou/maingo/elastic"
)

const INDEX_NAME = "mango"
const TYPE_NAME = "logos"

const BASE64_PNG_PREFIX = "data:image/png;base64,"
const DEFAULT_PHANTOM_HOST = "http://localhost:3333"

type PJRequest struct {
	Urls      []string `json:"urls"`
	Watermark bool     `json:"watermark"`
	Width     string   `json:"width"`
}

type LogoResult struct {
	S3link string `json:"s3link"`
}

func NewLogoSearchRequest(simple, exact string) *elastic.SearchRequest {
	var q elastic.Query
	if exact != "" {
		tq := elastic.NewTermFilter("simple", simple)
		mq := elastic.NewMatchQuery("exact", exact)
		q = elastic.NewFilteredQuery(mq)
		// type assertion
		q.(elastic.FilteredQuery).Filter(tq)
	} else {
		q = elastic.NewTermQuery("simple", simple)
	}
	return elastic.NewSearchRequest().Index(INDEX_NAME).Type(TYPE_NAME).
		Source(elastic.NewSearchSource().Query(q))
}

func searchLogos(srs ...*elastic.SearchRequest) []string {
	var urls []string = nil

	results, err := es.MultiSearch(srs...)

	if err != nil {
		return nil
	}

	for _, sres := range results.Responses {
		for _, hit := range sres.Hits.Hits {
			var logoResult LogoResult
			_, err := json.Marshal(hit.Source)
			if err != nil {
				panic(err)
			}
			if err := json.Unmarshal(*hit.Source, &logoResult); err != nil {
				log.Println("ERR:", err)
				return nil
			}
			urls = append(urls, logoResult.S3link)
		}
	}

	return urls
}

func renderViaPJ(urls []string, width string, watermark bool) []byte {
	phantomJSHost := os.Getenv("PHANTOM_SERVER")
	if phantomJSHost == "" {
		phantomJSHost = DEFAULT_PHANTOM_HOST
	}

	request := PJRequest{
		Urls:      urls,
		Width:     width,
		Watermark: watermark,
	}
	jsonStr, _ := json.Marshal(&request)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", phantomJSHost, bytes.NewBuffer(jsonStr))
	r.Header.Add("Content-Type", "application/json")
	resp, _ := client.Do(r)

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body
}
