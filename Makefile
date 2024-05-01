export GO111MODULE=on
export GOBIN=$(CURDIR)/bin
export BUF_BIN=$(GOBIN)/buf

run:
	go run ./cmd/authorization_server

generate: bin-deps deps vendor-proto
	$(BUF_BIN) generate --path=./api/authorization_server

bin-deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.0
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.19.1
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.19.1
	go install github.com/bufbuild/buf/cmd/buf@v1.31.0

	# добавляем скаченные бинарники в PATH
	export PATH=$(PATH):$(go env GOPATH)/bin

deps:
	go mod download

# TODO: найти или написать нормальную вендорилку прото
vendor-proto:
	mkdir -p vendor.proto/google/api
	curl https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto -o vendor.proto/google/api/annotations.proto
	curl https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto -o vendor.proto/google/api/http.proto
	mkdir -p vendor.proto/google/protobuf
	curl https://raw.githubusercontent.com/protocolbuffers/protobuf/main/src/google/protobuf/descriptor.proto -o vendor.proto/google/protobuf/descriptor.proto
