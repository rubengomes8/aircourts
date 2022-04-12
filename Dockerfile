FROM golang:1.18 as gobuilder

COPY . /app

WORKDIR /app

RUN GOGC=off go build -mod=vendor -a -installsuffix cgo -o ./bin/find_courts cmd/http/main.go

RUN chmod u+x ./bin/find_courts

ENTRYPOINT ["./bin/cron.sh"]