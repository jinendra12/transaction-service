FROM golang:1.23.4-alpine

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o main ./cmd/api

EXPOSE 8080

CMD ["./main"]
