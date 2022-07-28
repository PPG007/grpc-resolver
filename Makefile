start-server:
	go run server/server.go

start-client:
	go run client/client.go

gen-cert:
	cd cert && \
	./gen.sh
