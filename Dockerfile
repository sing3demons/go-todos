#build stage
FROM golang:1.18.3-alpine3.16 AS builder
WORKDIR /go/src/app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

ENV GOARCH=amd64
RUN go build -ldflags "-X main.buildcommit=`git rev-parse --short HEAD` \ 
    -X main.buildtime=`date "+%Y-%m-%dT%H:%M:%S%Z:00"`" \
    -o /go/bin/app 

#final stage
FROM alpine:3.16

WORKDIR /

COPY --from=builder /go/bin/app /app

COPY ./config/acl_model.conf ./config/acl_model.conf
COPY ./config/policy.csv ./config/policy.csv

LABEL Name=gotodos Version=0.0.1
EXPOSE 8080
CMD [ "/app" ]

