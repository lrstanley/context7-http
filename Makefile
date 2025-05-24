.DEFAULT_GOAL := build

export PROJECT := "context7-http"
export PACKAGE := "github.com/lrstanley/context7-http"

license:
	curl -sL https://liam.sh/-/gh/g/license-header.sh | bash -s

clean:
	/bin/rm -rfv ${PROJECT}

fetch:
	go mod download
	go mod tidy

up:
	go get -u ./...
	go get -u -t ./...
	go mod tidy

prepare: fetch
	go generate -x ./...

inspect:
	@echo "server will be available at: http://localhost:8081/?transport=sse&serverUrl=http://localhost:8080/sse#tools"
	@CLIENT_PORT=8081 pnpm dlx @modelcontextprotocol/inspector

dlv: prepare
	dlv debug \
		--headless --listen=:2345 \
		--api-version=2 --log \
		--allow-non-terminal-interactive \
		${PACKAGE} -- --debug

debug: prepare
	go run ${PACKAGE} \
		--debug

build: prepare
	CGO_ENABLED=0 \
	go build \
		-ldflags '-d -s -w -extldflags=-static' \
		-tags=netgo,osusergo,static_build \
		-installsuffix netgo \
		-trimpath \
		-o ${PROJECT} \
		${PACKAGE}
