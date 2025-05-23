FROM golang:1.23.7-alpine3.21

WORKDIR /app

RUN go install github.com/air-verse/air@latest

# COPY go.mod go.sum ./

# RUN go mod download

COPY go.* ./

RUN go mod download

COPY . .

# RUN go build -o main main.go

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]