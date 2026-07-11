FROM golang:1.25-alpine3.24 as builder

RUN apk add --no-cache make ca-certificates gcc musl-dev linux-headers git jq bash

COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum

WORKDIR /app

RUN go mod download

ARG CONFIG=config.yml

# build risk-server with the shared go.mod & go.sum files
COPY . /app/risk-server

WORKDIR /app/risk-server

RUN make

FROM alpine:3.18

COPY --from=builder /app/risk-server/risk-server /usr/local/bin

WORKDIR /app

ENTRYPOINT ["risk-server"]
