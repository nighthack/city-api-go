FROM golang:1.17.5-alpine3.15

RUN mkdir /app
ADD . /app
WORKDIR /app

# Copies and downloads necessary dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN go build -o main .

# Port 8080 exposed for use
EXPOSE 8080

# Command that starts the container
CMD ["/app/main"]