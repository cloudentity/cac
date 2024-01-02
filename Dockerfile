FROM golang:1.21

RUN go install github.com/cloudentity/cac@dev

ENTRYPOINT ["/go/bin/cac"]
