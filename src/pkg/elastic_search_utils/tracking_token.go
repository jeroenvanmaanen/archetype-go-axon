package elastic_search_utils

import (
	log "log"
	strconv "strconv"
	strings "strings"

	json "encoding/json"
	ioutil "io/ioutil"

	elasticSearch7 "github.com/elastic/go-elasticsearch/v7"
)

type ElasticSearchTokenStore struct {
	ProcessorName string
	ES7           *elasticSearch7.Client
}

func OpenTokenStore(processorName string) *ElasticSearchTokenStore {
	es7 := WaitForElasticSearch()
	return &ElasticSearchTokenStore{
		ProcessorName: processorName,
		ES7:           es7,
	}
}

func (tokenStore *ElasticSearchTokenStore) ReadToken() *int64 {
	es7 := tokenStore.ES7
	response, e := es7.Get("tracking-token", tokenStore.ProcessorName)
	if e != nil {
		log.Printf("Elastic search: Error while reading token: %v", e)
		return nil
	}
	log.Printf("Elastic search: token document: %v", response)

	if response.StatusCode == 404 {
		return nil
	}

	responseJson, e := ioutil.ReadAll(response.Body)
	if e != nil {
		log.Printf("Elastic search: Error while reading response body: %v", e)
		return nil
	}

	jsonMap := make(map[string]interface{})
	e = json.Unmarshal(responseJson, &jsonMap)
	if e != nil {
		log.Printf("Elastic search: Error while unmarshalling JSON, %v", e)
		return nil
	}

	hexToken := jsonMap["_source"].(map[string]interface{})["token"].(string)
	token, e := strconv.ParseInt(hexToken, 16, 64)
	if e != nil {
		log.Printf("Elastic search: Error while parsing hex token, %v", e)
		return nil
	}
	log.Printf("Elastic search: token: %v", token)

	return &token
}

func (tokenStore *ElasticSearchTokenStore) WriteToken(token int64) error {
	var b strings.Builder
	b.WriteString(`{"token" : "`)
	b.WriteString(strconv.FormatInt(token, 16))
	b.WriteString(`"}`)
	if e := AddToIndex("tracking-token", tokenStore.ProcessorName, b.String(), tokenStore.ES7); e != nil {
		log.Printf("Elastic search: Error while storing tracking token: %v: %v", tokenStore.ProcessorName, e)
		return e
	}
	return nil
}
