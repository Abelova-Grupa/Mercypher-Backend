PROTO_DIR = proto

OUT_GATEWAY = api-gateway/internal/grpc
GATEWAY_PROTO_FILES = proto/api-gateway.proto

OUT_SESSION = session-service/internal/grpc/pb
SESSION_PROTO_FILES = proto/session-service.proto

.PHONY: proto
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

