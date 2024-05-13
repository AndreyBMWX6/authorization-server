export GO111MODULE=on
export GOBIN=$(CURDIR)/bin
export BUF_BIN=$(GOBIN)/buf

LOCAL_DB_NAME:=authorization-server

run: run-auth run-client

run-auth:
	go run ./cmd/authorization_server

run-client:
	go run ./cmd/oauth2_client

generate: bin-deps deps vendor-proto
	$(BUF_BIN) generate --path=./api/authorization_server
	$(BUF_BIN) generate --path=./api/authentication_server

env-up:
	docker-compose up -d

env-down:
	docker-compose down

db-reset: db-create db-up

db-create:
	psql -U postgres -c "drop database if exists \"$(LOCAL_DB_NAME)\""
	psql -U postgres -c "create database \"$(LOCAL_DB_NAME)\""

db-up:
	goose -dir migrations postgres "postgres://postgres:postgres@localhost:5432/${LOCAL_DB_NAME}?sslmode=disable" up

jet:
	$(GOBIN)/jet -dsn "postgres://postgres:postgres@localhost:5432/${LOCAL_DB_NAME}?sslmode=disable" -path=./internal/generated -ignore-tables=goose_db_version

#todo: add installing docker and docker-compose
bin-deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.0
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.19.1
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.19.1
	go install github.com/bufbuild/buf/cmd/buf@v1.31.0
	go install github.com/pressly/goose/v3/cmd/goose@v3.20.0
	go install github.com/go-jet/jet/v2/cmd/jet@v2.11.1

	# добавляем скаченные бинарники в PATH
	export PATH=$(PATH):$(go env GOPATH)/bin

deps:
	go mod download
	go mod tidy

# TODO: найти или написать нормальную вендорилку прото
vendor-proto:
	mkdir -p vendor.proto/google/api
	curl https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto -o vendor.proto/google/api/annotations.proto
	curl https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto -o vendor.proto/google/api/http.proto
	mkdir -p vendor.proto/google/protobuf
	curl https://raw.githubusercontent.com/protocolbuffers/protobuf/main/src/google/protobuf/descriptor.proto -o vendor.proto/google/protobuf/descriptor.proto
