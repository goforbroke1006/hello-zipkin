package main

import (
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	httpreporter "github.com/openzipkin/zipkin-go/reporter/http"
)

var (
	argHandleAddr = flag.String("handle-addr", "localhost:9003", "")
)

const serviceName = "order-svc"

func main() {
	rand.Seed(time.Now().UnixNano())

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

	router := mux.NewRouter()
	router.Methods("POST").Path("/order").HandlerFunc(applyOrderHandler)
	_ = http.ListenAndServe(*argHandleAddr, serverMiddleware(router))
}

func applyOrderHandler(w http.ResponseWriter, req *http.Request) {
	latency := 500 + rand.Intn(100)
	time.Sleep(time.Duration(latency) * time.Millisecond)

	result := map[string]interface{}{
		"status": "IN_PROCESS",
		"time":   time.Now().UTC().Unix(),
	}

	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Println(err.Error())
		return
	}
}
