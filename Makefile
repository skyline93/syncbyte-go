all: syncbyte-engine syncbyte-agent syncbyte

docker-image: clean syncbyte-engine syncbyte-agent syncbyte
	docker rmi syncbyte:latest
	docker build -t syncbyte:latest .

syncbyte-engine:
	CGO_ENABLED=0 go build -o output/syncbyte-engine cmd/engine/main.go

syncbyte-agent:
	CGO_ENABLED=0 go build -o output/syncbyte-agent cmd/agent/main.go

syncbyte:
	CGO_ENABLED=0 go build -o output/syncbyte cmd/syncbyte/*

clean:
	rm -rf logs data output
