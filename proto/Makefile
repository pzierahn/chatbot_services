## Shipping package
## Assumes user is on MacOS, if other OS, please change PROTO_ROOT_DIR to the path of protobuf installation
PROTO_ROOT_DIR = $(shell brew --prefix)/Cellar/protobuf/27.1/include
APP_DIR = /Users/patrick/Code/brainboost_app/lib/generated

## Dart requires you to manually ship all google provided proto files too.
dart:
	@mkdir -p $(APP_DIR)
	@protoc -I=. --dart_out=grpc:$(APP_DIR) *.proto
	@protoc -I$(PROTO_ROOT_DIR) --dart_out=$(APP_DIR) $(PROTO_ROOT_DIR)/google/protobuf/*.proto

go:
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		--include_source_info --include_imports \
		--descriptor_set_out=api_descriptor.pb \
		*.proto

gen: dart go