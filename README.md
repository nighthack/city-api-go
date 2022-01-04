# city-api-go

## About this Project

A microservice built to handle user registration and searching a database for a list of cities based on a user query. MongoDB is used to store and hold all user details as well as the list of existing cities to filter through. Built with [Go](https://go.dev/).

## Getting Started

### Prequisites

- [Go](https://go.dev/)
- [gorilla/mux](https://github.com/gorilla/mux)
- [MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver)
- [GoDotEnv](https://github.com/joho/godotenv)

Any other remaining go packages, including the ones above, can be downloaded by running:

```go
go mod download
```

### Running

Simply run

```go
go run main.go
```

The available routes and how to access them can be found here in this [Postman Documentation](https://documenter.getpostman.com/view/12592433/UVREjQPW)

### Running on Docker

Import the DB on terminal

```shell
./localmon.sh
```
Run docker compose

```shell
docker-compose up
```

Insert appropriate request to search. Ensure API access credentials on x-api-key header.

```shell
http://localhost:54321/suggest?city_name=<insertcityname>
```
