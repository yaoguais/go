proto:
	protoc -I/usr/local/include \
		-I. \
  		-I${GOPATH}/src \
  		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  		--go_out=plugins=grpc:. \
  		main.proto

	protoc -I/usr/local/include \
		-I. \
		-I${GOPATH}/src  -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--grpc-gateway_out=logtostderr=true:. \
		main.proto

# make serve_grpc
# make serve_http
# make test

serve_grpc:
	go run grpc_main.go main.pb.go main.pb.gw.go

serve_http:
	go run http_main.go main.pb.go main.pb.gw.go

test:
	curl -X POST -k http://localhost:50052/v1/example/echo -d '{"name": " world"}'