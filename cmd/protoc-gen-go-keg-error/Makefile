.PHONY: build

build:
	go build -o protoc-gen-go-keg-error main.go && sudo cp ./protoc-gen-go-keg-error /usr/local/bin/protoc-gen-go-keg-error \
	&& sudo cp ./protoc-gen-go-keg-error /usr/local/bin/protoc-gen-go-keg-error

api:
	protoc --proto_path=. --proto_path=../../third_party \
	   --go_out=paths=source_relative:. \
	   --go-keg-error_out=paths=source_relative:. \
	   errors.proto