FROM golang:1.21 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cac .

FROM alpine:latest

WORKDIR /app

RUN apk add bash

COPY --from=build /app/cac .
COPY /scripts/action-entrypoint.sh /action-entrypoint.sh

ENTRYPOINT ["/app/cac"]