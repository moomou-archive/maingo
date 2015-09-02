package elastic

import (
	"gopkg.in/olivere/elastic.v2"

	"fmt"
	"log"
	"os"
)

var client *elastic.Client

func Index(indexName string, typeName string, id string, data interface{}) error {
	res, err := client.Index().
		Index(indexName).
		Type(typeName).
		Id(id).
		BodyJson(data).
		Do()
	log.Printf("INDEX: %s to index %s, type %s\n", res.Id, res.Index, res.Type)
	return err
}

func Search(indexName string, typeName string, fieldName string, fieldValue string) (*elastic.SearchResult, error) {
	matchQuery := elastic.NewMatchQuery(fieldName, fieldValue)
	matchQuery = matchQuery.Fuzziness("AUTO")
	result, err := client.Search().Index(indexName).Type(typeName).Query(&matchQuery).Pretty(true).Do()
	log.Printf("Found a total of %d %s\n", result.TotalHits(), typeName)
	return result, err
}

func SearchQuery(indexName string, typeName string, q elastic.Query) (*elastic.SearchResult, error) {
	result, err := client.Search().Index(indexName).Type(typeName).Query(q).Do()
	log.Printf("Found a total of %d %s\n", result.TotalHits(), typeName)
	return result, err
}

func MultiSearch(srs ...*elastic.SearchRequest) (*elastic.MultiSearchResult, error) {
	return client.MultiSearch().
		Add(srs...).
		Do()
}

func init() {
	var err error
	host := os.Getenv("ELASTIC_ADDR")
	port := os.Getenv("ELASTIC_PORT")
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "9200"
	}
	url := fmt.Sprintf("http://%s:%s", host, port)
	log.Printf("Trying to connect to %s", url)
	client, err = elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(false))
	if err != nil {
		log.Fatal(err)
	}
}
