build: 
	go build -o ./bin/market

run: build
	./bin/market

test: 
	go test -v ./...