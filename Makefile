GOLANG_VERSION := 1.17.8

test:
	docker run --rm -it -v "${PWD}:/app" -w /app golang:${GOLANG_VERSION} go test -v .

lint:
	docker run --rm -v "${PWD}:/app" -w /app golangci/golangci-lint:v1.43.0 golangci-lint run -v .
