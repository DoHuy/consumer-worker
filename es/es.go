package es

import (
	"context"
	"github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/esapi"
	"net/http"
	"strings"
	"time"
)

type ElasticSearch struct {
	Client 	*elasticsearch.Client
	Index	string
}

func New(esUris, index string) (*ElasticSearch, error){
	addrs := strings.Split(esUris, ",")
	cfg := elasticsearch.Config{
		Addresses: addrs,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 10,
			ResponseHeaderTimeout: time.Second,
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return & ElasticSearch{
		Client: es,
		Index: index,
	}, err
}

func (c *ElasticSearch) Insert(ctx context.Context, body string) error{
	// Instantiate a request object
	req := esapi.IndexRequest{
		Index:      c.Index,
		Body:       strings.NewReader(body),
		Refresh:    "true",
	}
	// Return an API response object from request
	res, err := req.Do(ctx, c.Client)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}