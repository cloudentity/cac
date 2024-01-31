FROM golang:1.21 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY scripts/entrypoint.sh /entrypoint.sh

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cac .

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/cac .
COPY --from=build /entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]