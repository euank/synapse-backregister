all:
	CGO_ENABLED=0 go build -a -tags netgo -ldflags '-w' -o ./bin/synapse-backregister ./cmd/synapse-backregister

docker:
	docker build -t euank/synapse-backregister:latest .
