package elastic_search_utils

import (
	context "context"
	errors "errors"
	log "log"
	strings "strings"
	time "time"

	json "encoding/json"

	elasticSearch7 "github.com/elastic/go-elasticsearch/v7"
	esapi "github.com/elastic/go-elasticsearch/v7/esapi"
)

func WaitForElasticSearch() *elasticSearch7.Client {
	cfg := elasticSearch7.Config{
		Addresses: []string{
			"http://elastic-search:9200",
		},
	}
	d, _ := time.ParseDuration("3s")
	log.Printf(".(es)")
	for true {
		es7, e := elasticSearch7.NewClient(cfg)
		if e == nil {
			if e = getElasticSearchInfo(es7); e == nil {
				waitForStatusReady(es7)
				return es7
			}
		}
		time.Sleep(d)
		log.Printf(".(es) %v", e)
	}
	return nil
}

func waitForStatusReady(es7 *elasticSearch7.Client) {
	var status string
	d, _ := time.ParseDuration("3s")
	for true {
		status = "???"
		statsRequest := esapi.ClusterStatsRequest{}
		statsResponse, e := statsRequest.Do(context.Background(), es7)
		if e == nil {
			stats, e := UnwrapElasticSearchResponse(statsResponse)
			if e == nil {
				status = stats["status"].(string)
				if status != "red" {
					log.Printf("Elastic Search: status: %v", status)
					return
				}
			}
		}

		time.Sleep(d)
		log.Printf(".(es %v)", status)
	}
}

func getElasticSearchInfo(es7 *elasticSearch7.Client) error {
	info, e := es7.Info(es7.Info.WithContext(context.Background()), es7.Info.WithHuman())
	if e != nil {
		log.Printf("Error while requesting Elastic Search info: %v", e)
		return e
	}
	log.Printf("Elastic Search: info: %v", info)

	result, e := UnwrapElasticSearchResponse(info)
	if e != nil {
		return e
	}

	// Print client and server version numbers.
	log.Printf("Client: %s", elasticSearch7.Version)
	log.Printf("Server: %s", result["version"].(map[string]interface{})["number"])
	log.Println(strings.Repeat("~", 37))
	return nil
}

func UnwrapElasticSearchResponse(response *esapi.Response) (map[string]interface{}, error) {
	if response.IsError() {
		log.Printf("Error: %s", response.String())
		return nil, errors.New("Elastic Search error: " + response.String())
	}

	if response.Body == nil {
		log.Printf("Missing body")
		return nil, errors.New("elastic Search error: response has no body")
	}

	defer response.Body.Close()
	// Deserialize the response into a map.
	var result map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Printf("Error parsing the response body: %s", err)
		return nil, err
	}
	return result, nil
}

func AddToIndex(indexName string, id string, body string, es7 *elasticSearch7.Client) error {
	log.Printf("Add to index: %s: Document ID: %v", indexName, id)

	// Set up the request object.
	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: id,
		Body:       strings.NewReader(body),
		Refresh:    "true",
	}

	// Perform the request with the client.
	res, err := req.Do(context.Background(), es7)
	if err != nil {
		log.Printf("Elastic Search: Error getting response: %s", err)
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Elastic Search: [%s] Error indexing document ID=%s", res.Status(), id)
		return errors.New("error indexing document")
	}

	// Deserialize the response into a map.
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Printf("Elastic Search: Error parsing the response body: %s", err)
		return errors.New("error parsing response body")
	}

	// Print the response status and indexed document version.
	log.Printf("Elastic Search: [%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
	return nil
}
