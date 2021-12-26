.PHONY: build postgres

build:
		go build -ldflags "-X main.buildcommit=`git rev-parse --short HEAD` \
		-X main.buildtime=`date "+%Y-%m-%dT%H:%M:%S%Z:00"`" \
		-o app

database-up:
			docker compose -f database/docker-compose.yml up -d
			docker compose up redis -d
clean-up:
			docker compose -f database/docker-compose.yml down
			docker compose down

image:	
		docker build -t gotodos:0.0.1 .

container:
		docker run --rm -p 8080:8080 -e PORT=8080 -e APP_ENV=dev -e DB_USER=postgres -e DB_PASSWORD=passw0rd \
			-e DB_NAME=todos -e DB_PORT=5432 -e DB_HOST=host.docker.internal -e REDIS_HOST=host.docker.internal \
			-v $(shell pwd)/logs:/logs -v $(shell pwd)/uploads:/uploads --name gotodos gotodos:0.0.1
