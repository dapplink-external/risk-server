risk-server:
	env GO111MODULE=on go build -v $(LDFLAGS) ./cmd/risk-server

clean:
	rm risk-server

test:
	go test -v ./...

lint:
	golangci-lint run ./...


.PHONY: \
	risk-server \
	clean \
	test \
	lint
