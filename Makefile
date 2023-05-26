.PHONY: build
build: server tools

.PHONY: fmt
fmt:
	gofmt -l -w `find . -type f -name '*.go' -not -path "./vendor/*"`
	goimports -l -w `find . -type f -name '*.go' -not -path "./vendor/*"`

.PHONY: local
local:
	LOGXI=* ./greenfield-deploy web

.PHONY: build
build:
	CGO_ENABLED=0 go build -v .

.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor

%:
	@:
