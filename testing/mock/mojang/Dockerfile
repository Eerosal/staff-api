FROM golang:1.19

ENV CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum mojang.go ./

RUN go build -o /app/main .

CMD ["/app/main"]