PROTO_DIR = proto

GATEWAY_PROTO_FILES = proto/api-gateway.proto
OUT_GATEWAY = api-gateway/external/grpc

SESSION_PROTO_FILES = proto/session-service.proto
OUT_SESSION = session-service/internal/grpc/pb

MESSAGE_PROTO_FILES = proto/message-service.proto
OUT_MESSAGE = message-service/external/grpc

.PHONY: proto 

# Make proto runs all services, Make gateway only runs gateway
proto: gateway session

gateway:
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(OUT_GATEWAY) \
		--go-grpc_out=$(OUT_GATEWAY) \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		$(GATEWAY_PROTO_FILES)

session:
	protoc \
		--proto_path=$(PROTO_DIR) \
		--proto_path=googleapis \
		--go_out=$(OUT_SESSION) \
		--go-grpc_out=$(OUT_SESSION) \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=$(OUT_SESSION) \
  		--grpc-gateway_opt=paths=source_relative \
		$(SESSION_PROTO_FILES)

message:
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(OUT_MESSAGE) \
		--go-grpc_out=$(OUT_MESSAGE) \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		$(MESSAGE_PROTO_FILES)


