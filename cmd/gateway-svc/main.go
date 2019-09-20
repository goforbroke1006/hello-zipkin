package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	httpreporter "github.com/openzipkin/zipkin-go/reporter/http"
)

var (
	argHandleAddr = flag.String("handle-addr", "localhost:8080", "")
)

const serviceName = "gateway-svc"

func main() {

	reporter := httpreporter.NewReporter("http://localhost:9411/api/v2/spans")
	defer reporter.Close()

	endpoint, err := zipkin.NewEndpoint(serviceName, *argHandleAddr)
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}

	// initialize our tracer
	tracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		log.Fatalf("unable to create tracer: %+v\n", err)
	}

	// create global zipkin http server middleware
	serverMiddleware := zipkinhttp.NewServerMiddleware(
		tracer, zipkinhttp.TagResponseSize(true),
	)

	// create global zipkin traced http client
	client, err := zipkinhttp.NewClient(tracer, zipkinhttp.ClientTrace(true))
	if err != nil {
		log.Fatalf("unable to create client: %+v\n", err)
	}

	paymentClient := NewPaymentClient("http://localhost:9002", client)
	orderClient := NewOrderClient("http://localhost:9003", client)

	// initialize router
	router := mux.NewRouter()
	router.Methods("POST").Path("/buy").HandlerFunc(payOrderHandler(tracer, paymentClient, orderClient))
	_ = http.ListenAndServe(*argHandleAddr, serverMiddleware(router))

}

func payOrderHandler(tracer *zipkin.Tracer, pc *PaymentClient, oc *OrderClient) (func(http.ResponseWriter, *http.Request)) {
	return func(w http.ResponseWriter, req *http.Request) {

		parentSpan := tracer.StartSpan("booking")

		deposit := pc.Deposit(123456, parentSpan)
		if deposit != nil {
			oc.Book(123, parentSpan)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("Success"))
		} else {
			w.WriteHeader(http.StatusPaymentRequired)
			_, _ = w.Write([]byte("Need payment"))
		}

		parentSpan.Finish()

	}
}
