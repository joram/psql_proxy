FROM golang:1.20.5-alpine3.18 as builder
WORKDIR /go/src/github.com/joram/psql_proxy/

ENV GOPATH=/go/

COPY go.mod go.sum ./
RUN go mod download

ADD . .
ADD ./server/* ./server/*
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o app .


# RUNNING IMAGE
FROM alpine:3.14
COPY --from=builder /go/src/github.com/joram/psql_proxy/app /
ENTRYPOINT ["/app"]