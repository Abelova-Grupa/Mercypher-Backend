PROTO_DIR = proto
OUT_GATEWAY = api-gateway/internal/grpc

PROTO_FILES = $(wildcard $(PROTO_DIR)/*.proto)

.PHONY: proto
gateway:
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(OUT_GATEWAY) \
		--go-grpc_out=$(OUT_GATEWAY) \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)
