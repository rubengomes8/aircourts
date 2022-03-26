PROJECT_PATH=github.com/rubengomes8/aircourts
PROJECT_NAME=aircourts

# DEPENDENCIES
install:
	go mod init ${PROJECT_PATH}
	@echo "=== Installing dependencies ==="
	go mod vendor
	@echo "Done"

update:
	@echo "=== Installing dependencies ==="
	go mod tidy
	go mod vendor
	@echo "Done"

clean:
	rm -f "go.mod"
	rm -f "go.sum"
	rm -rfv vendor/*


# BUILD
build:
	GOGC=off go build -mod=vendor -a -installsuffix cgo -o ./bin/scrapper	cmd/http/main.go

run-http:
	go run cmd/http/main.go
