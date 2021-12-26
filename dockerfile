FROM golang:1.17.5-alpine3.15

RUN mkdir /app
ADD . /app
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN go build -o main .

EXPOSE 8080

CMD ["/app/main"]