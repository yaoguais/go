proto:
	protoc --go_out=plugins=grpc:. *.proto

serve:
	go run main.go main.pb.go

test:
	go run test.go main.pb.go