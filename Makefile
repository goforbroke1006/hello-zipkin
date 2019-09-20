#source-run:
#	go run cmd/gateway-svc/main.go &
#	go run cmd/order-svc/main.go &
#	go run cmd/payment-svc/main.go &

send-test:
	curl -X POST -H "Content-Type: application/json" --url 'http://localhost:8080/buy' --data '{"product_id": 123}'--verbose