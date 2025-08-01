# About

API server sample repository built with Go and Gin.

# To-do

1. Find another tool to manage migration histories  
   Atlas can be used with Gorm
2. Add feature for Sigterm
   `https://oiasnak.hatenablog.com/entry/2023/09/10/003752`
3. switch from JWT to session

## Setup

1. install air  
   `go install github.com/air-verse/air@latest`
2. install testsum  
   `go install gotest.tools/gotestsum@latest`
3. install golangci-lint  
   `brew install golangci-lint`
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

### Memo

In GoLang, there is no method like asyncLocalStorage in Node. Is it better to pass context to service and repository for logging in a better way?
