.PHONY: build
build: server tools

.PHONY: fmt
fmt:
	gofmt -l -w `find . -type f -name '*.go' -not -path "./vendor/*"`
	goimports -l -w `find . -type f -name '*.go' -not -path "./vendor/*"`

.PHONY: local
local:
	LOGXI=* GITHUB_REPO=din-mukhammed GITHUB_TOKEN=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 ./greenfield-deploy web

.PHONY: bot
bot:
	cd bot && go run .

.PHONY: build
build:
	CGO_ENABLED=0 go build -v .

.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor

%:
	@:
