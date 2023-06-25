FROM golang:1.19

ENV CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum ./
COPY . .

# RUN go mod download

CMD ["go", "test", "-v", "./..."]