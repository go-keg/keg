keg.build: keg.build
protoc-gen-go-keg-error.build: protoc-gen-go-keg-error.build

lint:
	@clear
	golangci-lint run -c .golangci.yml --fix

%.gen:
	$(eval SERVICE:= $*)
	go generate ./cmd/$(SERVICE)/main.go

%.build:
	$(eval SERVICE:= $*)
	@echo "build: $(SERVICE):$(GIT_VERSION)"
	go env -w CGO_ENABLED=0 GOOS=linux GOARCH=amd64
	go build -ldflags "-X main.Version=$(GIT_VERSION)" -o ./bin/$(SERVICE) ./cmd/$(SERVICE)/
	#go build -ldflags "-s -w -X main.Version=$(GIT_VERSION)" -o ./bin/$(SERVICE) ./cmd/$(SERVICE)/