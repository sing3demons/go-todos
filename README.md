```
go run main.go
```

#### Liveness Probe

```
cat /tmp/live
echo $?
```

`output = 0 ,status = up`

| Method |  RelativePath     | CURL                              |
| ------ | ------------     | --------------------------------- |
| GET    | /healthz         | ` http://127.0.0.1:8080/healthz` |
| GET    | /x               | ` http://127.0.0.1:8080/x`       |
| GET    | /api/v1/todos           | ` http://127.0.0.1:8080/api/v1/todos`       |
| GET    | /api/v1/todos/:id           | ` http://127.0.0.1:8080/api/v1/todos/:id`       |
| DELETE    | /api/v1/todos/:id           | ` http://127.0.0.1:8080/api/v1/todos/:id`       |
| POST    | /api/v1/todos           | ` http://127.0.0.1:8080/api/v1/todos`       |

```
docker build -t gotodos:0.0.1 .
docker run --rm -p 8080:8080 -e PORT=8080 -e APP_ENV=dev --name gotodos gotodos:0.0.1
docker run --rm -p 8080:8080 -e PORT=8080 -e APP_ENV=production -v $(pwd)/logs:/logs --name gotodos gotodos:0.0.1
```

```
docker run --rm -p 8080:8080 -e PORT=8080 -e APP_ENV=dev -e DB_USER=postgres \
-e DB_PASSWORD=passw0rd -e DB_NAME=todos -e DB_PORT=5432 -e DB_HOST=host.docker.internal \
-v $(pwd)/logs:/logs -v $(pwd)/uploads:/uploads --name gotodos gotodos:0.0.1
```


## ```build database```
```
docker compose -f database/docker-compose.yml up -d
```
## clean database
```
docker compose -f database/docker-compose.yml down 
```
