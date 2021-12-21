```
go run main.go
```

#### Liveness Probe
```
cat /tmp/live
echo $?
```
```output = 0 ,status = up```

|       Method      |     RelativePath  | CURL |
| ------------- | ------------- | ------------- |
|   GET      | /healthz                  | ``` curl GET 127.0.0.1:8080/healthz ``` |
|   GET      | /x                        | ``` curl GET 127.0.0.1:8080/x ``` |