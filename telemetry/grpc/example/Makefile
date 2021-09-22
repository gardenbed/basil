.PHONY: tools
tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

.PHONY: zonePB
zonePB:
	protoc \
	  --proto_path=zonePB \
	  --go_out=zonePB --go_opt=paths=source_relative \
	  --go-grpc_out=zonePB --go-grpc_opt=paths=source_relative \
	  zonePB/zone.proto

.PHONY: server
server:
	@ cd server && go run .

.PHONY: client
client:
	@ cd client && go run .
