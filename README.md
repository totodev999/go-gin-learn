# About

API server sample repository built with Go and Gin.

# To-do

1. Tests for detecting race condition
2. Use Postgres for tests
3. Find another tool to manage migration histories

## Setup

1. install air  
   `go install github.com/air-verse/air@latest`
2. install testsum  
   `go install gotest.tools/gotestsum@latest`
3. install golangci-lint  
   `golangci-lint`
4. create database  
   `CREATE DATABASE fleamarket;`
5. create .env and .env.test

`.env`

```
ENV="prod"
DB_HOST="localhost"
DB_USER="postgres"
DB_PASSWORD="password"
DB_NAME="fleamarket"
DB_PORT="5432"
# openssl rand -hex 32
SECRET_KEY="47ed8f16a737a02a43b5211703e6452288961a15b4bebe5683ac862176df515b"
```

`.env.test`

```
ENV="test"
SECRET_KEY="test"
```

6. migrate  
   `make migrate`
7. run server  
   `make run`
