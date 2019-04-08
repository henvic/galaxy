.SILENT: server test
.PHONY: server test
server:
	go run ./cmd/server
test:
	./scripts/test.sh
