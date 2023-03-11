build:
	go build -o bin/exchange

run: build	
	./bin/exchange http-server --http-shutdown-delay=1s

test:
	go test -v ./...