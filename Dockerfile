## Build
FROM golang:1.22.3-alpine3.18 AS buildenv

ADD go.mod go.sum /

RUN go mod download

WORKDIR /app

ADD . .

ENV GO111MODULE=on


RUN  go build -o main cmd/main.go

## Deploy
FROM alpine

WORKDIR /

COPY --from=buildenv  /app/ /

EXPOSE 3000

CMD ["/main"]