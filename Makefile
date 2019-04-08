.SILENT: server apidocs write-apidocs test
.PHONY: server apidocs write-apidocs test
server:
	go run ./cmd/server
write-apidocs:
	swag init --generalInfo cmd/server/main.go
apidocs:
	go run -tags=swagon ./cmd/server
test:
	./scripts/test.sh
