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

The app is also available to use with [Docker](https://docs.docker.com/engine/install/) and can be used either locally or on a service like [Docker Playground](https://labs.play-with-docker.com/). The image for the app can be found [here](https://hub.docker.com/r/fiddler46/city-api-go).

Simply pull:

```go
docker pull fiddler46/city-api-go:<tag>
```

and then run:

```go
docker run fiddler46/sms-rails:<tag>
```

It will then show that the MongoDB has started successfully and the application running at the specified port. All Postman rules from above still follow.
