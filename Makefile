all: syncbyte-engine syncbyte-agent syncbyte

docker-image: clean syncbyte-engine syncbyte-agent syncbyte
	docker rmi -f syncbyte:latest
	docker build -t syncbyte:latest .

rpm: clean syncbyte-engine syncbyte-agent syncbyte
	rm -rf ~/rpmbuild

	rpmbuild -bb deploy/agent.spec
	rpmbuild -bb deploy/engine.spec
	rpmbuild -bb deploy/syncbyte.spec

	cp ~/rpmbuild/RPMS/x86_64/*.rpm output
	rm -rf ~/rpmbuild

syncbyte-engine:
	CGO_ENABLED=0 go build -o output/syncbyte-engine cmd/engine/main.go

syncbyte-agent:
	CGO_ENABLED=0 go build -o output/syncbyte-agent cmd/agent/main.go

syncbyte:
	CGO_ENABLED=0 go build -o output/syncbyte cmd/syncbyte/*

clean:
	rm -rf logs data output
