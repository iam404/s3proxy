FROM golang:1.19

USER root
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o ./s3proxy

EXPOSE 8080



CMD ["./s3proxy"]
