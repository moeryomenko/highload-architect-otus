# Social Network

Monolithic service with the backbone of a social network implementation.

## Usage

Run service and db.

```bash
$ cd deployments/compose
$ docker compose up --build
```

Request examples

```bash
$ curl -X POST http://localhost:8080/api/v1/signup -d '{"nickname":"moeryomenko", "password":"test"}'
# example of response see api/openapi.yaml
$ curl -X PUT http://localhost:8080/api/v1/profile -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzb2NpYWwiLCJleHAiOjE2NDg0NDgxMjMsImlhdCI6MTY0ODQ0ODEyMywidXNlcl9pZCI6IjIwMWY5M2VjLWRmYWUtNGIzYy04MDE5LTc4OGQwMzNhZTZhMCJ9.ed9pOxHM-AlyZWOuqdoDgQi7zf3AMWXX-cauBLwB-2I' -d '{"age":26,"city":"Krasnodar","first_name":"Maxim","last_name":"Eryomenko","gender":"male","interests":["programming","music"]}'
```
