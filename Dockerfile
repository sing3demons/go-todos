#build stage
FROM golang:1.18.3-alpine3.16 AS builder

RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates
WORKDIR /go/src/app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

ENV CGO_ENABLED=0 
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -a -installsuffix cgo \
-ldflags="-w -s -X main.buildcommit=`git rev-parse --short HEAD` -X main.buildtime=`date "+%Y-%m-%dT%H:%M:%S%Z:00"`" \
  -o /go/bin/app 

#final stage
FROM alpine:3.16

WORKDIR /

COPY --from=builder /go/bin/app /app
COPY ./config/acl_model.conf ./config/acl_model.conf
COPY ./config/policy.csv ./config/policy.csv

RUN adduser -u 1001 -D -s /bin/sh -g ping 1001
RUN chown 1001:1001 /app

# RUN chmod +x /app
USER 1001

EXPOSE 8080
# CMD [ "/app" ]
ENTRYPOINT [ "/app" ]
