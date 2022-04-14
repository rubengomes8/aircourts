FROM golang:1.18 as gobuilder

COPY . /app
WORKDIR /app
RUN GOGC=off go build -mod=vendor -a -installsuffix cgo -o ./find_courts cmd/http/main.go

FROM debian:buster
RUN apt-get update \
    && apt-get install -y locales \
    && apt-get install -y wget curl \
    && apt-get install -y ca-certificates \
    && apt-get install -y xsltproc \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

COPY --from=gobuilder /app/find_courts /find_courts
COPY --from=gobuilder /app/cron.sh /cron.sh
COPY --from=gobuilder /app/configuration /configuration

RUN chmod u+x ./find_courts

ENTRYPOINT ["./cron.sh"]

# docker build -f Dockerfile -t aircourts .
# docker run -p 5100:80 aircourts