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
	GOGC=off go build -mod=vendor -a -installsuffix cgo -o ./bin/find_courts	cmd/http/main.go

run-finder:
	./bin/find_courts

run-finder-email:
	./bin/find_courts -email=true

run-finder-outdoor:
	./bin/find_courts -indoor=false
