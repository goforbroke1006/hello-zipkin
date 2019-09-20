package main

import (
	"encoding/json"
	"github.com/openzipkin/zipkin-go"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
)

type OrderClient struct {
	baseUrl string
	client  *zipkinhttp.Client
}

type BookResponse struct {
	Status string    `json:"status"`
	Time   time.Time `json:"time"`
}

func (pc *OrderClient) Book(amount int, span zipkin.Span) *BookResponse {
	url := pc.baseUrl + "/order"

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Printf("Can't create register http request: %s", err)
	}

	resp, err := pc.client.Do(
		req.WithContext(zipkin.NewContext(req.Context(), span)),
	)
	if err != nil {
		log.Printf("Can't send register request: %s", err)
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Can't read register response: %s", err)
	}

	res := &BookResponse{}
	_ = json.Unmarshal(respBody, res)

	return res
}

func NewOrderClient(baseUrl string, client *zipkinhttp.Client) *OrderClient {
	return &OrderClient{
		baseUrl: baseUrl,
		client:  client,
	}
}
