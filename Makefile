tip:
	rm -f ./bin/tip
	go build -o ./bin/tip ./cmd/tip/main.go
test:
	go test -cover -race ./...
build:
	go build -o ./bin/tip ./cmd/tip/main.go
clean:
	rm -f ./bin/*
install:
	make build
	sudo mv ./bin/tip /usr/local/bin/