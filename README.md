# go-todos

## run kubernetes
```chmod +x k8s.sh
./k8s.sh 
kubectl get po -n gotodos 
```

```go
go run main.go
```

#### Liveness Probe

```
cat /tmp/live
echo $?
```

`output = 0 ,status = up`

| Method | RelativePath | CURL                             |
| ------ | ------------ | -------------------------------- |
| GET    | /healthz     | ` http://127.0.0.1:8080/healthz` |
| GET    | /x           | ` http://127.0.0.1:8080/x`       |

| Method | RelativePath      | CURL                                      |
| ------ | ----------------- | ----------------------------------------- |
| GET    | /api/v1/todos     | ` http://127.0.0.1:8080/api/v1/todos`     |
| GET    | /api/v1/todos/:id | ` http://127.0.0.1:8080/api/v1/todos/:id` |
| DELETE | /api/v1/todos/:id | ` http://127.0.0.1:8080/api/v1/todos/:id` |
| POST   | /api/v1/todos     | ` http://127.0.0.1:8080/api/v1/todos`     |
| PUT    | /api/v1/todos/:id | ` http://127.0.0.1:8080/api/v1/todos/:id` |

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

```
 curl -X POST -H "Content-Type: application/json" \
	  --data "{\"email\":\"sing@dev.com\",\"password\":\"passw0rd\"}" \
	  http://localhost:8080/api/v1/auth/sign-up
```

## `build database`

```
docker compose up -d
```

## clean database

```
docker compose down
```

```
docker run --rm -p 8080:8080 -e PORT=8080 -e APP_ENV=dev \
-e DSN='host=host.docker.internal user=postgres password=passw0rd dbname=todos port=5432  sslmode=disable TimeZone=Asia/Bangkok' \
		 -e REDIS_HOST=host.docker.internal \
		 -e JWT_SECRET_KEY=1F460676-D6C2-4E40-A93F-1F0790C0725C -v $(pwd)/uploads:/uploads --name gotodos sing3demons/gotodos:0.0.6
```
