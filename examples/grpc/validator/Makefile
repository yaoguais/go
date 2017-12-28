proto:
	protoc  \
		-I${GOPATH}/src \
		--proto_path=${GOPATH}/src \
		--proto_path=${GOPATH}/src/github.com/google/protobuf/src \
		--proto_path=. \
		--go_out=plugins=grpc:. \
		--govalidators_out=. \
		*.proto

serve:
	go run main.go main.pb.go main.validator.pb.go
