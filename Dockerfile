FROM golang:alpine as builder

ENV GOPATH=/

WORKDIR /app

COPY cmd cmd
COPY internal internal
COPY go.mod .
COPY go.sum .

RUN go get ./...
RUN go build cmd/main.go

FROM alpine:latest

COPY --from=builder ./app/main .

CMD ["./main"]

