.PHONY: build postgres

build:
		go build -ldflags "-X main.buildcommit=`git rev-parse --short HEAD` \
		-X main.buildtime=`date "+%Y-%m-%dT%H:%M:%S%Z:00"`" \
		-o app

database up:
			docker compose up db redis -d
clean:
			docker compose down

image:	
		docker build -t gotodos:0.0.1 .

container:
		docker run --rm -p 8080:8080 -e PORT=8080 -e APP_ENV=dev -e DSN='host=host.docker.internal user=postgres password=passw0rd dbname=todos port=5432  sslmode=disable TimeZone=Asia/Bangkok' \
		 -e REDIS_HOST=host.docker.internal -v $(shell pwd)/logs:/logs -v $(shell pwd)/uploads:/uploads --name gotodos sing3demons/gotodos:0.0.5
