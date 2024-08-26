FROM golang:1.23-bullseye AS build

WORKDIR /src

COPY go.* ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 go build -v -o go-app

FROM debian:bullseye-slim

RUN set -x && \
    apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=build /src/go-app /usr/bin/go-app

CMD ["/usr/bin/go-app"]
