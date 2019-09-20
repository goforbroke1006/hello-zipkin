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

type PaymentClient struct {
	baseUrl string
	client  *zipkinhttp.Client
}

type DepositResponse struct {
	Status string    `json:"status"`
	Time   time.Time `json:"time"`
}

func (pc *PaymentClient) Deposit(amount int, span zipkin.Span) *DepositResponse {
	url := pc.baseUrl + "/payment/deposit"

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Printf("Can't create register http request: %s", err)
	}

	//span := zipkin.SpanFromContext(r.Context())
	//span.Tag("custom_key", "some value")

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

	res := &DepositResponse{}
	_ = json.Unmarshal(respBody, res)

	return res
}

func NewPaymentClient(baseUrl string, client *zipkinhttp.Client) *PaymentClient {
	return &PaymentClient{
		baseUrl: baseUrl,
		client:  client,
	}
}
