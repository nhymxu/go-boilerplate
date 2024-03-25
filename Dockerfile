FROM golang:1.22-buster AS build

WORKDIR /src

COPY go.* ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 go build -v -o go-app

FROM debian:buster-slim

RUN set -x && \
    apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=build /src/go-app /usr/bin/go-app

CMD ["/usr/bin/go-app"]
