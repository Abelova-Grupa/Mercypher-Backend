PROTO_DIR = proto

GATEWAY_PROTO_FILES = proto/gateway/api-gateway.proto
OUT_GATEWAY = proto

SESSION_PROTO_FILES = proto/session/session-service.proto
OUT_SESSION = proto

MESSAGE_PROTO_FILES = proto/message/message-service.proto
OUT_MESSAGE = proto

USER_PROTO_FILES = proto/user/user-service.proto
OUT_USER = proto

.PHONY: proto 

# Make proto runs all services, Make gateway only runs gateway
proto: gateway session user message

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
		--go_out=$(OUT_SESSION) \
		--go-grpc_out=$(OUT_SESSION) \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		$(SESSION_PROTO_FILES)

message:
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(OUT_MESSAGE) \
		--go-grpc_out=$(OUT_MESSAGE) \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		$(MESSAGE_PROTO_FILES)

user:
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(OUT_USER) \
		--go-grpc_out=$(OUT_USER) \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		$(USER_PROTO_FILES)
