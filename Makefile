all:
	go build -o ./bin/synapse-backregister ./cmd/synapse-backregister

docker:
	docker build -t euank/synapse-backregister:latest .
