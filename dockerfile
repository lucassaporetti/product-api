# Dockerfile
FROM golang:1.21.0

WORKDIR /

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o product-api

EXPOSE 8080

CMD ["./product-api"]